package openfeature

import (
	"context"
	featurev1 "github.com/open-feature/feature-operator/pkg/apis/open-feature.dev/v1alpha1"
	"github.com/open-feature/feature-operator/pkg/common"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *OpenFeatureReconciler) reconcileService(ctx context.Context, openFeature *featurev1.OpenFeature) (*ctrl.Result, error) {
	return nil, nil
}

func makeService(feature featurev1.OpenFeature) *corev1.Service {
	labels := make(map[string]string)

	selectorLabels := map[string]string{
		"app":     common.FLAG_API_SERVICE_NAME,
		"version": feature.Spec.Version,
	}

	for k, v := range common.GetCommonLabels(feature) {
		labels[k] = v
	}

	for k, v := range feature.Spec.Labels {
		labels[k] = v
	}

	for k, v := range selectorLabels {
		labels[k] = v
	}

	annotations := map[string]string{}

	return &corev1.Service{
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
			Selector: common.GetCommonLabels(feature),
			Type:     "ClusterIP",
		},
	}
}
