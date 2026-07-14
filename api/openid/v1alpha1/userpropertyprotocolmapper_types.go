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

// UserPropertyProtocolMapperSpec defines the desired state of UserPropertyProtocolMapper.
type UserPropertyProtocolMapperSpec struct {
	// clientRef references the OpenIDClient this protocol mapper is attached to.
	// +required
	ClientRef *common.ClientRef `json:"clientRef"`

	// name is the mapper's name as shown in the Keycloak admin console.
	// +required
	Name string `json:"name"`

	// userProperty is the built-in Keycloak user property to read.
	// +required
	UserProperty string `json:"userProperty"`

	// claimName is the token claim that receives the mapped value.
	// +required
	ClaimName string `json:"claimName"`

	// claimValueType defines the JSON type label used for the claim.
	// +kubebuilder:validation:Enum=String;JSON;long;int;boolean
	// +optional
	ClaimValueType *string `json:"claimValueType,omitempty"`

	// addToIDToken controls whether the claim is included in ID tokens.
	// +optional
	AddToIDToken *bool `json:"addToIDToken,omitempty"`

	// addToAccessToken controls whether the claim is included in access tokens.
	// +optional
	AddToAccessToken *bool `json:"addToAccessToken,omitempty"`

	// addToUserInfo controls whether the claim is included in the userinfo response.
	// +optional
	AddToUserInfo *bool `json:"addToUserInfo,omitempty"`
}

// UserPropertyProtocolMapperStatus defines the observed state of UserPropertyProtocolMapper.
type UserPropertyProtocolMapperStatus struct {
	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// keycloakID is the UUID assigned by Keycloak for this protocol mapper.
	// +optional
	KeycloakID *string `json:"keycloakID,omitempty"`

	// conditions represent the current state of the UserPropertyProtocolMapper resource.
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

// UserPropertyProtocolMapper is the Schema for the userpropertyprotocolmappers API.
type UserPropertyProtocolMapper struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of UserPropertyProtocolMapper
	// +required
	Spec UserPropertyProtocolMapperSpec `json:"spec"`

	// status defines the observed state of UserPropertyProtocolMapper
	// +optional
	Status UserPropertyProtocolMapperStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// UserPropertyProtocolMapperList contains a list of UserPropertyProtocolMapper.
type UserPropertyProtocolMapperList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []UserPropertyProtocolMapper `json:"items"`
}

func init() {
	SchemeBuilder.Register(&UserPropertyProtocolMapper{}, &UserPropertyProtocolMapperList{})
}

// GetKeycloakID returns the Keycloak-assigned UUID stored in the status.
func (m *UserPropertyProtocolMapper) GetKeycloakID() *string { return m.Status.KeycloakID }

// SetKeycloakID stores a Keycloak-assigned UUID in the status.
func (m *UserPropertyProtocolMapper) SetKeycloakID(id *string) { m.Status.KeycloakID = id }
