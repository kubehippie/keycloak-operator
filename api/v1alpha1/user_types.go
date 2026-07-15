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

// UserSpec defines the desired state of User
type UserSpec struct {
	// realmRef is a reference to the Realm this user belongs to.
	// The Keycloak connection is resolved transitively through the Realm resource.
	// +required
	RealmRef *common.RealmRef `json:"realmRef"`

	// username specifies the login name of the user.
	// +required
	Username string `json:"username"`

	// enabled specifies whether the user account is active.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// emailVerified specifies whether the user's email address has been verified.
	// +optional
	EmailVerified *bool `json:"emailVerified,omitempty"`

	// firstName specifies the user's first name.
	// +optional
	FirstName *string `json:"firstName,omitempty"`

	// lastName specifies the user's last name.
	// +optional
	LastName *string `json:"lastName,omitempty"`

	// email specifies the user's email address.
	// +optional
	Email *string `json:"email,omitempty"`

	// password references a secret or direct value which contains the initial
	// password to set for the user. If unset, no password is set in Keycloak
	// and the user cannot authenticate until an administrator sets one
	// (e.g. via the Keycloak console/API or a "forgot password" email flow).
	// +optional
	Password *common.SecretKeyRefOrVal `json:"password,omitempty"`

	// temporary specifies whether the password set via the password field is
	// temporary, forcing the user to change it on next login. Defaults to
	// true and is only relevant when password is set.
	// +optional
	// +kubebuilder:default=true
	Temporary *bool `json:"temporary,omitempty"`

	// attributes holds custom user attributes as key/value pairs where each
	// key maps to one or more values.
	// +optional
	Attributes map[string][]string `json:"attributes,omitempty"`

	// requiredActions lists the actions the user must perform on next login
	// (e.g. UPDATE_PASSWORD, VERIFY_EMAIL).
	// +optional
	RequiredActions []string `json:"requiredActions,omitempty"`
}

// UserStatus defines the observed state of User.
type UserStatus struct {
	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// keycloakID is the UUID assigned by Keycloak for this user.
	// It is stored here so that update and delete operations can reference the
	// user directly without an additional lookup by username.
	// +optional
	KeycloakID *string `json:"keycloakID,omitempty"`

	// conditions represent the current state of the User resource.
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
// +kubebuilder:printcolumn:name="Username",type=string,JSONPath=`.spec.username`
// +kubebuilder:printcolumn:name="Email",type=string,JSONPath=`.spec.email`
// +kubebuilder:printcolumn:name="Enabled",type=boolean,JSONPath=`.spec.enabled`
// +kubebuilder:printcolumn:name="KeycloakID",type=string,JSONPath=`.status.keycloakID`

// User is the Schema for the users API
type User struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of User
	// +required
	Spec UserSpec `json:"spec"`

	// status defines the observed state of User
	// +optional
	Status UserStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// UserList contains a list of User
type UserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []User `json:"items"`
}

func init() {
	SchemeBuilder.Register(&User{}, &UserList{})
}

// GetKeycloakID returns the Keycloak-assigned UUID stored in the status.
func (u *User) GetKeycloakID() *string { return u.Status.KeycloakID }

// SetKeycloakID stores a Keycloak-assigned UUID in the status.
func (u *User) SetKeycloakID(id *string) { u.Status.KeycloakID = id }
