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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	keycloakoperatorwebhippiedev1alpha1 "github.com/kubehippie/keycloak-operator/api/v1alpha1"
)

// nolint:unused
// log is for logging in this package.
var userlog = logf.Log.WithName("user-resource")

// SetupUserWebhookWithManager registers the webhook for User in the manager.
func SetupUserWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&keycloakoperatorwebhippiedev1alpha1.User{}).
		WithValidator(&UserCustomValidator{}).
		WithDefaulter(&UserCustomDefaulter{}).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-keycloak-operator-webhippie-de-v1alpha1-user,mutating=true,failurePolicy=fail,sideEffects=None,groups=keycloak-operator.webhippie.de,resources=users,verbs=create;update,versions=v1alpha1,name=muser-v1alpha1.kb.io,admissionReviewVersions=v1

// UserCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind User when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type UserCustomDefaulter struct{}

var _ webhook.CustomDefaulter = &UserCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind User.
func (d *UserCustomDefaulter) Default(_ context.Context, obj runtime.Object) error {
	return nil
}

// +kubebuilder:webhook:path=/validate-keycloak-operator-webhippie-de-v1alpha1-user,mutating=false,failurePolicy=fail,sideEffects=None,groups=keycloak-operator.webhippie.de,resources=users,verbs=create;update,versions=v1alpha1,name=vuser-v1alpha1.kb.io,admissionReviewVersions=v1

// UserCustomValidator struct is responsible for validating the User resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type UserCustomValidator struct{}

var _ webhook.CustomValidator = &UserCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type User.
func (v *UserCustomValidator) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	user, ok := obj.(*keycloakoperatorwebhippiedev1alpha1.User)
	if !ok {
		return nil, fmt.Errorf("expected a User object but got %T", obj)
	}
	userlog.Info("Validation for User upon creation", "name", user.GetName())

	return nil, v.validate(user)
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type User.
func (v *UserCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	oldUser, ok := oldObj.(*keycloakoperatorwebhippiedev1alpha1.User)
	if !ok {
		return nil, fmt.Errorf("expected a User object for the oldObj but got %T", oldObj)
	}

	user, ok := newObj.(*keycloakoperatorwebhippiedev1alpha1.User)
	if !ok {
		return nil, fmt.Errorf("expected a User object for the newObj but got %T", newObj)
	}
	userlog.Info("Validation for User upon update", "name", user.GetName())

	if err := v.validate(user); err != nil {
		return nil, err
	}

	if oldUser.Spec.Username != user.Spec.Username {
		return nil, fmt.Errorf("spec.username is immutable and cannot be changed after creation")
	}

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type User.
func (v *UserCustomValidator) ValidateDelete(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

func (v *UserCustomValidator) validate(user *keycloakoperatorwebhippiedev1alpha1.User) error {
	if user.Spec.RealmRef == nil || strings.TrimSpace(user.Spec.RealmRef.Name) == "" {
		return fmt.Errorf("spec.realmRef.name must be set")
	}

	if strings.TrimSpace(user.Spec.Username) == "" {
		return fmt.Errorf("spec.username must be set")
	}

	return nil
}
