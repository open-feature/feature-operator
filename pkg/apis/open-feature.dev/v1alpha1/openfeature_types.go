/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OpenFeatureSpec defines the desired state of OpenFeature
type OpenFeatureSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Image if specified has precedence over baseImage, tag and sha
	// combinations. Specifying the version is still necessary to ensure the
	// Prometheus Operator knows what version of Prometheus is being
	// configured.
	// +kubebuilder:default=ghcr.io/open-feature/feature-operator
	Image string `json:"image,omitempty"`

	// +kubebuilder:validation:enum=Deployment;Daemonset
	Version string `json:"version,omitempty"`

	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	EnvVars     []corev1.EnvVar   `json:"env-vars,omitempty"`

	// Foo is an example field of OpenFeature. Edit openfeature_types.go to remove/update
	// +kubebuilder:default=Deployment
	// +kubebuilder:validation:enum=Deployment;Daemonset
	DeploymentType string `json:"foo,omitempty"`

	// +kubebuilder:default=3
	// +kubebuilder:validation:Minimum=1
	Replicas int32 `json:"replicas,omitempty"`

	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
}

// OpenFeatureStatus defines the observed state of OpenFeature
type OpenFeatureStatus struct {
	Version string   `json:"version,omitempty"`
	Pods    []string `json:"pods,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// OpenFeature is the Schema for the openfeatures API
type OpenFeature struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenFeatureSpec   `json:"spec,omitempty"`
	Status OpenFeatureStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// OpenFeatureList contains a list of OpenFeature
type OpenFeatureList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenFeature `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OpenFeature{}, &OpenFeatureList{})
}
