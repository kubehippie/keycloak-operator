/*
Copyright 2026 Thomas Boerger <thomas@webhippie.de>.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DefaultScopesSpec defines the desired state of DefaultScopes
type DefaultScopesSpec struct {
	// foo is an example field of DefaultScopes. Edit defaultscopes_types.go to remove/update
	// +optional
	Foo *string `json:"foo,omitempty"`
}

// DefaultScopesStatus defines the observed state of DefaultScopes.
type DefaultScopesStatus struct {
	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// conditions represent the current state of the DefaultScopes resource.
	// Each condition has a unique type and reflects the status of a specific aspect of the resource.
	//
	// Standard condition types include:
	// - "Available": the resource is fully functional
	// - "Progressing": the resource is being created or updated
	// - "Degraded": the resource failed to reach or maintain its desired state
	//
	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// DefaultScopes is the Schema for the defaultscopes API
type DefaultScopes struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of DefaultScopes
	// +required
	Spec DefaultScopesSpec `json:"spec"`

	// status defines the observed state of DefaultScopes
	// +optional
	Status DefaultScopesStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// DefaultScopesList contains a list of DefaultScopes
type DefaultScopesList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DefaultScopes `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DefaultScopes{}, &DefaultScopesList{})
}
