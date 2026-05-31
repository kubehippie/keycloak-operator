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
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// nolint:unused
// log is for logging in this package.
var hardcodedrolemapperlog = logf.Log.WithName("hardcodedrolemapper-resource")

// SetupHardcodedRoleMapperWebhookWithManager registers the webhook for HardcodedRoleMapper in the manager.
func SetupHardcodedRoleMapperWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &v1alpha1.HardcodedRoleMapper{}).
		WithValidator(&HardcodedRoleMapperCustomValidator{}).
		WithDefaulter(&HardcodedRoleMapperCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-identity-keycloak-operator-webhippie-de-v1alpha1-hardcodedrolemapper,mutating=true,failurePolicy=fail,sideEffects=None,groups=identity.keycloak-operator.webhippie.de,resources=hardcodedrolemappers,verbs=create;update,versions=v1alpha1,name=mhardcodedrolemapper-v1alpha1.kb.io,admissionReviewVersions=v1

// HardcodedRoleMapperCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind HardcodedRoleMapper when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type HardcodedRoleMapperCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ admission.Defaulter[*v1alpha1.HardcodedRoleMapper] = &HardcodedRoleMapperCustomDefaulter{}

// Default implements admission.Defaulter so a webhook will be registered for the Kind HardcodedRoleMapper.
func (d *HardcodedRoleMapperCustomDefaulter) Default(_ context.Context, mapper *v1alpha1.HardcodedRoleMapper) error {
	hardcodedrolemapperlog.Info("Defaulting for HardcodedRoleMapper", "name", mapper.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: If you want to customise the 'path', use the flags '--defaulting-path' or '--validation-path'.
// +kubebuilder:webhook:path=/validate-identity-keycloak-operator-webhippie-de-v1alpha1-hardcodedrolemapper,mutating=false,failurePolicy=fail,sideEffects=None,groups=identity.keycloak-operator.webhippie.de,resources=hardcodedrolemappers,verbs=create;update,versions=v1alpha1,name=vhardcodedrolemapper-v1alpha1.kb.io,admissionReviewVersions=v1

// HardcodedRoleMapperCustomValidator struct is responsible for validating the HardcodedRoleMapper resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type HardcodedRoleMapperCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ admission.Validator[*v1alpha1.HardcodedRoleMapper] = &HardcodedRoleMapperCustomValidator{}

// ValidateCreate implements admission.Validator so a webhook will be registered for the type HardcodedRoleMapper.
func (v *HardcodedRoleMapperCustomValidator) ValidateCreate(_ context.Context, mapper *v1alpha1.HardcodedRoleMapper) (admission.Warnings, error) {
	hardcodedrolemapperlog.Info("Validation for HardcodedRoleMapper upon creation", "name", mapper.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements admission.Validator so a webhook will be registered for the type HardcodedRoleMapper.
func (v *HardcodedRoleMapperCustomValidator) ValidateUpdate(_ context.Context, oldMapper, mapper *v1alpha1.HardcodedRoleMapper) (admission.Warnings, error) {
	_ = oldMapper
	hardcodedrolemapperlog.Info("Validation for HardcodedRoleMapper upon update", "name", mapper.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements admission.Validator so a webhook will be registered for the type HardcodedRoleMapper.
func (v *HardcodedRoleMapperCustomValidator) ValidateDelete(_ context.Context, mapper *v1alpha1.HardcodedRoleMapper) (admission.Warnings, error) {
	hardcodedrolemapperlog.Info("Validation for HardcodedRoleMapper upon deletion", "name", mapper.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
