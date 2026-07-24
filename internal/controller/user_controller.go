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

	"github.com/Nerzal/gocloak/v14"
	v1alpha1 "github.com/kubehippie/keycloak-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
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

	session, err := KeycloakSessionForRealm(ctx, r.Client, instance.Spec.RealmRef, req.Namespace)
	if err != nil {
		log.Error(err, "Unable to get Keycloak session")
		return ctrl.Result{RequeueAfter: FailedKeycloakConnectionRetryPeriod}, nil
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

func (r *UserReconciler) handleDeletion(ctx context.Context, instance *v1alpha1.User, session *KeycloakSession) (ctrl.Result, error) {
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

func (r *UserReconciler) reconcileUser(ctx context.Context, instance *v1alpha1.User, session *KeycloakSession) (ctrl.Result, error) {
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
			if err := UpdateKeycloakIDStatus(ctx, r.Client, instance, users[0].ID); err != nil {
				return ctrl.Result{}, err
			}
		} else {
			id, err := session.Client.CreateUser(ctx, session.Token.AccessToken, session.RealmName, desired)
			if err != nil {
				return ctrl.Result{}, fmt.Errorf("failed to create user in Keycloak: %w", err)
			}
			log.Info("User created in Keycloak", "id", id)
			if err := r.setPassword(ctx, instance, session, id); err != nil {
				return ctrl.Result{}, err
			}
			if err := UpdateKeycloakIDStatus(ctx, r.Client, instance, &id); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		desired.ID = instance.Status.KeycloakID
		if err := session.Client.UpdateUser(ctx, session.Token.AccessToken, session.RealmName, desired); err != nil {
			var apiErr *gocloak.APIError
			if errors.As(err, &apiErr) && apiErr.Code == 404 {
				log.Info("User missing in Keycloak, clearing status to recreate", "id", *instance.Status.KeycloakID)
				if err := UpdateKeycloakIDStatus(ctx, r.Client, instance, nil); err != nil {
					return ctrl.Result{}, err
				}
				return ctrl.Result{Requeue: true}, nil
			}
			return ctrl.Result{}, fmt.Errorf("failed to update user in Keycloak: %w", err)
		}

		if err := r.setPassword(ctx, instance, session, *instance.Status.KeycloakID); err != nil {
			return ctrl.Result{}, err
		}
	}

	if err := r.reconcileGroups(ctx, instance, session, *instance.Status.KeycloakID); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.reconcileRealmRoles(ctx, instance, session, *instance.Status.KeycloakID); err != nil {
		return ctrl.Result{}, err
	}

	log.Info("User reconciled", "id", *instance.Status.KeycloakID)
	return ctrl.Result{}, nil
}

// reconcileGroups ensures the user's Keycloak group memberships match
// instance.Spec.Groups. Membership changes are computed relative to
// instance.Status.Groups (the set previously applied by this operator), so
// memberships created outside of this resource are left untouched, while
// names removed from spec.groups have their membership revoked.
func (r *UserReconciler) reconcileGroups(ctx context.Context, instance *v1alpha1.User, session *KeycloakSession, keycloakID string) error {
	toAdd, toRemove := diffNames(instance.Spec.Groups, instance.Status.Groups)
	if len(toAdd) == 0 && len(toRemove) == 0 {
		return nil
	}

	for _, name := range toAdd {
		groupID, err := findGroupIDByName(ctx, session, name)
		if err != nil {
			return fmt.Errorf("failed to resolve group %q: %w", name, err)
		}
		if groupID == "" {
			return fmt.Errorf("group %q not found in realm %q", name, session.RealmName)
		}
		if err := session.Client.AddUserToGroup(ctx, session.Token.AccessToken, session.RealmName, keycloakID, groupID); err != nil {
			return fmt.Errorf("failed to add user to group %q: %w", name, err)
		}
	}

	for _, name := range toRemove {
		groupID, err := findGroupIDByName(ctx, session, name)
		if err != nil {
			return fmt.Errorf("failed to resolve group %q for removal: %w", name, err)
		}
		if groupID == "" {
			// Group no longer exists in Keycloak, nothing left to remove.
			continue
		}
		if err := session.Client.DeleteUserFromGroup(ctx, session.Token.AccessToken, session.RealmName, keycloakID, groupID); err != nil {
			return fmt.Errorf("failed to remove user from group %q: %w", name, err)
		}
	}

	return updateAppliedGroupsStatus(ctx, r.Client, instance)
}

// reconcileRealmRoles ensures the user's directly assigned realm roles match
// instance.Spec.RealmRoles. Assignment changes are computed relative to
// instance.Status.RealmRoles (the set previously applied by this operator),
// so roles assigned outside of this resource (e.g. default realm roles) are
// left untouched, while names removed from spec.realmRoles are revoked.
func (r *UserReconciler) reconcileRealmRoles(ctx context.Context, instance *v1alpha1.User, session *KeycloakSession, keycloakID string) error {
	toAdd, toRemove := diffNames(instance.Spec.RealmRoles, instance.Status.RealmRoles)
	if len(toAdd) == 0 && len(toRemove) == 0 {
		return nil
	}

	for _, name := range toAdd {
		role, err := session.Client.GetRealmRole(ctx, session.Token.AccessToken, session.RealmName, name)
		if err != nil {
			return fmt.Errorf("failed to resolve realm role %q: %w", name, err)
		}
		if err := session.Client.AddRealmRoleToUser(ctx, session.Token.AccessToken, session.RealmName, keycloakID, []gocloak.Role{*role}); err != nil {
			return fmt.Errorf("failed to assign realm role %q: %w", name, err)
		}
	}

	for _, name := range toRemove {
		role, err := session.Client.GetRealmRole(ctx, session.Token.AccessToken, session.RealmName, name)
		if err != nil {
			var apiErr *gocloak.APIError
			if errors.As(err, &apiErr) && apiErr.Code == 404 {
				// Role no longer exists in Keycloak, nothing left to remove.
				continue
			}
			return fmt.Errorf("failed to resolve realm role %q for removal: %w", name, err)
		}
		if err := session.Client.DeleteRealmRoleFromUser(ctx, session.Token.AccessToken, session.RealmName, keycloakID, []gocloak.Role{*role}); err != nil {
			return fmt.Errorf("failed to revoke realm role %q: %w", name, err)
		}
	}

	return updateAppliedRealmRolesStatus(ctx, r.Client, instance)
}

// findGroupIDByName searches the realm for a top-level group with an exact
// name match, returning an empty ID (without error) when no such group
// exists.
func findGroupIDByName(ctx context.Context, session *KeycloakSession, name string) (string, error) {
	groups, err := session.Client.GetGroups(ctx, session.Token.AccessToken, session.RealmName, gocloak.GetGroupsParams{
		Search: gocloak.StringP(name),
		Exact:  gocloak.BoolP(true),
	})
	if err != nil {
		return "", err
	}

	for _, g := range groups {
		if g.Name != nil && *g.Name == name && g.ID != nil {
			return *g.ID, nil
		}
	}

	return "", nil
}

// diffNames compares a desired set of names against the set previously
// applied by this operator (as recorded in status), returning the names that
// need to be added and the names that need to be removed.
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

// updateAppliedGroupsStatus persists instance.Spec.Groups into
// instance.Status.Groups, recording the set of memberships this operator has
// applied so that future reconciliations can compute an accurate diff.
func updateAppliedGroupsStatus(ctx context.Context, c client.Client, instance *v1alpha1.User) error {
	instance.Status.Groups = instance.Spec.Groups
	if err := c.Status().Update(ctx, instance); err != nil {
		return fmt.Errorf("failed to update applied groups status: %w", err)
	}
	return nil
}

// updateAppliedRealmRolesStatus persists instance.Spec.RealmRoles into
// instance.Status.RealmRoles, recording the set of role assignments this
// operator has applied so that future reconciliations can compute an
// accurate diff.
func updateAppliedRealmRolesStatus(ctx context.Context, c client.Client, instance *v1alpha1.User) error {
	instance.Status.RealmRoles = instance.Spec.RealmRoles
	if err := c.Status().Update(ctx, instance); err != nil {
		return fmt.Errorf("failed to update applied realm roles status: %w", err)
	}
	return nil
}

// setPassword resolves the spec.password reference, if set, and applies it to
// the Keycloak user via the admin API. It is a no-op when no password is
// configured, leaving the user without credentials until an administrator
// sets one manually.
func (r *UserReconciler) setPassword(ctx context.Context, instance *v1alpha1.User, session *KeycloakSession, keycloakID string) error {
	if instance.Spec.Password == nil {
		return nil
	}

	password, err := ResolveSecretKeyRefOrVal(ctx, r.Client, instance.Spec.Password, instance.Namespace)
	if err != nil {
		return fmt.Errorf("unable to resolve password: %w", err)
	}

	temporary := true
	if instance.Spec.Temporary != nil {
		temporary = *instance.Spec.Temporary
	}

	if err := session.Client.SetPassword(ctx, session.Token.AccessToken, keycloakID, session.RealmName, password, temporary); err != nil {
		return fmt.Errorf("failed to set password for user in Keycloak: %w", err)
	}

	return nil
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
		user.Attributes = attrs
	}

	if u.Spec.RequiredActions != nil {
		actions := u.Spec.RequiredActions
		user.RequiredActions = actions
	}

	return user
}

// SetupWithManager sets up the controller with the Manager.
func (r *UserReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := RegisterRefIndex(mgr, &v1alpha1.User{}, userSecretIndexField, userSecretRefKeys); err != nil {
		return err
	}

	newList := func() client.ObjectList { return &v1alpha1.UserList{} }

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.User{}).
		Watches(&corev1.Secret{}, RefEventHandler(mgr.GetClient(), newList, userSecretIndexField)).
		Named("user").
		Complete(r)
}

const userSecretIndexField = ".spec.secretRefs"

// userSecretRefKeys returns the RefIndexKey values for every Secret the
// given User instance may reference (currently the initial password).
func userSecretRefKeys(obj client.Object) []string {
	user, ok := obj.(*v1alpha1.User)
	if !ok || user.Spec.Password == nil || user.Spec.Password.SecretKeyRef == nil {
		return nil
	}

	ref := user.Spec.Password.SecretKeyRef
	ns := ref.Namespace
	if ns == "" {
		ns = user.Namespace
	}

	return []string{RefIndexKey(ns, ref.Name)}
}
