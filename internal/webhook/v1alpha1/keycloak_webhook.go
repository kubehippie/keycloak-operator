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

	keycloakoperatorwebhippiedev1alpha1 "github.com/kubehippie/keycloak-operator/api/v1alpha1"
)

// nolint:unused
// log is for logging in this package.
var keycloaklog = logf.Log.WithName("keycloak-resource")

// SetupKeycloakWebhookWithManager registers the webhook for Keycloak in the manager.
func SetupKeycloakWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&keycloakoperatorwebhippiedev1alpha1.Keycloak{}).
		WithValidator(&KeycloakCustomValidator{}).
		WithDefaulter(&KeycloakCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-keycloak-operator-webhippie-de-v1alpha1-keycloak,mutating=true,failurePolicy=fail,sideEffects=None,groups=keycloak-operator.webhippie.de,resources=keycloaks,verbs=create;update,versions=v1alpha1,name=mkeycloak-v1alpha1.kb.io,admissionReviewVersions=v1

// KeycloakCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind Keycloak when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type KeycloakCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &KeycloakCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind Keycloak.
func (d *KeycloakCustomDefaulter) Default(_ context.Context, obj runtime.Object) error {
	keycloak, ok := obj.(*keycloakoperatorwebhippiedev1alpha1.Keycloak)

	if !ok {
		return fmt.Errorf("expected an Keycloak object but got %T", obj)
	}
	keycloaklog.Info("Defaulting for Keycloak", "name", keycloak.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: If you want to customise the 'path', use the flags '--defaulting-path' or '--validation-path'.
// +kubebuilder:webhook:path=/validate-keycloak-operator-webhippie-de-v1alpha1-keycloak,mutating=false,failurePolicy=fail,sideEffects=None,groups=keycloak-operator.webhippie.de,resources=keycloaks,verbs=create;update,versions=v1alpha1,name=vkeycloak-v1alpha1.kb.io,admissionReviewVersions=v1

// KeycloakCustomValidator struct is responsible for validating the Keycloak resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type KeycloakCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &KeycloakCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type Keycloak.
func (v *KeycloakCustomValidator) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	keycloak, ok := obj.(*keycloakoperatorwebhippiedev1alpha1.Keycloak)
	if !ok {
		return nil, fmt.Errorf("expected a Keycloak object but got %T", obj)
	}
	keycloaklog.Info("Validation for Keycloak upon creation", "name", keycloak.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type Keycloak.
func (v *KeycloakCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	keycloak, ok := newObj.(*keycloakoperatorwebhippiedev1alpha1.Keycloak)
	if !ok {
		return nil, fmt.Errorf("expected a Keycloak object for the newObj but got %T", newObj)
	}
	keycloaklog.Info("Validation for Keycloak upon update", "name", keycloak.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type Keycloak.
func (v *KeycloakCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	keycloak, ok := obj.(*keycloakoperatorwebhippiedev1alpha1.Keycloak)
	if !ok {
		return nil, fmt.Errorf("expected a Keycloak object but got %T", obj)
	}
	keycloaklog.Info("Validation for Keycloak upon deletion", "name", keycloak.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
