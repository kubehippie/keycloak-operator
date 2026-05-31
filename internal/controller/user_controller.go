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

package controller

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v14"
	v1alpha1 "github.com/kubehippie/keycloak-operator/api/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const userFinalizer = "keycloak-operator.webhippie.de/user"

// UserReconciler reconciles a User object
type UserReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=keycloak-operator.webhippie.de,resources=users,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=keycloak-operator.webhippie.de,resources=users/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=keycloak-operator.webhippie.de,resources=users/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *UserReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)
	log.Info("Reconciling")

	instance := &v1alpha1.User{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("unable to fetch: %w", err)
	}

	session, err := keycloakSessionForRealm(ctx, r.Client, instance.Spec.RealmRef, req.Namespace)
	if err != nil {
		log.Error(err, "Unable to get Keycloak session")
		return ctrl.Result{RequeueAfter: failedKeycloakConnectionRetryPeriod}, nil
	}

	if !instance.DeletionTimestamp.IsZero() {
		return r.handleDeletion(ctx, instance, session)
	}

	if !controllerutil.ContainsFinalizer(instance, userFinalizer) {
		controllerutil.AddFinalizer(instance, userFinalizer)
		if err := r.Update(ctx, instance); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to add finalizer: %w", err)
		}
		return ctrl.Result{Requeue: true}, nil
	}

	return r.reconcileUser(ctx, instance, session)
}

func (r *UserReconciler) handleDeletion(ctx context.Context, instance *v1alpha1.User, session *keycloakSession) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	if !controllerutil.ContainsFinalizer(instance, userFinalizer) {
		return ctrl.Result{}, nil
	}

	if instance.Status.KeycloakID != nil {
		log.Info("Deleting user from Keycloak", "id", *instance.Status.KeycloakID)
		if err := session.Client.DeleteUser(ctx, session.Token.AccessToken, session.RealmName, *instance.Status.KeycloakID); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to delete user from Keycloak: %w", err)
		}
	}

	controllerutil.RemoveFinalizer(instance, userFinalizer)
	if err := r.Update(ctx, instance); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to remove finalizer: %w", err)
	}

	return ctrl.Result{}, nil
}

func (r *UserReconciler) reconcileUser(ctx context.Context, instance *v1alpha1.User, session *keycloakSession) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	desired := userToGocloak(instance)

	if instance.Status.KeycloakID == nil {
		users, err := session.Client.GetUsers(ctx, session.Token.AccessToken, session.RealmName, gocloak.GetUsersParams{
			Username: gocloak.StringP(instance.Spec.Username),
			Exact:    gocloak.BoolP(true),
		})
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to search for existing user: %w", err)
		}

		if len(users) > 0 {
			log.Info("Adopting existing Keycloak user", "id", *users[0].ID)
			if err := updateKeycloakIDStatus(ctx, r.Client, instance, users[0].ID); err != nil {
				return ctrl.Result{}, err
			}
		} else {
			id, err := session.Client.CreateUser(ctx, session.Token.AccessToken, session.RealmName, desired)
			if err != nil {
				return ctrl.Result{}, fmt.Errorf("failed to create user in Keycloak: %w", err)
			}
			log.Info("User created in Keycloak", "id", id)
			if err := updateKeycloakIDStatus(ctx, r.Client, instance, &id); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}
	}

	desired.ID = instance.Status.KeycloakID
	if err := session.Client.UpdateUser(ctx, session.Token.AccessToken, session.RealmName, desired); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to update user in Keycloak: %w", err)
	}

	log.Info("User reconciled", "id", *instance.Status.KeycloakID)
	return ctrl.Result{}, nil
}

// userToGocloak converts a User CR spec into the gocloak representation used
// for create and update API calls.
func userToGocloak(u *v1alpha1.User) gocloak.User {
	user := gocloak.User{
		Username:      gocloak.StringP(u.Spec.Username),
		Enabled:       u.Spec.Enabled,
		EmailVerified: u.Spec.EmailVerified,
		FirstName:     u.Spec.FirstName,
		LastName:      u.Spec.LastName,
		Email:         u.Spec.Email,
	}

	if u.Spec.Attributes != nil {
		attrs := u.Spec.Attributes
		user.Attributes = &attrs
	}

	if u.Spec.RequiredActions != nil {
		actions := u.Spec.RequiredActions
		user.RequiredActions = &actions
	}

	return user
}

// SetupWithManager sets up the controller with the Manager.
func (r *UserReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.User{}).
		Named("user").
		Complete(r)
}
