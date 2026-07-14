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

// OIDCAdvancedGroupMapperSpec defines the desired state of OIDCAdvancedGroupMapper.
//
// This mirrors Keycloak's built-in "Advanced Claim to Group" identity
// provider mapper (provider ID "oidc-advanced-group-idp-mapper"): if all
// configured claims match, the federated user is assigned to the given
// group.
type OIDCAdvancedGroupMapperSpec struct {
	// identityProviderRef references the OIDCIdentityProvider this mapper is attached to.
	// +required
	IdentityProviderRef *common.IdentityProviderRef `json:"identityProviderRef"`

	// name is the mapper's name as shown in the Keycloak admin console.
	// +required
	Name string `json:"name"`

	// claims maps claim names to the expected value that must be present in
	// the token for the group assignment to apply. Nested claims may be
	// referenced using a "." separator (e.g. "address.locality").
	// +required
	Claims map[string]string `json:"claims"`

	// claimValuesRegex controls whether the configured claim values are
	// interpreted as regular expressions. Defaults to false.
	// +optional
	ClaimValuesRegex *bool `json:"claimValuesRegex,omitempty"`

	// group is the full path of the Keycloak group to assign the user to
	// when all claims match (e.g. "/parent/child").
	// +required
	Group string `json:"group"`
}

// OIDCAdvancedGroupMapperStatus defines the observed state of OIDCAdvancedGroupMapper.
type OIDCAdvancedGroupMapperStatus struct {
	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// keycloakID is the UUID assigned by Keycloak for this identity provider mapper.
	// +optional
	KeycloakID *string `json:"keycloakID,omitempty"`

	// conditions represent the current state of the OIDCAdvancedGroupMapper resource.
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
// +kubebuilder:printcolumn:name="IdentityProvider",type=string,JSONPath=`.spec.identityProviderRef.name`
// +kubebuilder:printcolumn:name="Name",type=string,JSONPath=`.spec.name`
// +kubebuilder:printcolumn:name="Group",type=string,JSONPath=`.spec.group`
// +kubebuilder:printcolumn:name="KeycloakID",type=string,JSONPath=`.status.keycloakID`

// OIDCAdvancedGroupMapper is the Schema for the oidcadvancedgroupmappers API
type OIDCAdvancedGroupMapper struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of OIDCAdvancedGroupMapper
	// +required
	Spec OIDCAdvancedGroupMapperSpec `json:"spec"`

	// status defines the observed state of OIDCAdvancedGroupMapper
	// +optional
	Status OIDCAdvancedGroupMapperStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// OIDCAdvancedGroupMapperList contains a list of OIDCAdvancedGroupMapper
type OIDCAdvancedGroupMapperList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OIDCAdvancedGroupMapper `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OIDCAdvancedGroupMapper{}, &OIDCAdvancedGroupMapperList{})
}

// GetKeycloakID returns the Keycloak-assigned UUID stored in the status.
func (m *OIDCAdvancedGroupMapper) GetKeycloakID() *string { return m.Status.KeycloakID }

// SetKeycloakID stores a Keycloak-assigned UUID in the status.
func (m *OIDCAdvancedGroupMapper) SetKeycloakID(id *string) { m.Status.KeycloakID = id }
