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

	"github.com/kubehippie/keycloak-operator/api/common"
	v1alpha1 "github.com/kubehippie/keycloak-operator/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// nolint:unused
// log is for logging in this package.
var userlog = logf.Log.WithName("user-resource")

// SetupUserWebhookWithManager registers the webhook for User in the manager.
func SetupUserWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &v1alpha1.User{}).
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

var _ admission.Defaulter[*v1alpha1.User] = &UserCustomDefaulter{}

// Default implements admission.Defaulter so a webhook will be registered for the Kind User.
func (d *UserCustomDefaulter) Default(_ context.Context, user *v1alpha1.User) error {
	_ = user
	return nil
}

// +kubebuilder:webhook:path=/validate-keycloak-operator-webhippie-de-v1alpha1-user,mutating=false,failurePolicy=fail,sideEffects=None,groups=keycloak-operator.webhippie.de,resources=users,verbs=create;update,versions=v1alpha1,name=vuser-v1alpha1.kb.io,admissionReviewVersions=v1

// UserCustomValidator struct is responsible for validating the User resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type UserCustomValidator struct{}

var _ admission.Validator[*v1alpha1.User] = &UserCustomValidator{}

// ValidateCreate implements admission.Validator so a webhook will be registered for the type User.
func (v *UserCustomValidator) ValidateCreate(_ context.Context, user *v1alpha1.User) (admission.Warnings, error) {
	userlog.Info("Validation for User upon creation", "name", user.GetName())

	return nil, v.validate(user)
}

// ValidateUpdate implements admission.Validator so a webhook will be registered for the type User.
func (v *UserCustomValidator) ValidateUpdate(_ context.Context, oldUser, user *v1alpha1.User) (admission.Warnings, error) {
	userlog.Info("Validation for User upon update", "name", user.GetName())

	if err := v.validate(user); err != nil {
		return nil, err
	}

	if oldUser.Spec.Username != user.Spec.Username {
		return nil, fmt.Errorf("spec.username is immutable and cannot be changed after creation")
	}

	return nil, nil
}

// ValidateDelete implements admission.Validator so a webhook will be registered for the type User.
func (v *UserCustomValidator) ValidateDelete(_ context.Context, user *v1alpha1.User) (admission.Warnings, error) {
	_ = user
	return nil, nil
}

func (v *UserCustomValidator) validate(user *v1alpha1.User) error {
	if user.Spec.RealmRef == nil || strings.TrimSpace(user.Spec.RealmRef.Name) == "" {
		return fmt.Errorf("spec.realmRef.name must be set")
	}

	if strings.TrimSpace(user.Spec.Username) == "" {
		return fmt.Errorf("spec.username must be set")
	}

	if user.Spec.Password != nil {
		if err := v.validateSecretKeyRefOrVal("spec.password", user.Spec.Password); err != nil {
			return err
		}
	}

	return nil
}

func (v *UserCustomValidator) validateSecretKeyRefOrVal(field string, ref *common.SecretKeyRefOrVal) error {
	if strings.TrimSpace(ref.Value) != "" {
		return nil
	}

	if ref.SecretKeyRef == nil {
		return fmt.Errorf("%s must set value or secretKeyRef", field)
	}

	if strings.TrimSpace(ref.SecretKeyRef.Name) == "" || strings.TrimSpace(ref.SecretKeyRef.Key) == "" {
		return fmt.Errorf("%s.secretKeyRef must set name and key", field)
	}

	return nil
}
