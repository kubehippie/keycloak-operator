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
var attributeimportermapperlog = logf.Log.WithName("attributeimportermapper-resource")

// SetupAttributeImporterMapperWebhookWithManager registers the webhook for AttributeImporterMapper in the manager.
func SetupAttributeImporterMapperWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &v1alpha1.AttributeImporterMapper{}).
		WithValidator(&AttributeImporterMapperCustomValidator{}).
		WithDefaulter(&AttributeImporterMapperCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-identity-keycloak-operator-webhippie-de-v1alpha1-attributeimportermapper,mutating=true,failurePolicy=fail,sideEffects=None,groups=identity.keycloak-operator.webhippie.de,resources=attributeimportermappers,verbs=create;update,versions=v1alpha1,name=mattributeimportermapper-v1alpha1.kb.io,admissionReviewVersions=v1

// AttributeImporterMapperCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind AttributeImporterMapper when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type AttributeImporterMapperCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ admission.Defaulter[*v1alpha1.AttributeImporterMapper] = &AttributeImporterMapperCustomDefaulter{}

// Default implements admission.Defaulter so a webhook will be registered for the Kind AttributeImporterMapper.
func (d *AttributeImporterMapperCustomDefaulter) Default(_ context.Context, mapper *v1alpha1.AttributeImporterMapper) error {
	attributeimportermapperlog.Info("Defaulting for AttributeImporterMapper", "name", mapper.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: If you want to customise the 'path', use the flags '--defaulting-path' or '--validation-path'.
// +kubebuilder:webhook:path=/validate-identity-keycloak-operator-webhippie-de-v1alpha1-attributeimportermapper,mutating=false,failurePolicy=fail,sideEffects=None,groups=identity.keycloak-operator.webhippie.de,resources=attributeimportermappers,verbs=create;update,versions=v1alpha1,name=vattributeimportermapper-v1alpha1.kb.io,admissionReviewVersions=v1

// AttributeImporterMapperCustomValidator struct is responsible for validating the AttributeImporterMapper resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type AttributeImporterMapperCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ admission.Validator[*v1alpha1.AttributeImporterMapper] = &AttributeImporterMapperCustomValidator{}

// ValidateCreate implements admission.Validator so a webhook will be registered for the type AttributeImporterMapper.
func (v *AttributeImporterMapperCustomValidator) ValidateCreate(_ context.Context, mapper *v1alpha1.AttributeImporterMapper) (admission.Warnings, error) {
	attributeimportermapperlog.Info("Validation for AttributeImporterMapper upon creation", "name", mapper.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements admission.Validator so a webhook will be registered for the type AttributeImporterMapper.
func (v *AttributeImporterMapperCustomValidator) ValidateUpdate(_ context.Context, oldMapper, mapper *v1alpha1.AttributeImporterMapper) (admission.Warnings, error) {
	_ = oldMapper
	attributeimportermapperlog.Info("Validation for AttributeImporterMapper upon update", "name", mapper.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements admission.Validator so a webhook will be registered for the type AttributeImporterMapper.
func (v *AttributeImporterMapperCustomValidator) ValidateDelete(_ context.Context, mapper *v1alpha1.AttributeImporterMapper) (admission.Warnings, error) {
	attributeimportermapperlog.Info("Validation for AttributeImporterMapper upon deletion", "name", mapper.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
