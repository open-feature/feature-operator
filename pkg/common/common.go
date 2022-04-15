package common

import (
	"github.com/go-logr/logr"
	featurev1 "github.com/open-feature/feature-operator/pkg/apis/open-feature.dev/v1alpha1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func InitLog(request reconcile.Request, controllerName string, logName string) logr.Logger {
	log := logf.Log.WithName(logName)
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling " + controllerName)

	return reqLogger
}

func MergeMaps(a, b map[string]string) map[string]string {
	for k, v := range b {
		if _, present := a[k]; !present {
			a[k] = v
		}
	}

	return a
}

func GetCommonLabels(feature featurev1.OpenFeature) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       "flag-evaluation-api",
		"app.kubernetes.io/instance":   "flag-evaluation-api-" + feature.Name,
		"app.kubernetes.io/version":    feature.Spec.Version,
		"app.kubernetes.io/component":  "flag-evaluation-api",
		"app.kubernetes.io/part-of":    "open-feature",
		"app.kubernetes.io/managed-by": "flag-controller",
	}
}
