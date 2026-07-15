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
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/kubehippie/keycloak-operator/api/common"
	v1alpha1 "github.com/kubehippie/keycloak-operator/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// SetupKeycloakWebhookWithManager registers the webhook for Keycloak in the manager.
func SetupKeycloakWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &v1alpha1.Keycloak{}).
		WithValidator(&KeycloakCustomValidator{}).
		WithDefaulter(&KeycloakCustomDefaulter{}).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-keycloak-operator-webhippie-de-v1alpha1-keycloak,mutating=true,failurePolicy=fail,sideEffects=None,groups=keycloak-operator.webhippie.de,resources=keycloaks,verbs=create;update,versions=v1alpha1,name=mkeycloak-v1alpha1.kb.io,admissionReviewVersions=v1

// KeycloakCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind Keycloak when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type KeycloakCustomDefaulter struct{}

var _ admission.Defaulter[*v1alpha1.Keycloak] = &KeycloakCustomDefaulter{}

// Default implements admission.Defaulter so a webhook will be registered for the Kind Keycloak.
func (d *KeycloakCustomDefaulter) Default(_ context.Context, keycloak *v1alpha1.Keycloak) error {
	_ = keycloak
	return nil
}

// +kubebuilder:webhook:path=/validate-keycloak-operator-webhippie-de-v1alpha1-keycloak,mutating=false,failurePolicy=fail,sideEffects=None,groups=keycloak-operator.webhippie.de,resources=keycloaks,verbs=create;update,versions=v1alpha1,name=vkeycloak-v1alpha1.kb.io,admissionReviewVersions=v1

// KeycloakCustomValidator struct is responsible for validating the Keycloak resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type KeycloakCustomValidator struct{}

var _ admission.Validator[*v1alpha1.Keycloak] = &KeycloakCustomValidator{}

// ValidateCreate implements admission.Validator so a webhook will be registered for the type Keycloak.
func (v *KeycloakCustomValidator) ValidateCreate(_ context.Context, keycloak *v1alpha1.Keycloak) (admission.Warnings, error) {
	return nil, v.validate(keycloak)
}

// ValidateUpdate implements admission.Validator so a webhook will be registered for the type Keycloak.
func (v *KeycloakCustomValidator) ValidateUpdate(_ context.Context, oldKeycloak, keycloak *v1alpha1.Keycloak) (admission.Warnings, error) {
	_ = oldKeycloak
	return nil, v.validate(keycloak)
}

// ValidateDelete implements admission.Validator so a webhook will be registered for the type Keycloak.
func (v *KeycloakCustomValidator) ValidateDelete(_ context.Context, keycloak *v1alpha1.Keycloak) (admission.Warnings, error) {
	_ = keycloak
	return nil, nil
}

func (v *KeycloakCustomValidator) validate(keycloak *v1alpha1.Keycloak) error {
	spec := keycloak.Spec

	parsed, err := url.ParseRequestURI(
		strings.TrimSpace(spec.URL),
	)

	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("spec.url must be a valid URL")
	}

	if strings.TrimSpace(spec.RealmName) == "" {
		return fmt.Errorf("spec.realmName must be set")
	}

	usernamePasswordSet := spec.Username != nil || spec.Password != nil
	clientSecretSet := spec.Client != nil || spec.Secret != nil

	if usernamePasswordSet && clientSecretSet {
		return fmt.Errorf("only one of spec.username/spec.password or spec.client/spec.secret must be set")
	}

	if !usernamePasswordSet && !clientSecretSet {
		return fmt.Errorf("either spec.username/spec.password or spec.client/spec.secret must be set")
	}

	if usernamePasswordSet {
		if err := v.validateSecretKeyRefOrVal("spec.username", spec.Username); err != nil {
			return err
		}

		if err := v.validateSecretKeyRefOrVal("spec.password", spec.Password); err != nil {
			return err
		}
	}

	if clientSecretSet {
		if err := v.validateSecretKeyRefOrVal("spec.client", spec.Client); err != nil {
			return err
		}

		if err := v.validateSecretKeyRefOrVal("spec.secret", spec.Secret); err != nil {
			return err
		}
	}

	return nil
}

func (v *KeycloakCustomValidator) validateSecretKeyRefOrVal(field string, ref *common.SecretKeyRefOrVal) error {
	if ref == nil {
		return fmt.Errorf("%s must be set", field)
	}

	if strings.TrimSpace(ref.Value) != "" {
		return nil
	}

	if ref.SecretKeyRef == nil {
		return fmt.Errorf("%s must set value or secretKeyRef", field)
	}

	if strings.TrimSpace(ref.SecretKeyRef.Name) == "" || strings.TrimSpace(ref.SecretKeyRef.Key) == "" {
		return fmt.Errorf("%s.secretKeyRef must set name and key", field)
	}

	return nil
}
