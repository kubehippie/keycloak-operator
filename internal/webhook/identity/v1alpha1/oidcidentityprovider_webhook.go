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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	identityv1alpha1 "github.com/kubehippie/keycloak-operator/api/identity/v1alpha1"
)

// nolint:unused
// log is for logging in this package.
var oidcidentityproviderlog = logf.Log.WithName("oidcidentityprovider-resource")

// SetupOIDCIdentityProviderWebhookWithManager registers the webhook for OIDCIdentityProvider in the manager.
func SetupOIDCIdentityProviderWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&identityv1alpha1.OIDCIdentityProvider{}).
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
type OIDCIdentityProviderCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &OIDCIdentityProviderCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind OIDCIdentityProvider.
func (d *OIDCIdentityProviderCustomDefaulter) Default(_ context.Context, obj runtime.Object) error {
	oidcidentityprovider, ok := obj.(*identityv1alpha1.OIDCIdentityProvider)

	if !ok {
		return fmt.Errorf("expected an OIDCIdentityProvider object but got %T", obj)
	}
	oidcidentityproviderlog.Info("Defaulting for OIDCIdentityProvider", "name", oidcidentityprovider.GetName())

	// TODO(user): fill in your defaulting logic.

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
type OIDCIdentityProviderCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &OIDCIdentityProviderCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type OIDCIdentityProvider.
func (v *OIDCIdentityProviderCustomValidator) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	oidcidentityprovider, ok := obj.(*identityv1alpha1.OIDCIdentityProvider)
	if !ok {
		return nil, fmt.Errorf("expected a OIDCIdentityProvider object but got %T", obj)
	}
	oidcidentityproviderlog.Info("Validation for OIDCIdentityProvider upon creation", "name", oidcidentityprovider.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type OIDCIdentityProvider.
func (v *OIDCIdentityProviderCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	oidcidentityprovider, ok := newObj.(*identityv1alpha1.OIDCIdentityProvider)
	if !ok {
		return nil, fmt.Errorf("expected a OIDCIdentityProvider object for the newObj but got %T", newObj)
	}
	oidcidentityproviderlog.Info("Validation for OIDCIdentityProvider upon update", "name", oidcidentityprovider.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type OIDCIdentityProvider.
func (v *OIDCIdentityProviderCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	oidcidentityprovider, ok := obj.(*identityv1alpha1.OIDCIdentityProvider)
	if !ok {
		return nil, fmt.Errorf("expected a OIDCIdentityProvider object but got %T", obj)
	}
	oidcidentityproviderlog.Info("Validation for OIDCIdentityProvider upon deletion", "name", oidcidentityprovider.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
