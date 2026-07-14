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

// OIDCIdentityProviderSpec defines the desired state of OIDCIdentityProvider
type OIDCIdentityProviderSpec struct {
	// realmRef is a reference to the Realm this identity provider belongs to.
	// +required
	RealmRef *common.RealmRef `json:"realmRef"`

	// alias is the unique identifier of the identity provider within the realm.
	// This field is immutable after creation.
	// +required
	Alias string `json:"alias"`

	// displayName is the human-friendly name shown in the login form.
	// +optional
	DisplayName *string `json:"displayName,omitempty"`

	// enabled controls whether the identity provider is active.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// authorizationURL is the OIDC authorization endpoint.
	// +required
	AuthorizationURL string `json:"authorizationURL"`

	// tokenURL is the OIDC token endpoint.
	// +required
	TokenURL string `json:"tokenURL"`

	// clientID is the OIDC client identifier registered with the external provider.
	// +required
	ClientID string `json:"clientID"`

	// clientSecret references a secret or inline value holding the OIDC client secret.
	// +required
	ClientSecret *common.SecretKeyRefOrVal `json:"clientSecret"`

	// userInfoURL is the OIDC user info endpoint.
	// +optional
	UserInfoURL *string `json:"userInfoURL,omitempty"`

	// jwksURL is the JSON Web Key Set URL used to validate provider signatures.
	// +optional
	JwksURL *string `json:"jwksURL,omitempty"`

	// logoutURL is the OIDC end-session (logout) endpoint.
	// +optional
	LogoutURL *string `json:"logoutURL,omitempty"`

	// issuer is the expected OIDC issuer identifier.
	// +optional
	Issuer *string `json:"issuer,omitempty"`

	// defaultScopes is the space-separated list of scopes requested during
	// authorization. Defaults to "openid" when omitted.
	// +optional
	DefaultScopes *string `json:"defaultScopes,omitempty"`

	// validateSignature enables signature validation of tokens issued by the
	// external provider.
	// +optional
	ValidateSignature *bool `json:"validateSignature,omitempty"`

	// backchannelSupported indicates whether the external provider supports
	// backchannel logout.
	// +optional
	BackchannelSupported *bool `json:"backchannelSupported,omitempty"`

	// hideOnLoginPage hides the identity provider from the realm's login page.
	// +optional
	HideOnLoginPage *bool `json:"hideOnLoginPage,omitempty"`

	// storeToken enables storing tokens issued by the external provider.
	// +optional
	StoreToken *bool `json:"storeToken,omitempty"`

	// trustEmail trusts email addresses asserted by the external provider,
	// skipping Keycloak's own email verification.
	// +optional
	TrustEmail *bool `json:"trustEmail,omitempty"`

	// linkOnly restricts the identity provider to account linking only,
	// preventing it from being used for direct authentication.
	// +optional
	LinkOnly *bool `json:"linkOnly,omitempty"`

	// firstBrokerLoginFlowAlias is the authentication flow used the first
	// time a user logs in through this identity provider.
	// +optional
	FirstBrokerLoginFlowAlias *string `json:"firstBrokerLoginFlowAlias,omitempty"`

	// postBrokerLoginFlowAlias is the authentication flow executed after
	// every login through this identity provider.
	// +optional
	PostBrokerLoginFlowAlias *string `json:"postBrokerLoginFlowAlias,omitempty"`

	// syncMode controls when Keycloak imports/synchronizes user data from the
	// external provider (e.g. "IMPORT", "LEGACY", "FORCE").
	// +kubebuilder:validation:Enum=IMPORT;LEGACY;FORCE
	// +optional
	SyncMode *string `json:"syncMode,omitempty"`
}

// OIDCIdentityProviderStatus defines the observed state of OIDCIdentityProvider.
type OIDCIdentityProviderStatus struct {
	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// keycloakID is the internal UUID assigned by Keycloak for this identity provider.
	// +optional
	KeycloakID *string `json:"keycloakID,omitempty"`

	// conditions represent the current state of the OIDCIdentityProvider resource.
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
// +kubebuilder:printcolumn:name="Alias",type=string,JSONPath=`.spec.alias`
// +kubebuilder:printcolumn:name="Enabled",type=boolean,JSONPath=`.spec.enabled`
// +kubebuilder:printcolumn:name="KeycloakID",type=string,JSONPath=`.status.keycloakID`

// OIDCIdentityProvider is the Schema for the oidcidentityproviders API
type OIDCIdentityProvider struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of OIDCIdentityProvider
	// +required
	Spec OIDCIdentityProviderSpec `json:"spec"`

	// status defines the observed state of OIDCIdentityProvider
	// +optional
	Status OIDCIdentityProviderStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// OIDCIdentityProviderList contains a list of OIDCIdentityProvider
type OIDCIdentityProviderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OIDCIdentityProvider `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OIDCIdentityProvider{}, &OIDCIdentityProviderList{})
}

// GetKeycloakID returns the Keycloak-assigned UUID stored in the status.
func (p *OIDCIdentityProvider) GetKeycloakID() *string { return p.Status.KeycloakID }

// SetKeycloakID stores a Keycloak-assigned UUID in the status.
func (p *OIDCIdentityProvider) SetKeycloakID(id *string) { p.Status.KeycloakID = id }
