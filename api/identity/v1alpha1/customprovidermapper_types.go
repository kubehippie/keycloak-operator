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

// CustomProviderMapperSpec defines the desired state of CustomProviderMapper
type CustomProviderMapperSpec struct {
	// identityProviderRef references the OIDCIdentityProvider this mapper is attached to.
	// +required
	IdentityProviderRef *common.IdentityProviderRef `json:"identityProviderRef"`

	// name is the mapper's name as shown in the Keycloak admin console.
	// +required
	Name string `json:"name"`

	// mapperType is the free-form Keycloak identity provider mapper type
	// identifier, e.g. "oidc-username-idp-mapper" or a custom SPI provider ID.
	// +required
	MapperType string `json:"mapperType"`

	// config holds the arbitrary key/value configuration passed through to
	// the underlying Keycloak identity provider mapper.
	// +optional
	Config map[string]string `json:"config,omitempty"`
}

// CustomProviderMapperStatus defines the observed state of CustomProviderMapper.
type CustomProviderMapperStatus struct {
	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// keycloakID is the UUID assigned by Keycloak for this identity provider mapper.
	// +optional
	KeycloakID *string `json:"keycloakID,omitempty"`

	// conditions represent the current state of the CustomProviderMapper resource.
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
// +kubebuilder:printcolumn:name="MapperType",type=string,JSONPath=`.spec.mapperType`
// +kubebuilder:printcolumn:name="KeycloakID",type=string,JSONPath=`.status.keycloakID`

// CustomProviderMapper is the Schema for the customprovidermappers API
type CustomProviderMapper struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of CustomProviderMapper
	// +required
	Spec CustomProviderMapperSpec `json:"spec"`

	// status defines the observed state of CustomProviderMapper
	// +optional
	Status CustomProviderMapperStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// CustomProviderMapperList contains a list of CustomProviderMapper
type CustomProviderMapperList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CustomProviderMapper `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CustomProviderMapper{}, &CustomProviderMapperList{})
}

// GetKeycloakID returns the Keycloak-assigned UUID stored in the status.
func (m *CustomProviderMapper) GetKeycloakID() *string { return m.Status.KeycloakID }

// SetKeycloakID stores a Keycloak-assigned UUID in the status.
func (m *CustomProviderMapper) SetKeycloakID(id *string) { m.Status.KeycloakID = id }
