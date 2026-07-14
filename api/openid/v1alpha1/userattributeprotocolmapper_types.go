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

// UserAttributeProtocolMapperSpec defines the desired state of UserAttributeProtocolMapper.
type UserAttributeProtocolMapperSpec struct {
	// clientRef references the OpenIDClient this protocol mapper is attached to.
	// +required
	ClientRef *common.ClientRef `json:"clientRef"`

	// name is the mapper's name as shown in the Keycloak admin console.
	// +required
	Name string `json:"name"`

	// userAttribute is the custom Keycloak user attribute to read.
	// +required
	UserAttribute string `json:"userAttribute"`

	// claimName is the token claim that receives the mapped value.
	// +required
	ClaimName string `json:"claimName"`

	// claimValueType defines the JSON type label used for the claim.
	// +kubebuilder:validation:Enum=String;JSON;long;int;boolean
	// +optional
	ClaimValueType *string `json:"claimValueType,omitempty"`

	// multivalued controls whether multi-valued claims are emitted.
	// +optional
	Multivalued *bool `json:"multivalued,omitempty"`

	// aggregateAttributes controls whether attributes with the same name are aggregated.
	// +optional
	AggregateAttributes *bool `json:"aggregateAttributes,omitempty"`

	// addToIDToken controls whether the claim is included in ID tokens.
	// +optional
	AddToIDToken *bool `json:"addToIDToken,omitempty"`

	// addToAccessToken controls whether the claim is included in access tokens.
	// +optional
	AddToAccessToken *bool `json:"addToAccessToken,omitempty"`

	// addToUserInfo controls whether the claim is included in the userinfo response.
	// +optional
	AddToUserInfo *bool `json:"addToUserInfo,omitempty"`

	// addToTokenIntrospection controls whether the claim is included in token introspection responses.
	// +optional
	AddToTokenIntrospection *bool `json:"addToTokenIntrospection,omitempty"`
}

// UserAttributeProtocolMapperStatus defines the observed state of UserAttributeProtocolMapper.
type UserAttributeProtocolMapperStatus struct {
	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// keycloakID is the UUID assigned by Keycloak for this protocol mapper.
	// +optional
	KeycloakID *string `json:"keycloakID,omitempty"`

	// conditions represent the current state of the UserAttributeProtocolMapper resource.
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
// +kubebuilder:printcolumn:name="Client",type=string,JSONPath=`.spec.clientRef.name`
// +kubebuilder:printcolumn:name="Name",type=string,JSONPath=`.spec.name`
// +kubebuilder:printcolumn:name="KeycloakID",type=string,JSONPath=`.status.keycloakID`

// UserAttributeProtocolMapper is the Schema for the userattributeprotocolmappers API.
type UserAttributeProtocolMapper struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of UserAttributeProtocolMapper
	// +required
	Spec UserAttributeProtocolMapperSpec `json:"spec"`

	// status defines the observed state of UserAttributeProtocolMapper
	// +optional
	Status UserAttributeProtocolMapperStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// UserAttributeProtocolMapperList contains a list of UserAttributeProtocolMapper.
type UserAttributeProtocolMapperList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []UserAttributeProtocolMapper `json:"items"`
}

func init() {
	SchemeBuilder.Register(&UserAttributeProtocolMapper{}, &UserAttributeProtocolMapperList{})
}

// GetKeycloakID returns the Keycloak-assigned UUID stored in the status.
func (m *UserAttributeProtocolMapper) GetKeycloakID() *string { return m.Status.KeycloakID }

// SetKeycloakID stores a Keycloak-assigned UUID in the status.
func (m *UserAttributeProtocolMapper) SetKeycloakID(id *string) { m.Status.KeycloakID = id }
