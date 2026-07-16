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

// RoleSpec defines the desired state of Role
type RoleSpec struct {
	// realmRef is a reference to the Realm this role belongs to.
	// The Keycloak connection is resolved transitively through the Realm resource.
	// +required
	RealmRef *common.RealmRef `json:"realmRef"`

	// name specifies the realm-level role name within the realm.
	// +required
	Name string `json:"name"`

	// description holds a human-readable description of the role.
	// +optional
	Description *string `json:"description,omitempty"`

	// attributes holds custom role attributes as key/value pairs where each
	// key maps to one or more values.
	// +optional
	Attributes map[string][]string `json:"attributes,omitempty"`
}

// RoleStatus defines the observed state of Role.
type RoleStatus struct {
	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// keycloakID is the UUID assigned by Keycloak for this role.
	// It is stored here so that update and delete operations can reference the
	// role directly without an additional lookup by name.
	// +optional
	KeycloakID *string `json:"keycloakID,omitempty"`

	// conditions represent the current state of the Role resource.
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
// +kubebuilder:printcolumn:name="Realm",type=string,JSONPath=`.spec.realmRef.name`
// +kubebuilder:printcolumn:name="Name",type=string,JSONPath=`.spec.name`
// +kubebuilder:printcolumn:name="KeycloakID",type=string,JSONPath=`.status.keycloakID`

// Role is the Schema for the roles API
type Role struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of Role
	// +required
	Spec RoleSpec `json:"spec"`

	// status defines the observed state of Role
	// +optional
	Status RoleStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// RoleList contains a list of Role
type RoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Role `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Role{}, &RoleList{})
}

// GetKeycloakID returns the Keycloak-assigned UUID stored in the status.
func (r *Role) GetKeycloakID() *string { return r.Status.KeycloakID }

// SetKeycloakID stores a Keycloak-assigned UUID in the status.
func (r *Role) SetKeycloakID(id *string) { r.Status.KeycloakID = id }
