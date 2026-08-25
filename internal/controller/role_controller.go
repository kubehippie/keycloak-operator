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
	"errors"
	"fmt"
	"time"

	"github.com/Nerzal/gocloak/v14"
	v1alpha1 "github.com/kubehippie/keycloak-operator/api/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const roleFinalizer = "keycloak-operator.webhippie.de/role"

// RoleReconciler reconciles a Role object
type RoleReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=keycloak-operator.webhippie.de,resources=roles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=keycloak-operator.webhippie.de,resources=roles/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=keycloak-operator.webhippie.de,resources=roles/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *RoleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)
	log.Info("Reconciling")

	instance := &v1alpha1.Role{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("unable to fetch: %w", err)
	}

	session, err := KeycloakSessionForRealm(ctx, r.Client, instance.Spec.RealmRef, req.Namespace)
	if err != nil {
		if !instance.DeletionTimestamp.IsZero() && apierrors.IsNotFound(err) {
			log.Info("Realm no longer exists, skipping Keycloak cleanup", "error", err.Error())
			return r.handleDeletion(ctx, instance, nil)
		}
		log.Error(err, "Unable to get Keycloak session")
		return ctrl.Result{RequeueAfter: FailedKeycloakConnectionRetryPeriod}, nil
	}

	if !instance.DeletionTimestamp.IsZero() {
		return r.handleDeletion(ctx, instance, session)
	}

	if !controllerutil.ContainsFinalizer(instance, roleFinalizer) {
		controllerutil.AddFinalizer(instance, roleFinalizer)
		if err := r.Update(ctx, instance); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to add finalizer: %w", err)
		}
		return ctrl.Result{RequeueAfter: time.Second}, nil
	}

	return r.reconcileRole(ctx, instance, session)
}

func (r *RoleReconciler) handleDeletion(ctx context.Context, instance *v1alpha1.Role, session *KeycloakSession) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	if !controllerutil.ContainsFinalizer(instance, roleFinalizer) {
		return ctrl.Result{}, nil
	}

	if session != nil && instance.Status.KeycloakID != nil {
		log.Info("Deleting role from Keycloak", "id", *instance.Status.KeycloakID)
		if err := session.Client.DeleteRealmRole(ctx, session.Token.AccessToken, session.RealmName, instance.Spec.Name); err != nil {
			var apiErr *gocloak.APIError
			if errors.As(err, &apiErr) && apiErr.Code == 404 {
				log.Info("Role already absent in Keycloak, skipping delete")
			} else {
				return ctrl.Result{}, fmt.Errorf("failed to delete role from Keycloak: %w", err)
			}
		}
	}

	controllerutil.RemoveFinalizer(instance, roleFinalizer)
	if err := r.Update(ctx, instance); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to remove finalizer: %w", err)
	}

	return ctrl.Result{}, nil
}

func (r *RoleReconciler) reconcileRole(ctx context.Context, instance *v1alpha1.Role, session *KeycloakSession) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	desired := roleToGocloak(instance)

	if instance.Status.KeycloakID == nil {
		match, err := session.Client.GetRealmRole(ctx, session.Token.AccessToken, session.RealmName, instance.Spec.Name)
		if err != nil {
			var apiErr *gocloak.APIError
			if !errors.As(err, &apiErr) || apiErr.Code != 404 {
				return ctrl.Result{}, fmt.Errorf("failed to search for existing role: %w", err)
			}
		}

		if match != nil {
			log.Info("Adopting existing Keycloak role", "id", *match.ID)
			if err := UpdateKeycloakIDStatus(ctx, r.Client, instance, match.ID); err != nil {
				return ctrl.Result{}, err
			}
		} else {
			id, err := session.Client.CreateRealmRole(ctx, session.Token.AccessToken, session.RealmName, desired)
			if err != nil {
				return ctrl.Result{}, fmt.Errorf("failed to create role in Keycloak: %w", err)
			}
			log.Info("Role created in Keycloak", "id", id)
			if err := UpdateKeycloakIDStatus(ctx, r.Client, instance, &id); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}
	}

	desired.ID = instance.Status.KeycloakID
	if err := session.Client.UpdateRealmRoleByID(ctx, session.Token.AccessToken, session.RealmName, *instance.Status.KeycloakID, desired); err != nil {
		var apiErr *gocloak.APIError
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			log.Info("Role missing in Keycloak, clearing status to recreate", "id", *instance.Status.KeycloakID)
			if err := UpdateKeycloakIDStatus(ctx, r.Client, instance, nil); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{RequeueAfter: time.Second}, nil
		}
		return ctrl.Result{}, fmt.Errorf("failed to update role in Keycloak: %w", err)
	}

	log.Info("Role reconciled", "id", *instance.Status.KeycloakID)
	return ctrl.Result{}, nil
}

// roleToGocloak converts a Role CR spec into the gocloak representation used
// for create and update API calls.
func roleToGocloak(rl *v1alpha1.Role) gocloak.Role {
	role := gocloak.Role{
		Name: gocloak.StringP(rl.Spec.Name),
	}

	if rl.Spec.Description != nil {
		role.Description = rl.Spec.Description
	}

	if rl.Spec.Attributes != nil {
		attrs := rl.Spec.Attributes
		role.Attributes = attrs
	}

	return role
}

// SetupWithManager sets up the controller with the Manager.
func (r *RoleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Role{}).
		Named("role").
		Complete(r)
}
