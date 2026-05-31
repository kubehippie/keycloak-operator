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
var grouplog = logf.Log.WithName("group-resource")

// SetupGroupWebhookWithManager registers the webhook for Group in the manager.
func SetupGroupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &v1alpha1.Group{}).
		WithValidator(&GroupCustomValidator{}).
		WithDefaulter(&GroupCustomDefaulter{}).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-keycloak-operator-webhippie-de-v1alpha1-group,mutating=true,failurePolicy=fail,sideEffects=None,groups=keycloak-operator.webhippie.de,resources=groups,verbs=create;update,versions=v1alpha1,name=mgroup-v1alpha1.kb.io,admissionReviewVersions=v1

// GroupCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind Group when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type GroupCustomDefaulter struct{}

var _ admission.Defaulter[*v1alpha1.Group] = &GroupCustomDefaulter{}

// Default implements admission.Defaulter so a webhook will be registered for the Kind Group.
func (d *GroupCustomDefaulter) Default(_ context.Context, group *v1alpha1.Group) error {
	_ = group
	return nil
}

// +kubebuilder:webhook:path=/validate-keycloak-operator-webhippie-de-v1alpha1-group,mutating=false,failurePolicy=fail,sideEffects=None,groups=keycloak-operator.webhippie.de,resources=groups,verbs=create;update,versions=v1alpha1,name=vgroup-v1alpha1.kb.io,admissionReviewVersions=v1

// GroupCustomValidator struct is responsible for validating the Group resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type GroupCustomValidator struct{}

var _ admission.Validator[*v1alpha1.Group] = &GroupCustomValidator{}

// ValidateCreate implements admission.Validator so a webhook will be registered for the type Group.
func (v *GroupCustomValidator) ValidateCreate(_ context.Context, group *v1alpha1.Group) (admission.Warnings, error) {
	grouplog.Info("Validation for Group upon creation", "name", group.GetName())

	return nil, v.validate(group)
}

// ValidateUpdate implements admission.Validator so a webhook will be registered for the type Group.
func (v *GroupCustomValidator) ValidateUpdate(_ context.Context, oldGroup, group *v1alpha1.Group) (admission.Warnings, error) {
	grouplog.Info("Validation for Group upon update", "name", group.GetName())

	if err := v.validate(group); err != nil {
		return nil, err
	}

	if oldGroup.Spec.Name != group.Spec.Name {
		return nil, fmt.Errorf("spec.name is immutable and cannot be changed after creation")
	}

	return nil, nil
}

// ValidateDelete implements admission.Validator so a webhook will be registered for the type Group.
func (v *GroupCustomValidator) ValidateDelete(_ context.Context, group *v1alpha1.Group) (admission.Warnings, error) {
	_ = group
	return nil, nil
}

func (v *GroupCustomValidator) validate(group *v1alpha1.Group) error {
	if group.Spec.RealmRef == nil || strings.TrimSpace(group.Spec.RealmRef.Name) == "" {
		return fmt.Errorf("spec.realmRef.name must be set")
	}

	if strings.TrimSpace(group.Spec.Name) == "" {
		return fmt.Errorf("spec.name must be set")
	}

	return nil
}
