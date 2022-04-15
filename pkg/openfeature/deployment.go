package openfeature

import (
	"context"
	featurev1 "github.com/open-feature/feature-operator/pkg/apis/open-feature.dev/v1alpha1"
	"github.com/open-feature/feature-operator/pkg/common"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *OpenFeatureReconciler) reconcileDeployment(ctx context.Context, openFeature *featurev1.OpenFeature) (*ctrl.Result, error) {
	return nil, nil
}

func makeDeployment(feature featurev1.OpenFeature) *appsv1.Deployment {
	labels := make(map[string]string)

	selectorLabels := map[string]string{
		"app":     common.FLAG_API_DEPLOYMENT_NAME,
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

	return &appsv1.Deployment{
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
}
