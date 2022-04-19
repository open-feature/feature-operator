package openfeature

import (
	"context"
	"fmt"
	featurev1 "github.com/open-feature/feature-operator/pkg/apis/open-feature.dev/v1alpha1"
	"github.com/open-feature/feature-operator/pkg/common"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *OpenFeatureReconciler) reconcileDeployment(instance *featurev1.OpenFeature, ctx context.Context) (bool, error) {
	newDeploy, err := r.makeDeployment(instance)
	if err != nil {
		return false, err
	}

	deploy := &appsv1.Deployment{}

	deploy, created, err := r.getCreateDeployment(ctx, common.FLAG_API_DEPLOYMENT_NAME, instance, deploy, newDeploy)
	if err != nil {
		return false, err
	} else if created {
		return true, nil
	}

	err = r.compareUpdateDeployment(ctx, instance, deploy, newDeploy)
	if err != nil {
		return false, nil
	}
	return false, nil
}

func (r *OpenFeatureReconciler) getCreateDeployment(ctx context.Context, name string, instance *featurev1.OpenFeature, obj *appsv1.Deployment, newObj *appsv1.Deployment) (*appsv1.Deployment, bool, error) {
	outObj := &appsv1.Deployment{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: name, Namespace: instance.Namespace}, outObj)
	if err != nil && errors.IsNotFound(err) {
		err = r.Client.Create(ctx, newObj)
		if err != nil {
			r.Recorder.Event(instance, "Warning", "FailedCreate", fmt.Sprintf("Failed to create %s %s/%s (Reason: %s changed)", "Deployment", instance.Namespace, common.FLAG_API_DEPLOYMENT_NAME, common.CRD_OPENFEATURE_NAME))
			return outObj, false, fmt.Errorf("failed to create new deployment: %w", err)
		}
		r.Recorder.Event(instance, "Normal", "Created", fmt.Sprintf("Created %s %s/%s (Reason: %s changed)", "Deployment", instance.Namespace, common.FLAG_API_DEPLOYMENT_NAME, common.CRD_OPENFEATURE_NAME))
		return outObj, true, nil
	} else if err != nil {
		return outObj, false, fmt.Errorf("failed to get deployment: %w", err)
	}
	return outObj, false, nil
}

func (r *OpenFeatureReconciler) compareUpdateDeployment(ctx context.Context, instance *featurev1.OpenFeature, obj *appsv1.Deployment, newObj *appsv1.Deployment) error {
	const TYPE = "Deployment"
	if !common.CompareHashStructure(obj.ObjectMeta.Annotations[common.APPLIED_HASH_ANNOTATION], newObj.Spec) {
		if obj.ObjectMeta.Annotations == nil {
			obj.ObjectMeta.Annotations = make(map[string]string)
		}
		obj.Spec = newObj.Spec
		obj.ObjectMeta.Annotations[common.APPLIED_HASH_ANNOTATION] = common.GetHashStructure(newObj.Spec)
		err := r.Client.Update(ctx, obj)
		if err != nil {
			r.Log.Error(err, "Failed to update "+TYPE, common.CRD_OPENFEATURE_NAME+".Namespace", instance.Namespace, common.CRD_OPENFEATURE_NAME+".Name", common.FLAG_API_DEPLOYMENT_NAME)
			r.Recorder.Event(instance, "Warning", "FailedUpdate", fmt.Sprintf("Failed to update %s %s/%s (Reason: %s changed)", TYPE, instance.Namespace, common.FLAG_API_DEPLOYMENT_NAME, common.CRD_OPENFEATURE_NAME))
			return err
		}
		r.Recorder.Event(instance, "Normal", "Updated", fmt.Sprintf("Updated %s %s/%s (Reason: %s changed)", TYPE, instance.Namespace, common.FLAG_API_DEPLOYMENT_NAME, common.CRD_OPENFEATURE_NAME))
		r.Log.Info(TYPE + " updated")
	}
	return nil
}

func (r *OpenFeatureReconciler) makeDeployment(feature *featurev1.OpenFeature) (*appsv1.Deployment, error) {
	labels := make(map[string]string)

	selectorLabels := map[string]string{
		"app":     common.FLAG_API_DEPLOYMENT_NAME,
		"version": feature.Spec.Version,
	}

	for k, v := range common.GetCommonLabels(*feature) {
		labels[k] = v
	}

	for k, v := range feature.Spec.Labels {
		labels[k] = v
	}

	for k, v := range selectorLabels {
		labels[k] = v
	}

	annotations := map[string]string{}

	obj := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        common.FLAG_API_DEPLOYMENT_NAME,
			Labels:      labels,
			Annotations: annotations,
			Namespace:   feature.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &feature.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: selectorLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: annotations,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName:           common.FLAG_API_SERVICEACCOUNT_NAME,
					AutomountServiceAccountToken: common.Bool(true),
					Containers: []corev1.Container{
						{
							Name:      "openfeature-flag-api",
							Image:     fmt.Sprintf("%s:%s", feature.Spec.Image, feature.Spec.Version),
							Env:       feature.Spec.EnvVars,
							Resources: corev1.ResourceRequirements{},
						},
					},
					SecurityContext: &corev1.PodSecurityContext{
						RunAsUser:    common.Int64(65532),
						RunAsGroup:   common.Int64(65532),
						RunAsNonRoot: common.Bool(true),
						FSGroup:      common.Int64(65532),
					},
				},
			},
		},
	}
	obj.ObjectMeta.Annotations[common.APPLIED_HASH_ANNOTATION] = common.GetHashStructure(obj.Spec)

	err := controllerutil.SetControllerReference(feature, obj, r.Scheme)
	if err != nil {
		return nil, fmt.Errorf("could not set controller reference: %w", err)
	}

	return obj, nil
}
