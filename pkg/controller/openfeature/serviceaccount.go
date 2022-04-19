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

func (r *OpenFeatureReconciler) reconcileServiceAccount(instance *featurev1.OpenFeature, ctx context.Context) (bool, error) {
	serviceAccount := &corev1.ServiceAccount{}
	newServiceAccount, err := r.makeServiceAccount(instance)

	if err != nil {
		return false, err
	}

	_, created, err := r.getCreateServiceAccount(ctx, common.FLAG_API_SERVICE_NAME, instance, serviceAccount, newServiceAccount)
	if err != nil {
		return false, err
	} else if created {
		return true, nil
	}
	return false, nil
}

func (r *OpenFeatureReconciler) getCreateServiceAccount(ctx context.Context, name string, instance *featurev1.OpenFeature, obj *corev1.ServiceAccount, newObj *corev1.ServiceAccount) (*corev1.ServiceAccount, bool, error) {
	outObj := &corev1.ServiceAccount{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: name, Namespace: instance.Namespace}, obj)
	if err != nil && errors.IsNotFound(err) {
		err = r.Client.Create(ctx, newObj)
		if err != nil {
			r.Recorder.Event(instance, "Warning", "FailedCreate", fmt.Sprintf("Failed to create %s %s/%s (Reason: %s changed)", "ServiceAccount", instance.Namespace, common.FLAG_API_SERVICEACCOUNT_NAME, common.CRD_OPENFEATURE_NAME))
			return outObj, false, fmt.Errorf("failed to create new ServiceAccount: %w", err)
		}
		r.Recorder.Event(instance, "Normal", "Created", fmt.Sprintf("Created %s %s/%s (Reason: %s changed)", "ServiceAccount", instance.Namespace, common.FLAG_API_SERVICEACCOUNT_NAME, common.CRD_OPENFEATURE_NAME))
		return outObj, true, nil
	} else if err != nil {
		return outObj, false, fmt.Errorf("failed to get ServiceAccount: %w", err)
	}
	return outObj, false, nil
}

func (r *OpenFeatureReconciler) makeServiceAccount(feature *featurev1.OpenFeature) (*corev1.ServiceAccount, error) {
	labels := make(map[string]string)
	annotations := make(map[string]string)

	for k, v := range common.GetCommonLabels(*feature) {
		labels[k] = v

	}

	for k, v := range feature.Spec.Labels {
		labels[k] = v
	}

	obj := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:        common.FLAG_API_SERVICEACCOUNT_NAME,
			Labels:      labels,
			Annotations: annotations,
			Namespace:   feature.Namespace,
		},
	}

	err := controllerutil.SetControllerReference(feature, obj, r.Scheme)
	if err != nil {
		return nil, fmt.Errorf("could not set controller reference: %w", err)
	}

	return obj, err

}
