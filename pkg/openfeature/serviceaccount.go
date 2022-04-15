package openfeature

import (
	"context"
	featurev1 "github.com/open-feature/feature-operator/pkg/apis/open-feature.dev/v1alpha1"
	"github.com/open-feature/feature-operator/pkg/common"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *OpenFeatureReconciler) reconcileServiceAccount(ctx context.Context, openFeature *featurev1.OpenFeature) (*ctrl.Result, error) {
	return nil, nil
}

func makeServiceAccount(feature featurev1.OpenFeature) *corev1.ServiceAccount {
	labels := make(map[string]string)
	annotations := make(map[string]string)

	for k, v := range common.GetCommonLabels(feature) {
		labels[k] = v

	}

	for k, v := range feature.Spec.Labels {
		labels[k] = v
	}

	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:        common.FLAG_API_SERVICEACCOUNT_NAME,
			Labels:      labels,
			Annotations: annotations,
			Namespace:   feature.Namespace,
		},
	}
}
