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

const groupFinalizer = "keycloak-operator.webhippie.de/group"

// GroupReconciler reconciles a Group object
type GroupReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=keycloak-operator.webhippie.de,resources=groups,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=keycloak-operator.webhippie.de,resources=groups/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=keycloak-operator.webhippie.de,resources=groups/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *GroupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)
	log.Info("Reconciling")

	instance := &v1alpha1.Group{}
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

	if !controllerutil.ContainsFinalizer(instance, groupFinalizer) {
		controllerutil.AddFinalizer(instance, groupFinalizer)
		if err := r.Update(ctx, instance); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to add finalizer: %w", err)
		}
		return ctrl.Result{RequeueAfter: time.Second}, nil
	}

	return r.reconcileGroup(ctx, instance, session)
}

func (r *GroupReconciler) handleDeletion(ctx context.Context, instance *v1alpha1.Group, session *KeycloakSession) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	if !controllerutil.ContainsFinalizer(instance, groupFinalizer) {
		return ctrl.Result{}, nil
	}

	if session != nil && instance.Status.KeycloakID != nil {
		log.Info("Deleting group from Keycloak", "id", *instance.Status.KeycloakID)
		if err := session.Client.DeleteGroup(ctx, session.Token.AccessToken, session.RealmName, *instance.Status.KeycloakID); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to delete group from Keycloak: %w", err)
		}
	}

	controllerutil.RemoveFinalizer(instance, groupFinalizer)
	if err := r.Update(ctx, instance); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to remove finalizer: %w", err)
	}

	return ctrl.Result{}, nil
}

func (r *GroupReconciler) reconcileGroup(ctx context.Context, instance *v1alpha1.Group, session *KeycloakSession) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	desired := groupToGocloak(instance)

	if instance.Status.KeycloakID == nil {
		groups, err := session.Client.GetGroups(ctx, session.Token.AccessToken, session.RealmName, gocloak.GetGroupsParams{
			Search: gocloak.StringP(instance.Spec.Name),
			Exact:  gocloak.BoolP(true),
		})
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to search for existing group: %w", err)
		}

		var match *gocloak.Group
		for _, g := range groups {
			if g.Name != nil && *g.Name == instance.Spec.Name {
				match = g
				break
			}
		}

		if match != nil {
			log.Info("Adopting existing Keycloak group", "id", *match.ID)
			if err := UpdateKeycloakIDStatus(ctx, r.Client, instance, match.ID); err != nil {
				return ctrl.Result{}, err
			}
		} else {
			id, err := session.Client.CreateGroup(ctx, session.Token.AccessToken, session.RealmName, desired)
			if err != nil {
				return ctrl.Result{}, fmt.Errorf("failed to create group in Keycloak: %w", err)
			}
			log.Info("Group created in Keycloak", "id", id)
			if err := UpdateKeycloakIDStatus(ctx, r.Client, instance, &id); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}
	}

	desired.ID = instance.Status.KeycloakID
	if err := session.Client.UpdateGroup(ctx, session.Token.AccessToken, session.RealmName, desired); err != nil {
		var apiErr *gocloak.APIError
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			log.Info("Group missing in Keycloak, clearing status to recreate", "id", *instance.Status.KeycloakID)
			if err := UpdateKeycloakIDStatus(ctx, r.Client, instance, nil); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{RequeueAfter: time.Second}, nil
		}
		return ctrl.Result{}, fmt.Errorf("failed to update group in Keycloak: %w", err)
	}

	log.Info("Group reconciled", "id", *instance.Status.KeycloakID)
	return ctrl.Result{}, nil
}

// groupToGocloak converts a Group CR spec into the gocloak representation used
// for create and update API calls.
func groupToGocloak(g *v1alpha1.Group) gocloak.Group {
	group := gocloak.Group{
		Name: gocloak.StringP(g.Spec.Name),
	}

	if g.Spec.Attributes != nil {
		attrs := g.Spec.Attributes
		group.Attributes = attrs
	}

	if g.Spec.RealmRoles != nil {
		roles := g.Spec.RealmRoles
		group.RealmRoles = roles
	}

	if g.Spec.ClientRoles != nil {
		cr := g.Spec.ClientRoles
		group.ClientRoles = cr
	}

	return group
}

// SetupWithManager sets up the controller with the Manager.
func (r *GroupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Group{}).
		Named("group").
		Complete(r)
}
