package common

// KeycloakRef is a reference to a Keycloak instance.
type KeycloakRef struct {
	// Kind specifies the kind of the Keycloak resource.
	// +kubebuilder:validation:Enum=Keycloak
	// +kubebuilder:default=Keycloak
	// +optional
	Kind string `json:"kind,omitempty"`

	// Name specifies the name of the Keycloak resource.
	// +required
	Name string `json:"name,omitempty"`

	// Namespace specifies the namespace of the Keycloak resource.
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// RealmRef is a reference to a Realm resource.
type RealmRef struct {
	// Kind specifies the kind of the Realm resource.
	// +kubebuilder:validation:Enum=Realm
	// +kubebuilder:default=Realm
	// +optional
	Kind string `json:"kind,omitempty"`

	// Name specifies the name of the Realm resource.
	// +required
	Name string `json:"name,omitempty"`

	// Namespace specifies the namespace of the Realm resource.
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// ConfigMapKeySelector selects a key of a ConfigMap.
// +kubebuilder:object:generate=true
type ConfigMapKeySelector struct {
	// Name is the name of the ConfigMap.
	// +required
	Name string `json:"name"`

	// Key is the key in the ConfigMap to select from.
	// +required
	Key string `json:"key"`

	// Namespace of the ConfigMap.
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// SecretKeySelector selects a key of a Secret.
// +kubebuilder:object:generate=true
type SecretKeySelector struct {
	// Name is the name of the Secret.
	// +required
	Name string `json:"name"`

	// Key is the key of the Secret to select from.
	// +required
	Key string `json:"key"`

	// Namespace of the Secret.
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// SecretKeyRefOrVal holds either an inline value or a reference to a Secret key.
// Exactly one of value or secretKeyRef must be set.
// +kubebuilder:object:generate=true
type SecretKeyRefOrVal struct {
	// SecretKeyRef selects a key of a Kubernetes Secret.
	// +optional
	SecretKeyRef *SecretKeySelector `json:"secretKeyRef,omitempty"`

	// Value directly specifies the value as a plain string.
	// +optional
	Value string `json:"value,omitempty"`
}

// ConfigMapRefOrVal holds either an inline value or a reference to a ConfigMap key.
// Exactly one of value or configMapKeyRef must be set.
// +kubebuilder:object:generate=true
type ConfigMapRefOrVal struct {
	// ConfigMapKeyRef selects a key of a Kubernetes ConfigMap.
	// +optional
	ConfigMapKeyRef *ConfigMapKeySelector `json:"configMapKeyRef,omitempty"`

	// Value directly specifies the value as a plain string.
	// +optional
	Value string `json:"value,omitempty"`
}

// SourceRef is a reference to a key in a ConfigMap or a Secret.
// +kubebuilder:object:generate=true
type SourceRef struct {
	// ConfigMapKeyRef selects a key of a ConfigMap.
	// +optional
	ConfigMapKeyRef *ConfigMapKeySelector `json:"configMapKeyRef,omitempty"`

	// SecretKeyRef selects a key of a Secret.
	// +optional
	SecretKeyRef *SecretKeySelector `json:"secretKeyRef,omitempty"`
}

// SourceRefOrVal holds either an inline value or a reference to a ConfigMap or Secret key.
// Exactly one of value, configMapKeyRef, or secretKeyRef must be set.
// +kubebuilder:object:generate=true
type SourceRefOrVal struct {
	// ConfigMapKeyRef selects a key of a ConfigMap.
	// +optional
	ConfigMapKeyRef *ConfigMapKeySelector `json:"configMapKeyRef,omitempty"`

	// SecretKeyRef selects a key of a Secret.
	// +optional
	SecretKeyRef *SecretKeySelector `json:"secretKeyRef,omitempty"`

	// Value directly specifies the value as a plain string.
	// +optional
	Value string `json:"value,omitempty"`
}
