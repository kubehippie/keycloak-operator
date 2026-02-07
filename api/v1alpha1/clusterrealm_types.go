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
	"github.com/kubehippie/keycloak-operator/api/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterRealmSpec defines the desired state of ClusterRealm
type ClusterRealmSpec struct {
	// KeycloakRef is a reference to the Keycloak instance that owns the realm.
	// +required
	KeycloakRef *common.KeycloakRef `json:"keycloakRef"`

	// realmName specifies the name of the realm.
	// +required
	RealmName string `json:"realmName"`
}

// ClusterRealmStatus defines the observed state of ClusterRealm.
type ClusterRealmStatus struct {
	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// conditions represent the current state of the ClusterRealm resource.
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
// +kubebuilder:resource:scope=Cluster

// ClusterRealm is the Schema for the clusterrealms API
type ClusterRealm struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of ClusterRealm
	// +required
	Spec ClusterRealmSpec `json:"spec"`

	// status defines the observed state of ClusterRealm
	// +optional
	Status ClusterRealmStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// ClusterRealmList contains a list of ClusterRealm
type ClusterRealmList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterRealm `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterRealm{}, &ClusterRealmList{})
}
