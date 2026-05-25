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

// RealmLoginSettings groups login-page and registration settings for the realm.
type RealmLoginSettings struct {
	// registrationAllowed allows users to self-register.
	// +optional
	RegistrationAllowed *bool `json:"registrationAllowed,omitempty"`

	// registrationEmailAsUsername uses the email address as username during
	// self-registration.
	// +optional
	RegistrationEmailAsUsername *bool `json:"registrationEmailAsUsername,omitempty"`

	// editUsernameAllowed allows users to change their own username.
	// +optional
	EditUsernameAllowed *bool `json:"editUsernameAllowed,omitempty"`

	// resetPasswordAllowed allows users to reset their password via email.
	// +optional
	ResetPasswordAllowed *bool `json:"resetPasswordAllowed,omitempty"`

	// rememberMe enables the "Remember me" login option.
	// +optional
	RememberMe *bool `json:"rememberMe,omitempty"`

	// verifyEmail requires new users to verify their email address.
	// +optional
	VerifyEmail *bool `json:"verifyEmail,omitempty"`

	// loginWithEmailAllowed allows users to authenticate using their email.
	// +optional
	LoginWithEmailAllowed *bool `json:"loginWithEmailAllowed,omitempty"`

	// duplicateEmailsAllowed allows multiple accounts to share the same email.
	// +optional
	DuplicateEmailsAllowed *bool `json:"duplicateEmailsAllowed,omitempty"`
}

// RealmThemes groups UI theme settings for the realm.
type RealmThemes struct {
	// login sets the theme used on the login page.
	// +optional
	Login *string `json:"login,omitempty"`

	// account sets the theme used on the account management page.
	// +optional
	Account *string `json:"account,omitempty"`

	// admin sets the theme used in the admin console.
	// +optional
	Admin *string `json:"admin,omitempty"`

	// email sets the theme used for emails sent by the realm.
	// +optional
	Email *string `json:"email,omitempty"`
}

// RealmInternationalization groups i18n settings for the realm.
type RealmInternationalization struct {
	// enabled enables multi-language support.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// defaultLocale sets the default locale when internationalization is enabled.
	// +optional
	DefaultLocale *string `json:"defaultLocale,omitempty"`

	// supportedLocales lists the locales available to users.
	// +optional
	SupportedLocales []string `json:"supportedLocales,omitempty"`
}

// RealmSessionTimeouts groups session and token lifetime settings for the realm.
type RealmSessionTimeouts struct {
	// accessTokenLifespan sets how long (in seconds) an access token is valid.
	// +optional
	AccessTokenLifespan *int `json:"accessTokenLifespan,omitempty"`

	// ssoSessionIdleTimeout sets the idle timeout (in seconds) for SSO sessions.
	// +optional
	SsoSessionIdleTimeout *int `json:"ssoSessionIdleTimeout,omitempty"`

	// ssoSessionMaxLifespan sets the maximum lifespan (in seconds) of an SSO
	// session regardless of activity.
	// +optional
	SsoSessionMaxLifespan *int `json:"ssoSessionMaxLifespan,omitempty"`

	// offlineSessionIdleTimeout sets the idle timeout (in seconds) for offline
	// sessions.
	// +optional
	OfflineSessionIdleTimeout *int `json:"offlineSessionIdleTimeout,omitempty"`
}

// RealmSmtpServer holds the SMTP server configuration for outgoing realm emails.
type RealmSmtpServer struct {
	// host is the SMTP server hostname or IP address.
	// +required
	Host string `json:"host"`

	// port is the SMTP server port. Defaults to 25 when omitted.
	// +optional
	Port *int `json:"port,omitempty"`

	// from is the envelope-from email address used for outgoing mail.
	// +required
	From string `json:"from"`

	// fromDisplayName is the human-readable name shown as the sender.
	// +optional
	FromDisplayName *string `json:"fromDisplayName,omitempty"`

	// replyTo is the reply-to email address.
	// +optional
	ReplyTo *string `json:"replyTo,omitempty"`

	// replyToDisplayName is the human-readable name for the reply-to address.
	// +optional
	ReplyToDisplayName *string `json:"replyToDisplayName,omitempty"`

	// envelopeFrom overrides the envelope sender address.
	// +optional
	EnvelopeFrom *string `json:"envelopeFrom,omitempty"`

	// ssl enables SSL/TLS for the SMTP connection.
	// +optional
	Ssl *bool `json:"ssl,omitempty"`

	// startTls enables STARTTLS negotiation.
	// +optional
	StartTls *bool `json:"startTls,omitempty"`

	// auth enables SMTP authentication.
	// +optional
	Auth *bool `json:"auth,omitempty"`

	// user is the SMTP authentication username.
	// +optional
	User *string `json:"user,omitempty"`

	// password references a secret or inline value for the SMTP authentication
	// password.
	// +optional
	Password *common.SecretKeyRefOrVal `json:"password,omitempty"`
}

// RealmSpec defines the desired state of Realm
type RealmSpec struct {
	// keycloakRef is a reference to the Keycloak instance that owns the realm.
	// +required
	KeycloakRef *common.KeycloakRef `json:"keycloakRef"`

	// realmName specifies the internal name of the realm.
	// This field is immutable after creation.
	// +required
	RealmName string `json:"realmName"`

	// displayName is a human-readable name shown in the UI.
	// +optional
	DisplayName *string `json:"displayName,omitempty"`

	// displayNameHtml is an HTML-formatted display name shown in the login page.
	// +optional
	DisplayNameHtml *string `json:"displayNameHtml,omitempty"`

	// enabled controls whether the realm is active.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// sslRequired sets the SSL/TLS enforcement policy.
	// +kubebuilder:validation:Enum=none;external;all
	// +optional
	SslRequired *string `json:"sslRequired,omitempty"`

	// passwordPolicy sets the password policy expression for the realm
	// (e.g. "length(8) and digits(1)").
	// +optional
	PasswordPolicy *string `json:"passwordPolicy,omitempty"`

	// bruteForceProtected enables automatic account lockout after repeated
	// failed login attempts.
	// +optional
	BruteForceProtected *bool `json:"bruteForceProtected,omitempty"`

	// login groups login-page and registration behaviour settings.
	// +optional
	Login *RealmLoginSettings `json:"login,omitempty"`

	// themes groups UI theme assignments for the realm.
	// +optional
	Themes *RealmThemes `json:"themes,omitempty"`

	// internationalization groups locale and multi-language settings.
	// +optional
	Internationalization *RealmInternationalization `json:"internationalization,omitempty"`

	// sessionTimeouts groups session and token lifetime settings.
	// +optional
	SessionTimeouts *RealmSessionTimeouts `json:"sessionTimeouts,omitempty"`

	// smtpServer configures the SMTP server used for outgoing realm emails.
	// +optional
	SmtpServer *RealmSmtpServer `json:"smtpServer,omitempty"`

	// attributes holds arbitrary realm-level key/value attributes.
	// +optional
	Attributes map[string]string `json:"attributes,omitempty"`
}

// RealmStatus defines the observed state of Realm.
type RealmStatus struct {
	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// keycloakID is the UUID assigned by Keycloak for this realm.
	// +optional
	KeycloakID *string `json:"keycloakID,omitempty"`

	// conditions represent the current state of the Realm resource.
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
// +kubebuilder:printcolumn:name="Realm",type=string,JSONPath=`.spec.realmName`
// +kubebuilder:printcolumn:name="Enabled",type=boolean,JSONPath=`.spec.enabled`
// +kubebuilder:printcolumn:name="KeycloakID",type=string,JSONPath=`.status.keycloakID`

// Realm is the Schema for the realms API
type Realm struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of Realm
	// +required
	Spec RealmSpec `json:"spec"`

	// status defines the observed state of Realm
	// +optional
	Status RealmStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// RealmList contains a list of Realm
type RealmList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Realm `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Realm{}, &RealmList{})
}

// GetKeycloakID returns the Keycloak-assigned UUID stored in the status.
func (r *Realm) GetKeycloakID() *string { return r.Status.KeycloakID }

// SetKeycloakID stores a Keycloak-assigned UUID in the status.
func (r *Realm) SetKeycloakID(id *string) { r.Status.KeycloakID = id }
