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

	"github.com/kubehippie/keycloak-operator/api/openid/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// nolint:unused
// log is for logging in this package.
var openidclientlog = logf.Log.WithName("openidclient-resource")

const (
	accessTypeConfidential = "CONFIDENTIAL"
	accessTypePublic       = "PUBLIC"
	accessTypeBearerOnly   = "BEARER-ONLY"
)

// SetupOpenIDClientWebhookWithManager registers the webhook for OpenIDClient in the manager.
func SetupOpenIDClientWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &v1alpha1.OpenIDClient{}).
		WithValidator(&OpenIDClientCustomValidator{}).
		WithDefaulter(&OpenIDClientCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-openid-keycloak-operator-webhippie-de-v1alpha1-openidclient,mutating=true,failurePolicy=fail,sideEffects=None,groups=openid.keycloak-operator.webhippie.de,resources=openidclients,verbs=create;update,versions=v1alpha1,name=mopenidclient-v1alpha1.kb.io,admissionReviewVersions=v1

// OpenIDClientCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind OpenIDClient when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type OpenIDClientCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ admission.Defaulter[*v1alpha1.OpenIDClient] = &OpenIDClientCustomDefaulter{}

// Default implements admission.Defaulter so a webhook will be registered for the Kind OpenIDClient.
func (d *OpenIDClientCustomDefaulter) Default(_ context.Context, client *v1alpha1.OpenIDClient) error {
	openidclientlog.Info("Defaulting for OpenIDClient", "name", client.GetName())

	if client.Spec.AccessType == "" {
		client.Spec.AccessType = accessTypeConfidential
	}
	if client.Spec.Enabled == nil {
		client.Spec.Enabled = boolPtr(true)
	}
	if client.Spec.StandardFlowEnabled == nil {
		client.Spec.StandardFlowEnabled = boolPtr(true)
	}
	if client.Spec.ImplicitFlowEnabled == nil {
		client.Spec.ImplicitFlowEnabled = boolPtr(false)
	}
	if client.Spec.DirectAccessGrantsEnabled == nil {
		client.Spec.DirectAccessGrantsEnabled = boolPtr(true)
	}
	if client.Spec.ServiceAccountsEnabled == nil {
		client.Spec.ServiceAccountsEnabled = boolPtr(false)
	}
	if client.Spec.FrontChannelLogoutEnabled == nil {
		client.Spec.FrontChannelLogoutEnabled = boolPtr(false)
	}
	if client.Spec.FullScopeAllowed == nil {
		client.Spec.FullScopeAllowed = boolPtr(true)
	}
	if client.Spec.ClientAuthenticatorType == nil {
		client.Spec.ClientAuthenticatorType = stringPtr("client-secret")
	}

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: If you want to customise the 'path', use the flags '--defaulting-path' or '--validation-path'.
// +kubebuilder:webhook:path=/validate-openid-keycloak-operator-webhippie-de-v1alpha1-openidclient,mutating=false,failurePolicy=fail,sideEffects=None,groups=openid.keycloak-operator.webhippie.de,resources=openidclients,verbs=create;update,versions=v1alpha1,name=vopenidclient-v1alpha1.kb.io,admissionReviewVersions=v1

// OpenIDClientCustomValidator struct is responsible for validating the OpenIDClient resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type OpenIDClientCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ admission.Validator[*v1alpha1.OpenIDClient] = &OpenIDClientCustomValidator{}

// ValidateCreate implements admission.Validator so a webhook will be registered for the type OpenIDClient.
func (v *OpenIDClientCustomValidator) ValidateCreate(_ context.Context, client *v1alpha1.OpenIDClient) (admission.Warnings, error) {
	openidclientlog.Info("Validation for OpenIDClient upon creation", "name", client.GetName())

	return nil, validateOpenIDClient(client)
}

// ValidateUpdate implements admission.Validator so a webhook will be registered for the type OpenIDClient.
func (v *OpenIDClientCustomValidator) ValidateUpdate(_ context.Context, oldClient, client *v1alpha1.OpenIDClient) (admission.Warnings, error) {
	openidclientlog.Info("Validation for OpenIDClient upon update", "name", client.GetName())

	if err := validateOpenIDClient(client); err != nil {
		return nil, err
	}

	if oldClient.Spec.ClientID != client.Spec.ClientID {
		return nil, field.Invalid(field.NewPath("spec").Child("clientID"), client.Spec.ClientID, "clientID is immutable")
	}

	return nil, nil
}

// ValidateDelete implements admission.Validator so a webhook will be registered for the type OpenIDClient.
func (v *OpenIDClientCustomValidator) ValidateDelete(_ context.Context, client *v1alpha1.OpenIDClient) (admission.Warnings, error) {
	openidclientlog.Info("Validation for OpenIDClient upon deletion", "name", client.GetName())

	return nil, nil
}

func validateOpenIDClient(client *v1alpha1.OpenIDClient) error {
	if client.Spec.RealmRef == nil {
		return field.Required(field.NewPath("spec").Child("realmRef"), "realmRef is required")
	}
	if client.Spec.ClientID == "" {
		return field.Required(field.NewPath("spec").Child("clientID"), "clientID is required")
	}
	switch client.Spec.AccessType {
	case "", accessTypeConfidential, accessTypePublic, accessTypeBearerOnly:
	default:
		return field.Invalid(field.NewPath("spec").Child("accessType"), client.Spec.AccessType, "must be one of CONFIDENTIAL, PUBLIC, BEARER-ONLY")
	}
	if client.Spec.AccessType == accessTypeConfidential && client.Spec.ClientSecret == nil {
		return field.Required(field.NewPath("spec").Child("clientSecret"), "clientSecret is required when accessType is CONFIDENTIAL")
	}
	return nil
}
