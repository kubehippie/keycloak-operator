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
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// isServiceAccountsEnabled reports whether the client's service account
// (client credentials grant) is enabled.
func isServiceAccountsEnabled(o *v1alpha1.OpenIDClient) bool {
	return o.Spec.ServiceAccountsEnabled != nil && *o.Spec.ServiceAccountsEnabled
}

const openIDClientFinalizer = "keycloak-operator.webhippie.de/openidclient"

// OpenIDClientReconciler reconciles a OpenIDClient object
type OpenIDClientReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=openid.keycloak-operator.webhippie.de,resources=openidclients,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=openid.keycloak-operator.webhippie.de,resources=openidclients/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=openid.keycloak-operator.webhippie.de,resources=openidclients/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *OpenIDClientReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)
	log.Info("Reconciling")

	instance := &v1alpha1.OpenIDClient{}
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

	if !controllerutil.ContainsFinalizer(instance, openIDClientFinalizer) {
		controllerutil.AddFinalizer(instance, openIDClientFinalizer)
		if err := r.Update(ctx, instance); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to add finalizer: %w", err)
		}
		return ctrl.Result{RequeueAfter: time.Second}, nil
	}

	return r.reconcileOpenIDClient(ctx, instance, session)
}

func (r *OpenIDClientReconciler) handleDeletion(ctx context.Context, instance *v1alpha1.OpenIDClient, session *controller.KeycloakSession) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	if !controllerutil.ContainsFinalizer(instance, openIDClientFinalizer) {
		return ctrl.Result{}, nil
	}

	if session != nil && instance.Status.KeycloakID != nil {
		log.Info("Deleting client from Keycloak", "id", *instance.Status.KeycloakID)
		if err := session.Client.DeleteClient(ctx, session.Token.AccessToken, session.RealmName, *instance.Status.KeycloakID); err != nil {
			var apiErr *gocloak.APIError
			if errors.As(err, &apiErr) && apiErr.Code == 404 {
				log.Info("Client already absent in Keycloak, skipping delete")
			} else {
				return ctrl.Result{}, fmt.Errorf("failed to delete client from Keycloak: %w", err)
			}
		}
	}

	controllerutil.RemoveFinalizer(instance, openIDClientFinalizer)
	if err := r.Update(ctx, instance); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to remove finalizer: %w", err)
	}

	return ctrl.Result{}, nil
}

func (r *OpenIDClientReconciler) reconcileOpenIDClient(ctx context.Context, instance *v1alpha1.OpenIDClient, session *controller.KeycloakSession) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	desired, err := openIDClientToGocloak(ctx, r.Client, instance, instance.Namespace)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to build client representation: %w", err)
	}

	if instance.Status.KeycloakID == nil {
		clients, err := session.Client.GetClients(ctx, session.Token.AccessToken, session.RealmName, gocloak.GetClientsParams{
			ClientID: gocloak.StringP(instance.Spec.ClientID),
		})
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to search for existing client: %w", err)
		}

		var existing *gocloak.Client
		for _, c := range clients {
			if c.ClientID != nil && *c.ClientID == instance.Spec.ClientID {
				existing = c
				break
			}
		}

		if existing != nil {
			log.Info("Adopting existing Keycloak client", "id", *existing.ID)
			if err := controller.UpdateKeycloakIDStatus(ctx, r.Client, instance, existing.ID); err != nil {
				return ctrl.Result{}, err
			}
		} else {
			id, err := session.Client.CreateClient(ctx, session.Token.AccessToken, session.RealmName, desired)
			if err != nil {
				return ctrl.Result{}, fmt.Errorf("failed to create client in Keycloak: %w", err)
			}
			log.Info("Client created in Keycloak", "id", id)
			if err := controller.UpdateKeycloakIDStatus(ctx, r.Client, instance, &id); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}
	}

	desired.ID = instance.Status.KeycloakID
	if err := session.Client.UpdateClient(ctx, session.Token.AccessToken, session.RealmName, desired); err != nil {
		var apiErr *gocloak.APIError
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			log.Info("Client missing in Keycloak, clearing status to recreate", "id", *instance.Status.KeycloakID)
			if err := controller.UpdateKeycloakIDStatus(ctx, r.Client, instance, nil); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{RequeueAfter: time.Second}, nil
		}
		return ctrl.Result{}, fmt.Errorf("failed to update client in Keycloak: %w", err)
	}

	if err := r.reconcileServiceAccountRoles(ctx, instance, session); err != nil {
		return ctrl.Result{}, err
	}

	log.Info("Client reconciled", "id", *instance.Status.KeycloakID)
	return ctrl.Result{}, nil
}

// reconcileServiceAccountRoles ensures the client's service account user has
// the realm and client roles configured in
// spec.serviceAccountRealmRoles/spec.serviceAccountClientRoles. It is a
// no-op when serviceAccountsEnabled is not set.
func (r *OpenIDClientReconciler) reconcileServiceAccountRoles(ctx context.Context, instance *v1alpha1.OpenIDClient, session *controller.KeycloakSession) error {
	if !isServiceAccountsEnabled(instance) {
		return nil
	}

	if len(instance.Spec.ServiceAccountRealmRoles) == 0 && len(instance.Spec.ServiceAccountClientRoles) == 0 &&
		len(instance.Status.ServiceAccountRealmRoles) == 0 && len(instance.Status.ServiceAccountClientRoles) == 0 {
		return nil
	}

	serviceAccount, err := session.Client.GetClientServiceAccount(ctx, session.Token.AccessToken, session.RealmName, *instance.Status.KeycloakID)
	if err != nil {
		return fmt.Errorf("failed to resolve service account user: %w", err)
	}
	if serviceAccount == nil || serviceAccount.ID == nil {
		return fmt.Errorf("service account user not found for client %q", instance.Spec.ClientID)
	}

	if err := r.reconcileServiceAccountRealmRoles(ctx, instance, session, *serviceAccount.ID); err != nil {
		return err
	}

	return r.reconcileServiceAccountClientRoles(ctx, instance, session, *serviceAccount.ID)
}

// reconcileServiceAccountRealmRoles ensures the service account user's
// directly assigned realm roles match spec.serviceAccountRealmRoles.
// Assignment changes are computed relative to
// status.serviceAccountRealmRoles (the set previously applied by this
// operator), so roles assigned outside of this resource (e.g. default realm
// roles) are left untouched, while names removed from
// spec.serviceAccountRealmRoles are revoked.
func (r *OpenIDClientReconciler) reconcileServiceAccountRealmRoles(ctx context.Context, instance *v1alpha1.OpenIDClient, session *controller.KeycloakSession, userID string) error {
	toAdd, toRemove := diffNames(instance.Spec.ServiceAccountRealmRoles, instance.Status.ServiceAccountRealmRoles)
	if len(toAdd) == 0 && len(toRemove) == 0 {
		return nil
	}

	for _, name := range toAdd {
		role, err := session.Client.GetRealmRole(ctx, session.Token.AccessToken, session.RealmName, name)
		if err != nil {
			return fmt.Errorf("failed to resolve service account realm role %q: %w", name, err)
		}
		if err := session.Client.AddRealmRoleToUser(ctx, session.Token.AccessToken, session.RealmName, userID, []gocloak.Role{*role}); err != nil {
			return fmt.Errorf("failed to assign service account realm role %q: %w", name, err)
		}
	}

	for _, name := range toRemove {
		role, err := session.Client.GetRealmRole(ctx, session.Token.AccessToken, session.RealmName, name)
		if err != nil {
			if isNotFoundAPIError(err) {
				continue
			}
			return fmt.Errorf("failed to resolve service account realm role %q for removal: %w", name, err)
		}
		if err := session.Client.DeleteRealmRoleFromUser(ctx, session.Token.AccessToken, session.RealmName, userID, []gocloak.Role{*role}); err != nil {
			return fmt.Errorf("failed to revoke service account realm role %q: %w", name, err)
		}
	}

	instance.Status.ServiceAccountRealmRoles = instance.Spec.ServiceAccountRealmRoles
	if err := r.Status().Update(ctx, instance); err != nil {
		return fmt.Errorf("failed to update applied service account realm roles status: %w", err)
	}

	return nil
}

// reconcileServiceAccountClientRoles ensures the service account user's
// client role assignments match spec.serviceAccountClientRoles, a map of
// target clientID (e.g. "realm-management") to role names on that client
// (e.g. "view-users", "manage-users", "query-users"). Assignment changes are
// computed per target client relative to
// status.serviceAccountClientRoles, following the same "leave
// externally-managed roles untouched" semantics as realm roles.
func (r *OpenIDClientReconciler) reconcileServiceAccountClientRoles(ctx context.Context, instance *v1alpha1.OpenIDClient, session *controller.KeycloakSession, userID string) error {
	desired := instance.Spec.ServiceAccountClientRoles
	previous := instance.Status.ServiceAccountClientRoles

	targetClientIDs := make(map[string]struct{}, len(desired)+len(previous))
	for name := range desired {
		targetClientIDs[name] = struct{}{}
	}
	for name := range previous {
		targetClientIDs[name] = struct{}{}
	}

	changed := false
	for targetClientID := range targetClientIDs {
		toAdd, toRemove := diffNames(desired[targetClientID], previous[targetClientID])
		if len(toAdd) == 0 && len(toRemove) == 0 {
			continue
		}
		changed = true

		idOfTargetClient, err := findClientIDByClientID(ctx, session, targetClientID)
		if err != nil {
			return fmt.Errorf("failed to resolve client %q: %w", targetClientID, err)
		}
		if idOfTargetClient == "" {
			return fmt.Errorf("client %q not found in realm %q", targetClientID, session.RealmName)
		}

		for _, name := range toAdd {
			role, err := session.Client.GetClientRole(ctx, session.Token.AccessToken, session.RealmName, idOfTargetClient, name)
			if err != nil {
				return fmt.Errorf("failed to resolve client role %q on client %q: %w", name, targetClientID, err)
			}
			if err := session.Client.AddClientRolesToUser(ctx, session.Token.AccessToken, session.RealmName, idOfTargetClient, userID, []gocloak.Role{*role}); err != nil {
				return fmt.Errorf("failed to assign client role %q on client %q: %w", name, targetClientID, err)
			}
		}

		for _, name := range toRemove {
			role, err := session.Client.GetClientRole(ctx, session.Token.AccessToken, session.RealmName, idOfTargetClient, name)
			if err != nil {
				if isNotFoundAPIError(err) {
					continue
				}
				return fmt.Errorf("failed to resolve client role %q on client %q for removal: %w", name, targetClientID, err)
			}
			if err := session.Client.DeleteClientRolesFromUser(ctx, session.Token.AccessToken, session.RealmName, idOfTargetClient, userID, []gocloak.Role{*role}); err != nil {
				return fmt.Errorf("failed to revoke client role %q on client %q: %w", name, targetClientID, err)
			}
		}
	}

	if changed {
		instance.Status.ServiceAccountClientRoles = desired
		if err := r.Status().Update(ctx, instance); err != nil {
			return fmt.Errorf("failed to update applied service account client roles status: %w", err)
		}
	}

	return nil
}

// findClientIDByClientID resolves the internal Keycloak UUID for the client
// with the given clientID, returning an empty string (without error) when no
// such client exists.
func findClientIDByClientID(ctx context.Context, session *controller.KeycloakSession, clientID string) (string, error) {
	clients, err := session.Client.GetClients(ctx, session.Token.AccessToken, session.RealmName, gocloak.GetClientsParams{
		ClientID: gocloak.StringP(clientID),
	})
	if err != nil {
		return "", err
	}

	for _, c := range clients {
		if c.ClientID != nil && *c.ClientID == clientID && c.ID != nil {
			return *c.ID, nil
		}
	}

	return "", nil
}

// diffNames compares a desired set of names against the set previously
// applied by this operator (as recorded in status), returning the names
// that need to be added and the names that need to be removed.
func diffNames(desired, previouslyApplied []string) (toAdd, toRemove []string) {
	desiredSet := make(map[string]struct{}, len(desired))
	for _, name := range desired {
		desiredSet[name] = struct{}{}
	}

	previousSet := make(map[string]struct{}, len(previouslyApplied))
	for _, name := range previouslyApplied {
		previousSet[name] = struct{}{}
	}

	for _, name := range desired {
		if _, ok := previousSet[name]; !ok {
			toAdd = append(toAdd, name)
		}
	}

	for _, name := range previouslyApplied {
		if _, ok := desiredSet[name]; !ok {
			toRemove = append(toRemove, name)
		}
	}

	return toAdd, toRemove
}

// openIDClientToGocloak converts an OpenIDClient CR spec into the gocloak
// representation used for create and update API calls. The client secret (if
// configured) is resolved from the cluster at call time.
func openIDClientToGocloak(ctx context.Context, cl client.Client, o *v1alpha1.OpenIDClient, ns string) (gocloak.Client, error) {
	c := gocloak.Client{
		ClientID:                  gocloak.StringP(o.Spec.ClientID),
		Name:                      o.Spec.Name,
		Description:               o.Spec.Description,
		Enabled:                   o.Spec.Enabled,
		ClientAuthenticatorType:   o.Spec.ClientAuthenticatorType,
		StandardFlowEnabled:       o.Spec.StandardFlowEnabled,
		ImplicitFlowEnabled:       o.Spec.ImplicitFlowEnabled,
		DirectAccessGrantsEnabled: o.Spec.DirectAccessGrantsEnabled,
		ServiceAccountsEnabled:    o.Spec.ServiceAccountsEnabled,
		FrontChannelLogout:        o.Spec.FrontChannelLogoutEnabled,
		RedirectURIs:              o.Spec.ValidRedirectURIs,
		WebOrigins:                o.Spec.WebOrigins,
		RootURL:                   o.Spec.RootURL,
		BaseURL:                   o.Spec.BaseURL,
		AdminURL:                  o.Spec.AdminURL,
		FullScopeAllowed:          o.Spec.FullScopeAllowed,
		Attributes:                o.Spec.Attributes,
	}

	switch o.Spec.AccessType {
	case "PUBLIC":
		c.PublicClient = gocloak.BoolP(true)
		c.BearerOnly = gocloak.BoolP(false)
	case "BEARER-ONLY":
		c.PublicClient = gocloak.BoolP(false)
		c.BearerOnly = gocloak.BoolP(true)
	default: // CONFIDENTIAL
		c.PublicClient = gocloak.BoolP(false)
		c.BearerOnly = gocloak.BoolP(false)
	}

	if o.Spec.ClientSecret != nil {
		secret, err := controller.ResolveSecretKeyRefOrVal(ctx, cl, o.Spec.ClientSecret, ns)
		if err != nil {
			return gocloak.Client{}, fmt.Errorf("unable to resolve clientSecret: %w", err)
		}
		c.Secret = gocloak.StringP(secret)
	}

	return c, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OpenIDClientReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := controller.RegisterRefIndex(mgr, &v1alpha1.OpenIDClient{}, openIDClientSecretIndexField, openIDClientSecretRefKeys); err != nil {
		return err
	}

	newList := func() client.ObjectList { return &v1alpha1.OpenIDClientList{} }

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.OpenIDClient{}).
		Watches(&corev1.Secret{}, controller.RefEventHandler(mgr.GetClient(), newList, openIDClientSecretIndexField)).
		Named("openid-openidclient").
		Complete(r)
}

const openIDClientSecretIndexField = ".spec.secretRefs"

// openIDClientSecretRefKeys returns the RefIndexKey values for every Secret
// the given OpenIDClient instance may reference (currently the clientSecret).
func openIDClientSecretRefKeys(obj client.Object) []string {
	c, ok := obj.(*v1alpha1.OpenIDClient)
	if !ok || c.Spec.ClientSecret == nil || c.Spec.ClientSecret.SecretKeyRef == nil {
		return nil
	}

	ref := c.Spec.ClientSecret.SecretKeyRef
	ns := ref.Namespace
	if ns == "" {
		ns = c.Namespace
	}

	return []string{controller.RefIndexKey(ns, ref.Name)}
}
