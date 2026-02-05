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
var oidcadvancedgroupmapperlog = logf.Log.WithName("oidcadvancedgroupmapper-resource")

// SetupOIDCAdvancedGroupMapperWebhookWithManager registers the webhook for OIDCAdvancedGroupMapper in the manager.
func SetupOIDCAdvancedGroupMapperWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&identityv1alpha1.OIDCAdvancedGroupMapper{}).
		WithValidator(&OIDCAdvancedGroupMapperCustomValidator{}).
		WithDefaulter(&OIDCAdvancedGroupMapperCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-identity-keycloak-operator-webhippie-de-v1alpha1-oidcadvancedgroupmapper,mutating=true,failurePolicy=fail,sideEffects=None,groups=identity.keycloak-operator.webhippie.de,resources=oidcadvancedgroupmappers,verbs=create;update,versions=v1alpha1,name=moidcadvancedgroupmapper-v1alpha1.kb.io,admissionReviewVersions=v1

// OIDCAdvancedGroupMapperCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind OIDCAdvancedGroupMapper when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type OIDCAdvancedGroupMapperCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &OIDCAdvancedGroupMapperCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind OIDCAdvancedGroupMapper.
func (d *OIDCAdvancedGroupMapperCustomDefaulter) Default(_ context.Context, obj runtime.Object) error {
	oidcadvancedgroupmapper, ok := obj.(*identityv1alpha1.OIDCAdvancedGroupMapper)

	if !ok {
		return fmt.Errorf("expected an OIDCAdvancedGroupMapper object but got %T", obj)
	}
	oidcadvancedgroupmapperlog.Info("Defaulting for OIDCAdvancedGroupMapper", "name", oidcadvancedgroupmapper.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: If you want to customise the 'path', use the flags '--defaulting-path' or '--validation-path'.
// +kubebuilder:webhook:path=/validate-identity-keycloak-operator-webhippie-de-v1alpha1-oidcadvancedgroupmapper,mutating=false,failurePolicy=fail,sideEffects=None,groups=identity.keycloak-operator.webhippie.de,resources=oidcadvancedgroupmappers,verbs=create;update,versions=v1alpha1,name=voidcadvancedgroupmapper-v1alpha1.kb.io,admissionReviewVersions=v1

// OIDCAdvancedGroupMapperCustomValidator struct is responsible for validating the OIDCAdvancedGroupMapper resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type OIDCAdvancedGroupMapperCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &OIDCAdvancedGroupMapperCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type OIDCAdvancedGroupMapper.
func (v *OIDCAdvancedGroupMapperCustomValidator) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	oidcadvancedgroupmapper, ok := obj.(*identityv1alpha1.OIDCAdvancedGroupMapper)
	if !ok {
		return nil, fmt.Errorf("expected a OIDCAdvancedGroupMapper object but got %T", obj)
	}
	oidcadvancedgroupmapperlog.Info("Validation for OIDCAdvancedGroupMapper upon creation", "name", oidcadvancedgroupmapper.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type OIDCAdvancedGroupMapper.
func (v *OIDCAdvancedGroupMapperCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	oidcadvancedgroupmapper, ok := newObj.(*identityv1alpha1.OIDCAdvancedGroupMapper)
	if !ok {
		return nil, fmt.Errorf("expected a OIDCAdvancedGroupMapper object for the newObj but got %T", newObj)
	}
	oidcadvancedgroupmapperlog.Info("Validation for OIDCAdvancedGroupMapper upon update", "name", oidcadvancedgroupmapper.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type OIDCAdvancedGroupMapper.
func (v *OIDCAdvancedGroupMapperCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	oidcadvancedgroupmapper, ok := obj.(*identityv1alpha1.OIDCAdvancedGroupMapper)
	if !ok {
		return nil, fmt.Errorf("expected a OIDCAdvancedGroupMapper object but got %T", obj)
	}
	oidcadvancedgroupmapperlog.Info("Validation for OIDCAdvancedGroupMapper upon deletion", "name", oidcadvancedgroupmapper.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
