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

package openid

import (
	"context"
	"fmt"
	"time"

	v1alpha1 "github.com/kubehippie/keycloak-operator/api/openid/v1alpha1"
	controller "github.com/kubehippie/keycloak-operator/internal/controller"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const defaultScopesFinalizer = "keycloak-operator.webhippie.de/defaultscopes"

// DefaultScopesReconciler reconciles a DefaultScopes object.
type DefaultScopesReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=openid.keycloak-operator.webhippie.de,resources=defaultscopes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=openid.keycloak-operator.webhippie.de,resources=defaultscopes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=openid.keycloak-operator.webhippie.de,resources=defaultscopes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *DefaultScopesReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)
	log.Info("Reconciling")

	instance := &v1alpha1.DefaultScopes{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("unable to fetch: %w", err)
	}

	session, idOfClient, err := controller.KeycloakSessionForClient(ctx, r.Client, instance.Spec.ClientRef, req.Namespace)
	if err != nil {
		if !instance.DeletionTimestamp.IsZero() && apierrors.IsNotFound(err) {
			log.Info("OpenID client no longer exists, skipping Keycloak cleanup", "error", err.Error())
			return r.handleDeletion(ctx, instance, nil, "")
		}
		log.Error(err, "Unable to get Keycloak session")
		return ctrl.Result{RequeueAfter: controller.FailedKeycloakConnectionRetryPeriod}, nil
	}

	if !instance.DeletionTimestamp.IsZero() {
		return r.handleDeletion(ctx, instance, session, idOfClient)
	}

	if !controllerutil.ContainsFinalizer(instance, defaultScopesFinalizer) {
		controllerutil.AddFinalizer(instance, defaultScopesFinalizer)
		if err := r.Update(ctx, instance); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to add finalizer: %w", err)
		}
		return ctrl.Result{RequeueAfter: time.Second}, nil
	}

	return r.reconcileDefaultScopes(ctx, instance, session, idOfClient)
}

func (r *DefaultScopesReconciler) handleDeletion(ctx context.Context, instance *v1alpha1.DefaultScopes, session *controller.KeycloakSession, idOfClient string) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	if !controllerutil.ContainsFinalizer(instance, defaultScopesFinalizer) {
		return ctrl.Result{}, nil
	}

	if session != nil {
		currentScopes, err := session.Client.GetClientsDefaultScopes(ctx, session.Token.AccessToken, session.RealmName, idOfClient)
		if err != nil {
			if isNotFoundAPIError(err) {
				log.Info("Client already absent in Keycloak, skipping default scope cleanup")
			} else {
				return ctrl.Result{}, fmt.Errorf("failed to read default scopes from Keycloak: %w", err)
			}
			currentScopes = nil
		}

		desiredNames := make(map[string]struct{}, len(instance.Spec.DefaultScopes))
		for _, scopeName := range instance.Spec.DefaultScopes {
			desiredNames[scopeName] = struct{}{}
		}

		for _, scope := range currentScopes {
			if scope == nil || scope.Name == nil || scope.ID == nil {
				continue
			}
			if _, ok := desiredNames[*scope.Name]; !ok {
				continue
			}
			if err := session.Client.RemoveDefaultScopeFromClient(ctx, session.Token.AccessToken, session.RealmName, idOfClient, *scope.ID); err != nil {
				log.Error(err, "Failed to remove default scope during deletion", "scope", *scope.Name)
			}
		}
	}

	controllerutil.RemoveFinalizer(instance, defaultScopesFinalizer)
	if err := r.Update(ctx, instance); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to remove finalizer: %w", err)
	}

	return ctrl.Result{}, nil
}

func (r *DefaultScopesReconciler) reconcileDefaultScopes(ctx context.Context, instance *v1alpha1.DefaultScopes, session *controller.KeycloakSession, idOfClient string) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	realmScopes, err := session.Client.GetClientScopes(ctx, session.Token.AccessToken, session.RealmName)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to list client scopes from Keycloak: %w", err)
	}

	realmScopeIDs := make(map[string]string, len(realmScopes))
	for _, scope := range realmScopes {
		if scope == nil || scope.Name == nil || scope.ID == nil {
			continue
		}
		realmScopeIDs[*scope.Name] = *scope.ID
	}

	currentScopes, err := session.Client.GetClientsDefaultScopes(ctx, session.Token.AccessToken, session.RealmName, idOfClient)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to list current default scopes from Keycloak: %w", err)
	}

	toAdd, toRemove, err := resolveDefaultScopesPlan(realmScopeIDs, instance.Spec.DefaultScopes, currentScopes)
	if err != nil {
		return ctrl.Result{}, err
	}

	for _, scopeID := range toAdd {
		if err := session.Client.AddDefaultScopeToClient(ctx, session.Token.AccessToken, session.RealmName, idOfClient, scopeID); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to attach default scope %q: %w", scopeID, err)
		}
	}

	for _, scopeID := range toRemove {
		if err := session.Client.RemoveDefaultScopeFromClient(ctx, session.Token.AccessToken, session.RealmName, idOfClient, scopeID); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to detach default scope %q: %w", scopeID, err)
		}
	}

	log.Info("Default scopes reconciled", "added", len(toAdd), "removed", len(toRemove))
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DefaultScopesReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.DefaultScopes{}).
		Named("openid-defaultscopes").
		Complete(r)
}
