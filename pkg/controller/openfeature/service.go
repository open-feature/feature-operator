package openfeature

import (
	"context"
	"fmt"
	featurev1 "github.com/open-feature/feature-operator/pkg/apis/open-feature.dev/v1alpha1"
	"github.com/open-feature/feature-operator/pkg/common"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *OpenFeatureReconciler) reconcileService(instance *featurev1.OpenFeature, ctx context.Context) (bool, error) {
	newService, err := r.makeService(instance)
	if err != nil {
		return false, err
	}

	service := &corev1.Service{}

	service, created, err := r.getCreateService(ctx, common.FLAG_API_SERVICE_NAME, instance, service, newService)
	if err != nil {
		return false, err
	} else if created {
		return true, nil
	}

	err = r.compareUpdateService(ctx, instance, service, newService)
	if err != nil {
		return false, nil
	}
	return false, nil
}

func (r *OpenFeatureReconciler) getCreateService(ctx context.Context, name string, instance *featurev1.OpenFeature, obj *corev1.Service, newObj *corev1.Service) (*corev1.Service, bool, error) {
	outObj := &corev1.Service{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: name, Namespace: instance.Namespace}, outObj)
	if err != nil && errors.IsNotFound(err) {
		err = r.Client.Create(ctx, newObj)
		if err != nil {
			r.Recorder.Event(instance, "Warning", "FailedCreate", fmt.Sprintf("Failed to create %s %s/%s (Reason: %s changed)", "Service", instance.Namespace, common.FLAG_API_SERVICE_NAME, common.CRD_OPENFEATURE_NAME))
			return outObj, false, fmt.Errorf("failed to create new service: %w", err)
		}
		r.Recorder.Event(instance, "Normal", "Created", fmt.Sprintf("Created %s %s/%s (Reason: %s changed)", "Service", instance.Namespace, common.FLAG_API_SERVICE_NAME, common.CRD_OPENFEATURE_NAME))
		return outObj, true, nil
	} else if err != nil {
		return outObj, false, fmt.Errorf("failed to get service: %w", err)
	}
	return outObj, false, nil
}

func (r *OpenFeatureReconciler) compareUpdateService(ctx context.Context, instance *featurev1.OpenFeature, obj *corev1.Service, newObj *corev1.Service) error {
	const TYPE = "Service"
	if !common.CompareHashStructure(obj.ObjectMeta.Annotations[common.APPLIED_HASH_ANNOTATION], newObj.Spec) {
		obj.Spec = newObj.Spec
		obj.ObjectMeta.Annotations[common.APPLIED_HASH_ANNOTATION] = common.GetHashStructure(newObj.Spec)
		err := r.Client.Update(ctx, obj)
		if err != nil {
			r.Log.Error(err, "Failed to update "+TYPE, common.CRD_OPENFEATURE_NAME+".Namespace", instance.Namespace, common.CRD_OPENFEATURE_NAME+".Name", common.FLAG_API_SERVICE_NAME)
			r.Recorder.Event(instance, "Warning", "FailedUpdate", fmt.Sprintf("Failed to update %s %s/%s (Reason: %s changed)", TYPE, instance.Namespace, common.FLAG_API_SERVICE_NAME, common.CRD_OPENFEATURE_NAME))
			return err
		}
		r.Recorder.Event(instance, "Normal", "Updated", fmt.Sprintf("Updated %s %s/%s (Reason: %s changed)", TYPE, instance.Namespace, common.FLAG_API_SERVICE_NAME, common.CRD_OPENFEATURE_NAME))
		r.Log.Info(TYPE + " updated")
	}
	return nil
}

func (r *OpenFeatureReconciler) makeService(feature *featurev1.OpenFeature) (*corev1.Service, error) {
	labels := make(map[string]string)

	selectorLabels := map[string]string{
		"app":     common.FLAG_API_SERVICE_NAME,
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

	obj := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        common.FLAG_API_SERVICE_NAME,
			Labels:      labels,
			Annotations: annotations,
			Namespace:   feature.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port:     8080,
					Protocol: "TCP",
				},
			},
			Selector: common.GetCommonLabels(*feature),
			Type:     "ClusterIP",
		},
	}

	obj.ObjectMeta.Annotations[common.APPLIED_HASH_ANNOTATION] = common.GetHashStructure(obj.Spec)

	err := controllerutil.SetControllerReference(feature, obj, r.Scheme)
	if err != nil {
		return nil, fmt.Errorf("could not set controller reference: %w", err)
	}

	return obj, err
}
