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
	"strings"

	v1alpha1 "github.com/kubehippie/keycloak-operator/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// nolint:unused
// log is for logging in this package.
var rolelog = logf.Log.WithName("role-resource")

// SetupRoleWebhookWithManager registers the webhook for Role in the manager.
func SetupRoleWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &v1alpha1.Role{}).
		WithValidator(&RoleCustomValidator{}).
		WithDefaulter(&RoleCustomDefaulter{}).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-keycloak-operator-webhippie-de-v1alpha1-role,mutating=true,failurePolicy=fail,sideEffects=None,groups=keycloak-operator.webhippie.de,resources=roles,verbs=create;update,versions=v1alpha1,name=mrole-v1alpha1.kb.io,admissionReviewVersions=v1

// RoleCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind Role when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type RoleCustomDefaulter struct{}

var _ admission.Defaulter[*v1alpha1.Role] = &RoleCustomDefaulter{}

// Default implements admission.Defaulter so a webhook will be registered for the Kind Role.
func (d *RoleCustomDefaulter) Default(_ context.Context, role *v1alpha1.Role) error {
	_ = role
	return nil
}

// +kubebuilder:webhook:path=/validate-keycloak-operator-webhippie-de-v1alpha1-role,mutating=false,failurePolicy=fail,sideEffects=None,groups=keycloak-operator.webhippie.de,resources=roles,verbs=create;update,versions=v1alpha1,name=vrole-v1alpha1.kb.io,admissionReviewVersions=v1

// RoleCustomValidator struct is responsible for validating the Role resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type RoleCustomValidator struct{}

var _ admission.Validator[*v1alpha1.Role] = &RoleCustomValidator{}

// ValidateCreate implements admission.Validator so a webhook will be registered for the type Role.
func (v *RoleCustomValidator) ValidateCreate(_ context.Context, role *v1alpha1.Role) (admission.Warnings, error) {
	rolelog.Info("Validation for Role upon creation", "name", role.GetName())

	return nil, v.validate(role)
}

// ValidateUpdate implements admission.Validator so a webhook will be registered for the type Role.
func (v *RoleCustomValidator) ValidateUpdate(_ context.Context, oldRole, role *v1alpha1.Role) (admission.Warnings, error) {
	rolelog.Info("Validation for Role upon update", "name", role.GetName())

	if err := v.validate(role); err != nil {
		return nil, err
	}

	if oldRole.Spec.Name != role.Spec.Name {
		return nil, fmt.Errorf("spec.name is immutable and cannot be changed after creation")
	}

	return nil, nil
}

// ValidateDelete implements admission.Validator so a webhook will be registered for the type Role.
func (v *RoleCustomValidator) ValidateDelete(_ context.Context, role *v1alpha1.Role) (admission.Warnings, error) {
	_ = role
	return nil, nil
}

func (v *RoleCustomValidator) validate(role *v1alpha1.Role) error {
	if role.Spec.RealmRef == nil || strings.TrimSpace(role.Spec.RealmRef.Name) == "" {
		return fmt.Errorf("spec.realmRef.name must be set")
	}

	if strings.TrimSpace(role.Spec.Name) == "" {
		return fmt.Errorf("spec.name must be set")
	}

	return nil
}
