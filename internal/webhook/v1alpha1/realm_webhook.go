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

	v1alpha1 "github.com/kubehippie/keycloak-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// nolint:unused
// log is for logging in this package.
var realmlog = logf.Log.WithName("realm-resource")

// SetupRealmWebhookWithManager registers the webhook for Realm in the manager.
func SetupRealmWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &v1alpha1.Realm{}).
		WithValidator(&RealmCustomValidator{}).
		WithDefaulter(&RealmCustomDefaulter{}).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-keycloak-operator-webhippie-de-v1alpha1-realm,mutating=true,failurePolicy=fail,sideEffects=None,groups=keycloak-operator.webhippie.de,resources=realms,verbs=create;update,versions=v1alpha1,name=mrealm-v1alpha1.kb.io,admissionReviewVersions=v1

// RealmCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind Realm when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type RealmCustomDefaulter struct{}

var _ admission.Defaulter[*v1alpha1.Realm] = &RealmCustomDefaulter{}

// Default implements admission.Defaulter so a webhook will be registered for the Kind Realm.
func (d *RealmCustomDefaulter) Default(_ context.Context, realm *v1alpha1.Realm) error {
	realmlog.Info("Defaulting for Realm", "name", realm.GetName())
	return nil
}

// +kubebuilder:webhook:path=/validate-keycloak-operator-webhippie-de-v1alpha1-realm,mutating=false,failurePolicy=fail,sideEffects=None,groups=keycloak-operator.webhippie.de,resources=realms,verbs=create;update,versions=v1alpha1,name=vrealm-v1alpha1.kb.io,admissionReviewVersions=v1

// RealmCustomValidator struct is responsible for validating the Realm resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type RealmCustomValidator struct{}

var _ admission.Validator[*v1alpha1.Realm] = &RealmCustomValidator{}

// ValidateCreate implements admission.Validator so a webhook will be registered for the type Realm.
func (v *RealmCustomValidator) ValidateCreate(_ context.Context, realm *v1alpha1.Realm) (admission.Warnings, error) {
	realmlog.Info("Validation for Realm upon creation", "name", realm.GetName())

	if errs := validateRealm(realm); len(errs) > 0 {
		return nil, errs.ToAggregate()
	}

	return nil, nil
}

// ValidateUpdate implements admission.Validator so a webhook will be registered for the type Realm.
func (v *RealmCustomValidator) ValidateUpdate(_ context.Context, oldRealm, realm *v1alpha1.Realm) (admission.Warnings, error) {
	realmlog.Info("Validation for Realm upon update", "name", realm.GetName())

	var allErrs field.ErrorList

	if errs := validateRealm(realm); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if oldRealm.Spec.Name != realm.Spec.Name {
		allErrs = append(allErrs, field.Forbidden(
			field.NewPath("spec", "name"),
			"name is immutable and cannot be changed after creation",
		))
	}

	if len(allErrs) > 0 {
		return nil, allErrs.ToAggregate()
	}

	return nil, nil
}

// ValidateDelete implements admission.Validator so a webhook will be registered for the type Realm.
func (v *RealmCustomValidator) ValidateDelete(_ context.Context, realm *v1alpha1.Realm) (admission.Warnings, error) {
	realmlog.Info("Validation for Realm upon deletion", "name", realm.GetName())
	return nil, nil
}

func validateRealm(realm *v1alpha1.Realm) field.ErrorList {
	var errs field.ErrorList

	if realm.Spec.KeycloakRef == nil {
		errs = append(errs, field.Required(
			field.NewPath("spec", "keycloakRef"),
			"keycloakRef is required",
		))
	} else if realm.Spec.KeycloakRef.Name == "" {
		errs = append(errs, field.Required(
			field.NewPath("spec", "keycloakRef", "name"),
			"keycloakRef.name is required",
		))
	}

	if realm.Spec.Name == "" {
		errs = append(errs, field.Required(
			field.NewPath("spec", "name"),
			"name is required",
		))
	}

	return errs
}
