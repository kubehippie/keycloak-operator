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

// ClientScopeSpec defines the desired state of ClientScope, modeled after
// the terraform-provider-keycloak "keycloak_openid_client_scope" resource.
type ClientScopeSpec struct {
	// realmRef references the Realm this client scope belongs to. Immutable
	// after creation.
	// +required
	RealmRef *common.RealmRef `json:"realmRef"`

	// name is the name of the client scope.
	// +required
	Name string `json:"name"`

	// description is a human-readable description of the client scope.
	// +optional
	Description *string `json:"description,omitempty"`

	// consentScreenText is the text displayed on the consent screen when
	// this scope is requested. Setting this implies displayOnConsentScreen.
	// +optional
	ConsentScreenText *string `json:"consentScreenText,omitempty"`

	// includeInTokenScope controls whether this client scope is included in
	// the access token's "scope" claim. Defaults to true.
	// +optional
	IncludeInTokenScope *bool `json:"includeInTokenScope,omitempty"`
}

// ClientScopeStatus defines the observed state of ClientScope.
type ClientScopeStatus struct {
	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// keycloakID is the UUID assigned by Keycloak for this client scope.
	// +optional
	KeycloakID *string `json:"keycloakID,omitempty"`

	// conditions represent the current state of the ClientScope resource.
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

// ClientScope is the Schema for the clientscopes API
type ClientScope struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of ClientScope
	// +required
	Spec ClientScopeSpec `json:"spec"`

	// status defines the observed state of ClientScope
	// +optional
	Status ClientScopeStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// ClientScopeList contains a list of ClientScope
type ClientScopeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClientScope `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClientScope{}, &ClientScopeList{})
}

// GetKeycloakID returns the Keycloak-assigned UUID stored in the status.
func (c *ClientScope) GetKeycloakID() *string { return c.Status.KeycloakID }

// SetKeycloakID stores a Keycloak-assigned UUID in the status.
func (c *ClientScope) SetKeycloakID(id *string) { c.Status.KeycloakID = id }
