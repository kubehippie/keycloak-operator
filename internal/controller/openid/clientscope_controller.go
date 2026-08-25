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
	"errors"
	"fmt"
	"time"

	"github.com/Nerzal/gocloak/v14"
	"github.com/kubehippie/keycloak-operator/api/openid/v1alpha1"
	"github.com/kubehippie/keycloak-operator/internal/controller"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const clientScopeFinalizer = "keycloak-operator.webhippie.de/clientscope"

// ClientScopeReconciler reconciles a ClientScope object.
type ClientScopeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=openid.keycloak-operator.webhippie.de,resources=clientscopes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=openid.keycloak-operator.webhippie.de,resources=clientscopes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=openid.keycloak-operator.webhippie.de,resources=clientscopes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ClientScopeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)
	log.Info("Reconciling")

	instance := &v1alpha1.ClientScope{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("unable to fetch: %w", err)
	}

	session, err := controller.KeycloakSessionForRealm(ctx, r.Client, instance.Spec.RealmRef, req.Namespace)
	if err != nil {
		if !instance.DeletionTimestamp.IsZero() && apierrors.IsNotFound(err) {
			log.Info("Realm no longer exists, skipping Keycloak cleanup", "error", err.Error())
			return r.handleDeletion(ctx, instance, nil)
		}
		log.Error(err, "Unable to get Keycloak session")
		return ctrl.Result{RequeueAfter: controller.FailedKeycloakConnectionRetryPeriod}, nil
	}

	if !instance.DeletionTimestamp.IsZero() {
		return r.handleDeletion(ctx, instance, session)
	}

	if !controllerutil.ContainsFinalizer(instance, clientScopeFinalizer) {
		controllerutil.AddFinalizer(instance, clientScopeFinalizer)
		if err := r.Update(ctx, instance); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to add finalizer: %w", err)
		}
		return ctrl.Result{RequeueAfter: time.Second}, nil
	}

	return r.reconcileClientScope(ctx, instance, session)
}

func (r *ClientScopeReconciler) handleDeletion(ctx context.Context, instance *v1alpha1.ClientScope, session *controller.KeycloakSession) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	if !controllerutil.ContainsFinalizer(instance, clientScopeFinalizer) {
		return ctrl.Result{}, nil
	}

	if session != nil && instance.Status.KeycloakID != nil {
		log.Info("Deleting client scope from Keycloak", "id", *instance.Status.KeycloakID)
		if err := session.Client.DeleteClientScope(ctx, session.Token.AccessToken, session.RealmName, *instance.Status.KeycloakID); err != nil {
			var apiErr *gocloak.APIError
			if errors.As(err, &apiErr) && apiErr.Code == 404 {
				log.Info("Client scope already absent in Keycloak, skipping delete")
			} else {
				return ctrl.Result{}, fmt.Errorf("failed to delete client scope from Keycloak: %w", err)
			}
		}
	}

	controllerutil.RemoveFinalizer(instance, clientScopeFinalizer)
	if err := r.Update(ctx, instance); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to remove finalizer: %w", err)
	}

	return ctrl.Result{}, nil
}

func (r *ClientScopeReconciler) reconcileClientScope(ctx context.Context, instance *v1alpha1.ClientScope, session *controller.KeycloakSession) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	desired := clientScopeToGocloak(instance)

	if instance.Status.KeycloakID == nil {
		scopes, err := session.Client.GetClientScopes(ctx, session.Token.AccessToken, session.RealmName)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to search for existing client scope: %w", err)
		}

		var existing *gocloak.ClientScope
		for _, scope := range scopes {
			if scope != nil && scope.Name != nil && *scope.Name == instance.Spec.Name {
				existing = scope
				break
			}
		}

		if existing != nil {
			if existing.ID == nil {
				return ctrl.Result{}, fmt.Errorf("found existing client scope %q without ID", instance.Spec.Name)
			}
			log.Info("Adopting existing Keycloak client scope", "id", *existing.ID)
			if err := controller.UpdateKeycloakIDStatus(ctx, r.Client, instance, existing.ID); err != nil {
				return ctrl.Result{}, err
			}
		} else {
			id, err := session.Client.CreateClientScope(ctx, session.Token.AccessToken, session.RealmName, desired)
			if err != nil {
				return ctrl.Result{}, fmt.Errorf("failed to create client scope in Keycloak: %w", err)
			}
			log.Info("Client scope created in Keycloak", "id", id)
			if err := controller.UpdateKeycloakIDStatus(ctx, r.Client, instance, &id); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}
	}

	desired.ID = instance.Status.KeycloakID
	if err := session.Client.UpdateClientScope(ctx, session.Token.AccessToken, session.RealmName, desired); err != nil {
		var apiErr *gocloak.APIError
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			log.Info("Client scope missing in Keycloak, clearing status to recreate", "id", *instance.Status.KeycloakID)
			if err := controller.UpdateKeycloakIDStatus(ctx, r.Client, instance, nil); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{RequeueAfter: time.Second}, nil
		}
		return ctrl.Result{}, fmt.Errorf("failed to update client scope in Keycloak: %w", err)
	}

	log.Info("Client scope reconciled", "id", *instance.Status.KeycloakID)
	return ctrl.Result{}, nil
}

func clientScopeToGocloak(o *v1alpha1.ClientScope) gocloak.ClientScope {
	includeInTokenScope := boolStringDefaultTrue(o.Spec.IncludeInTokenScope)

	attributes := &gocloak.ClientScopeAttributes{
		IncludeInTokenScope: &includeInTokenScope,
	}

	if o.Spec.ConsentScreenText != nil {
		attributes.ConsentScreenText = o.Spec.ConsentScreenText
		displayOnConsentScreen := "true"
		attributes.DisplayOnConsentScreen = &displayOnConsentScreen
	} else {
		displayOnConsentScreen := "false"
		attributes.DisplayOnConsentScreen = &displayOnConsentScreen
	}

	return gocloak.ClientScope{
		Name:                  gocloak.StringP(o.Spec.Name),
		Description:           o.Spec.Description,
		Protocol:              gocloak.StringP(openidConnectProtocol),
		ClientScopeAttributes: attributes,
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClientScopeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.ClientScope{}).
		Named("openid-clientscope").
		Complete(r)
}
