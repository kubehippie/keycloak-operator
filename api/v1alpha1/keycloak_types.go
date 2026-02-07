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

// KeycloakSpec defines the desired state of Keycloak
type KeycloakSpec struct {
	// url of keycloak instance.
	URL string `json:"url"`

	// secret references a secret which contains admin credentials.
	Secret *common.SecretKeySelector `json:"secret"`

	// caCert defines the root certificate authority that API clients use when
	// verifying server certificates.
	// +optional
	CACert *common.SourceRef `json:"caCert,omitempty"`

	// insecureSkipVerify controls whether API client verifies the servers
	// certificate chain and hostname. If insecureSkipVerify is true, API client
	// accepts any certificate presented by the server and any hostname in that
	// certificate.
	InsecureSkipVerify bool `json:"insecureSkipVerify,omitempty"`
}

// KeycloakStatus defines the observed state of Keycloak.
type KeycloakStatus struct {
	// connected shows if the Keycloak instance is connected successfully.
	Connected bool `json:"connected"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Connected",type="boolean",JSONPath=".status.connected",description="Connected to Keycloak"

// Keycloak is the Schema for the keycloaks API
type Keycloak struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of Keycloak
	// +required
	Spec KeycloakSpec `json:"spec"`

	// status defines the observed state of Keycloak
	// +kubebuilder:default={connected:false}
	Status KeycloakStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// KeycloakList contains a list of Keycloak
type KeycloakList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Keycloak `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Keycloak{}, &KeycloakList{})
}
