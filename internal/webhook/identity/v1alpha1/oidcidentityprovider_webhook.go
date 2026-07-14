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

	"github.com/kubehippie/keycloak-operator/api/identity/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// nolint:unused
// log is for logging in this package.
var oidcidentityproviderlog = logf.Log.WithName("oidcidentityprovider-resource")

// SetupOIDCIdentityProviderWebhookWithManager registers the webhook for OIDCIdentityProvider in the manager.
func SetupOIDCIdentityProviderWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &v1alpha1.OIDCIdentityProvider{}).
		WithValidator(&OIDCIdentityProviderCustomValidator{}).
		WithDefaulter(&OIDCIdentityProviderCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-identity-keycloak-operator-webhippie-de-v1alpha1-oidcidentityprovider,mutating=true,failurePolicy=fail,sideEffects=None,groups=identity.keycloak-operator.webhippie.de,resources=oidcidentityproviders,verbs=create;update,versions=v1alpha1,name=moidcidentityprovider-v1alpha1.kb.io,admissionReviewVersions=v1

// OIDCIdentityProviderCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind OIDCIdentityProvider when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type OIDCIdentityProviderCustomDefaulter struct{}

var _ admission.Defaulter[*v1alpha1.OIDCIdentityProvider] = &OIDCIdentityProviderCustomDefaulter{}

// Default implements admission.Defaulter so a webhook will be registered for the Kind OIDCIdentityProvider.
func (d *OIDCIdentityProviderCustomDefaulter) Default(_ context.Context, provider *v1alpha1.OIDCIdentityProvider) error {
	oidcidentityproviderlog.Info("Defaulting for OIDCIdentityProvider", "name", provider.GetName())

	if provider.Spec.DefaultScopes == nil {
		scopes := "openid"
		provider.Spec.DefaultScopes = &scopes
	}

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: If you want to customise the 'path', use the flags '--defaulting-path' or '--validation-path'.
// +kubebuilder:webhook:path=/validate-identity-keycloak-operator-webhippie-de-v1alpha1-oidcidentityprovider,mutating=false,failurePolicy=fail,sideEffects=None,groups=identity.keycloak-operator.webhippie.de,resources=oidcidentityproviders,verbs=create;update,versions=v1alpha1,name=voidcidentityprovider-v1alpha1.kb.io,admissionReviewVersions=v1

// OIDCIdentityProviderCustomValidator struct is responsible for validating the OIDCIdentityProvider resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type OIDCIdentityProviderCustomValidator struct{}

var _ admission.Validator[*v1alpha1.OIDCIdentityProvider] = &OIDCIdentityProviderCustomValidator{}

// ValidateCreate implements admission.Validator so a webhook will be registered for the type OIDCIdentityProvider.
func (v *OIDCIdentityProviderCustomValidator) ValidateCreate(_ context.Context, provider *v1alpha1.OIDCIdentityProvider) (admission.Warnings, error) {
	oidcidentityproviderlog.Info("Validation for OIDCIdentityProvider upon creation", "name", provider.GetName())

	if errs := validateOIDCIdentityProvider(provider); len(errs) > 0 {
		return nil, errs.ToAggregate()
	}

	return nil, nil
}

// ValidateUpdate implements admission.Validator so a webhook will be registered for the type OIDCIdentityProvider.
func (v *OIDCIdentityProviderCustomValidator) ValidateUpdate(_ context.Context, oldProvider, provider *v1alpha1.OIDCIdentityProvider) (admission.Warnings, error) {
	oidcidentityproviderlog.Info("Validation for OIDCIdentityProvider upon update", "name", provider.GetName())

	var allErrs field.ErrorList

	if errs := validateOIDCIdentityProvider(provider); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if oldProvider.Spec.Alias != provider.Spec.Alias {
		allErrs = append(allErrs, field.Forbidden(
			field.NewPath("spec", "alias"),
			"alias is immutable and cannot be changed after creation",
		))
	}

	if len(allErrs) > 0 {
		return nil, allErrs.ToAggregate()
	}

	return nil, nil
}

// ValidateDelete implements admission.Validator so a webhook will be registered for the type OIDCIdentityProvider.
func (v *OIDCIdentityProviderCustomValidator) ValidateDelete(_ context.Context, provider *v1alpha1.OIDCIdentityProvider) (admission.Warnings, error) {
	oidcidentityproviderlog.Info("Validation for OIDCIdentityProvider upon deletion", "name", provider.GetName())
	return nil, nil
}

func validateOIDCIdentityProvider(provider *v1alpha1.OIDCIdentityProvider) field.ErrorList {
	var errs field.ErrorList

	if provider.Spec.RealmRef == nil {
		errs = append(errs, field.Required(
			field.NewPath("spec", "realmRef"),
			"realmRef is required",
		))
	} else if provider.Spec.RealmRef.Name == "" {
		errs = append(errs, field.Required(
			field.NewPath("spec", "realmRef", "name"),
			"realmRef.name is required",
		))
	}

	if provider.Spec.Alias == "" {
		errs = append(errs, field.Required(
			field.NewPath("spec", "alias"),
			"alias is required",
		))
	}

	if provider.Spec.AuthorizationURL == "" {
		errs = append(errs, field.Required(
			field.NewPath("spec", "authorizationURL"),
			"authorizationURL is required",
		))
	}

	if provider.Spec.TokenURL == "" {
		errs = append(errs, field.Required(
			field.NewPath("spec", "tokenURL"),
			"tokenURL is required",
		))
	}

	if provider.Spec.ClientID == "" {
		errs = append(errs, field.Required(
			field.NewPath("spec", "clientID"),
			"clientID is required",
		))
	}

	if provider.Spec.ClientSecret == nil {
		errs = append(errs, field.Required(
			field.NewPath("spec", "clientSecret"),
			"clientSecret is required",
		))
	} else if provider.Spec.ClientSecret.Value == "" && provider.Spec.ClientSecret.SecretKeyRef == nil {
		errs = append(errs, field.Required(
			field.NewPath("spec", "clientSecret"),
			"either clientSecret.value or clientSecret.secretKeyRef must be set",
		))
	}

	return errs
}
