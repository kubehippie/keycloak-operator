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

// OpenIDClientSpec defines the desired state of OpenIDClient
type OpenIDClientSpec struct {
	// realmRef is a reference to the Realm this client belongs to.
	// +required
	RealmRef *common.RealmRef `json:"realmRef"`

	// clientID is the client identifier registered in Keycloak.
	// This field is immutable after creation.
	// +required
	ClientID string `json:"clientID"`

	// name is the human-readable name of the client shown in the admin console.
	// +optional
	Name *string `json:"name,omitempty"`

	// description describes the purpose of the client.
	// +optional
	Description *string `json:"description,omitempty"`

	// enabled controls whether the client is active.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// accessType controls the client's authentication requirements.
	// +kubebuilder:validation:Enum=CONFIDENTIAL;PUBLIC;BEARER-ONLY
	// +required
	AccessType string `json:"accessType"`

	// clientSecret references a secret or inline value holding the client
	// secret. Only used when accessType is CONFIDENTIAL.
	// +optional
	ClientSecret *common.SecretKeyRefOrVal `json:"clientSecret,omitempty"`

	// clientAuthenticatorType sets the authenticator used to validate client
	// credentials. Defaults to "client-secret" when omitted.
	// +optional
	ClientAuthenticatorType *string `json:"clientAuthenticatorType,omitempty"`

	// standardFlowEnabled enables the OIDC authorization code flow.
	// +optional
	StandardFlowEnabled *bool `json:"standardFlowEnabled,omitempty"`

	// implicitFlowEnabled enables the OIDC implicit flow.
	// +optional
	ImplicitFlowEnabled *bool `json:"implicitFlowEnabled,omitempty"`

	// directAccessGrantsEnabled enables the OIDC resource owner password
	// credentials grant.
	// +optional
	DirectAccessGrantsEnabled *bool `json:"directAccessGrantsEnabled,omitempty"`

	// serviceAccountsEnabled enables the OIDC client credentials grant.
	// +optional
	ServiceAccountsEnabled *bool `json:"serviceAccountsEnabled,omitempty"`

	// frontChannelLogoutEnabled enables front-channel logout for the client.
	// +optional
	FrontChannelLogoutEnabled *bool `json:"frontChannelLogoutEnabled,omitempty"`

	// validRedirectURIs lists the URIs allowed as OIDC redirect targets.
	// +optional
	ValidRedirectURIs []string `json:"validRedirectURIs,omitempty"`

	// webOrigins lists the allowed CORS origins.
	// +optional
	WebOrigins []string `json:"webOrigins,omitempty"`

	// rootURL is the root URL used for relative redirect/base/admin URLs.
	// +optional
	RootURL *string `json:"rootURL,omitempty"`

	// baseURL is the default URL the client is redirected to after a
	// successful login.
	// +optional
	BaseURL *string `json:"baseURL,omitempty"`

	// adminURL is the URL used by Keycloak to send admin management
	// requests (e.g. logout, sync) to the client.
	// +optional
	AdminURL *string `json:"adminURL,omitempty"`

	// fullScopeAllowed enables all realm/client roles to be included in
	// access tokens by default.
	// +optional
	FullScopeAllowed *bool `json:"fullScopeAllowed,omitempty"`

	// attributes holds arbitrary client-level key/value attributes.
	// +optional
	Attributes map[string]string `json:"attributes,omitempty"`
}

// OpenIDClientStatus defines the observed state of OpenIDClient.
type OpenIDClientStatus struct {
	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// keycloakID is the internal UUID assigned by Keycloak for this client.
	// +optional
	KeycloakID *string `json:"keycloakID,omitempty"`

	// conditions represent the current state of the OpenIDClient resource.
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
// +kubebuilder:printcolumn:name="ClientID",type=string,JSONPath=`.spec.clientID`
// +kubebuilder:printcolumn:name="AccessType",type=string,JSONPath=`.spec.accessType`
// +kubebuilder:printcolumn:name="KeycloakID",type=string,JSONPath=`.status.keycloakID`

// OpenIDClient is the Schema for the openidclients API
type OpenIDClient struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of OpenIDClient
	// +required
	Spec OpenIDClientSpec `json:"spec"`

	// status defines the observed state of OpenIDClient
	// +optional
	Status OpenIDClientStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// OpenIDClientList contains a list of OpenIDClient
type OpenIDClientList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenIDClient `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OpenIDClient{}, &OpenIDClientList{})
}

// GetKeycloakID returns the Keycloak-assigned UUID stored in the status.
func (c *OpenIDClient) GetKeycloakID() *string { return c.Status.KeycloakID }

// SetKeycloakID stores a Keycloak-assigned UUID in the status.
func (c *OpenIDClient) SetKeycloakID(id *string) { c.Status.KeycloakID = id }
