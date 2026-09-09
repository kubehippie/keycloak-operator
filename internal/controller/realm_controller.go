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
	"strconv"
	"time"

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

const realmFinalizer = "keycloak-operator.webhippie.de/realm"

// RealmReconciler reconciles a Realm object
type RealmReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=keycloak-operator.webhippie.de,resources=realms,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=keycloak-operator.webhippie.de,resources=realms/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=keycloak-operator.webhippie.de,resources=realms/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *RealmReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)
	log.Info("Reconciling")

	instance := &v1alpha1.Realm{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("unable to fetch: %w", err)
	}

	session, err := KeycloakSessionForKeycloak(ctx, r.Client, instance.Spec.KeycloakRef, req.Namespace)
	if err != nil {
		if !instance.DeletionTimestamp.IsZero() && apierrors.IsNotFound(err) {
			log.Info("Keycloak no longer exists, skipping Keycloak cleanup", "error", err.Error())
			return r.handleDeletion(ctx, instance, nil)
		}
		log.Error(err, "Unable to get Keycloak session")
		return ctrl.Result{RequeueAfter: FailedKeycloakConnectionRetryPeriod}, nil
	}

	if !instance.DeletionTimestamp.IsZero() {
		return r.handleDeletion(ctx, instance, session)
	}

	if !controllerutil.ContainsFinalizer(instance, realmFinalizer) {
		controllerutil.AddFinalizer(instance, realmFinalizer)
		if err := r.Update(ctx, instance); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to add finalizer: %w", err)
		}
		return ctrl.Result{RequeueAfter: time.Second}, nil
	}

	return r.reconcileRealm(ctx, instance, session)
}

func (r *RealmReconciler) handleDeletion(ctx context.Context, instance *v1alpha1.Realm, session *KeycloakSession) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	if !controllerutil.ContainsFinalizer(instance, realmFinalizer) {
		return ctrl.Result{}, nil
	}

	if session != nil {
		log.Info("Deleting realm from Keycloak", "realm", instance.Spec.Name)
		if err := session.Client.DeleteRealm(ctx, session.Token.AccessToken, instance.Spec.Name); err != nil {
			var apiErr *gocloak.APIError
			if errors.As(err, &apiErr) && apiErr.Code == 404 {
				log.Info("Realm already absent in Keycloak, skipping delete")
			} else {
				return ctrl.Result{}, fmt.Errorf("failed to delete realm from Keycloak: %w", err)
			}
		}
	}

	controllerutil.RemoveFinalizer(instance, realmFinalizer)
	if err := r.Update(ctx, instance); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to remove finalizer: %w", err)
	}

	return ctrl.Result{}, nil
}

func (r *RealmReconciler) reconcileRealm(ctx context.Context, instance *v1alpha1.Realm, session *KeycloakSession) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	existing, err := session.Client.GetRealm(ctx, session.Token.AccessToken, instance.Spec.Name)
	if err != nil {
		var apiErr *gocloak.APIError
		if !errors.As(err, &apiErr) || apiErr.Code != 404 {
			return ctrl.Result{}, fmt.Errorf("failed to check for existing realm: %w", err)
		}

		desired, err := realmToGocloak(ctx, r.Client, instance, instance.Namespace, nil)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to build realm representation: %w", err)
		}

		id, err := session.Client.CreateRealm(ctx, session.Token.AccessToken, desired)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to create realm in Keycloak: %w", err)
		}
		log.Info("Realm created in Keycloak", "id", id)
		if err := UpdateKeycloakIDStatus(ctx, r.Client, instance, &id); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	if err := UpdateKeycloakIDStatus(ctx, r.Client, instance, existing.ID); err != nil {
		return ctrl.Result{}, err
	}

	// Build the desired representation on top of the realm's current state in
	// Keycloak. Keycloak's realm update endpoint treats the PUT body as the
	// full representation: any field left nil/omitted is reset to its
	// default rather than left untouched. Starting from `existing` and only
	// overlaying fields explicitly managed by the CR spec preserves anything
	// not modeled/set here (e.g. settings changed outside the operator).
	desired, err := realmToGocloak(ctx, r.Client, instance, instance.Namespace, existing)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to build realm representation: %w", err)
	}

	if err := session.Client.UpdateRealm(ctx, session.Token.AccessToken, desired); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to update realm in Keycloak: %w", err)
	}

	log.Info("Realm reconciled", "realm", instance.Spec.Name)
	return ctrl.Result{}, nil
}

// realmToGocloak converts a Realm CR spec into the gocloak representation
// used for create and update API calls. When existing is non-nil (i.e. the
// realm already exists in Keycloak), it is used as the base representation
// so that fields not managed by the CR spec are preserved rather than reset
// to their Keycloak defaults by the full-representation PUT semantics of the
// realm update endpoint. The SMTP password (if configured as a Secret
// reference) is resolved from the cluster at call time.
func realmToGocloak(ctx context.Context, cl client.Client, r *v1alpha1.Realm, ns string, existing *gocloak.RealmRepresentation) (gocloak.RealmRepresentation, error) {
	var realm gocloak.RealmRepresentation
	if existing != nil {
		realm = *existing
	}

	realm.Realm = gocloak.StringP(r.Spec.Name)

	applyRealmBasicFields(&realm, r)
	applyRealmLoginSettings(&realm, r.Spec.Login)
	applyRealmThemes(&realm, r.Spec.Themes)
	applyRealmInternationalization(&realm, r.Spec.Internationalization)
	applyRealmSessionTimeouts(&realm, r.Spec.SessionTimeouts)

	if err := applyRealmSmtpServer(ctx, cl, &realm, r.Spec.SmtpServer, ns); err != nil {
		return gocloak.RealmRepresentation{}, err
	}

	if r.Spec.Attributes != nil {
		realm.Attributes = r.Spec.Attributes
	}

	return realm, nil
}

// applyRealmBasicFields overlays the top-level, non-grouped Realm spec
// fields onto the given representation.
func applyRealmBasicFields(realm *gocloak.RealmRepresentation, r *v1alpha1.Realm) {
	if r.Spec.DisplayName != nil {
		realm.DisplayName = r.Spec.DisplayName
	}
	if r.Spec.DisplayNameHtml != nil {
		realm.DisplayNameHTML = r.Spec.DisplayNameHtml
	}
	if r.Spec.Enabled != nil {
		realm.Enabled = r.Spec.Enabled
	}
	if r.Spec.SslRequired != nil {
		realm.SSLRequired = r.Spec.SslRequired
	}
	if r.Spec.PasswordPolicy != nil {
		realm.PasswordPolicy = r.Spec.PasswordPolicy
	}
	if r.Spec.BruteForceProtected != nil {
		realm.BruteForceProtected = r.Spec.BruteForceProtected
	}
}

// applyRealmLoginSettings overlays the Login spec section onto the given
// representation.
func applyRealmLoginSettings(realm *gocloak.RealmRepresentation, l *v1alpha1.RealmLoginSettings) {
	if l == nil {
		return
	}

	if l.RegistrationAllowed != nil {
		realm.RegistrationAllowed = l.RegistrationAllowed
	}
	if l.RegistrationEmailAsUsername != nil {
		realm.RegistrationEmailAsUsername = l.RegistrationEmailAsUsername
	}
	if l.EditUsernameAllowed != nil {
		realm.EditUsernameAllowed = l.EditUsernameAllowed
	}
	if l.ResetPasswordAllowed != nil {
		realm.ResetPasswordAllowed = l.ResetPasswordAllowed
	}
	if l.RememberMe != nil {
		realm.RememberMe = l.RememberMe
	}
	if l.VerifyEmail != nil {
		realm.VerifyEmail = l.VerifyEmail
	}
	if l.LoginWithEmailAllowed != nil {
		realm.LoginWithEmailAllowed = l.LoginWithEmailAllowed
	}
	if l.DuplicateEmailsAllowed != nil {
		realm.DuplicateEmailsAllowed = l.DuplicateEmailsAllowed
	}
}

// applyRealmThemes overlays the Themes spec section onto the given
// representation.
func applyRealmThemes(realm *gocloak.RealmRepresentation, t *v1alpha1.RealmThemes) {
	if t == nil {
		return
	}

	if t.Login != nil {
		realm.LoginTheme = t.Login
	}
	if t.Account != nil {
		realm.AccountTheme = t.Account
	}
	if t.Admin != nil {
		realm.AdminTheme = t.Admin
	}
	if t.Email != nil {
		realm.EmailTheme = t.Email
	}
}

// applyRealmInternationalization overlays the Internationalization spec
// section onto the given representation.
func applyRealmInternationalization(realm *gocloak.RealmRepresentation, i *v1alpha1.RealmInternationalization) {
	if i == nil {
		return
	}

	if i.Enabled != nil {
		realm.InternationalizationEnabled = i.Enabled
	}
	if i.DefaultLocale != nil {
		realm.DefaultLocale = i.DefaultLocale
	}
	if len(i.SupportedLocales) > 0 {
		realm.SupportedLocales = i.SupportedLocales
	}
}

// applyRealmSessionTimeouts overlays the SessionTimeouts spec section onto
// the given representation.
func applyRealmSessionTimeouts(realm *gocloak.RealmRepresentation, st *v1alpha1.RealmSessionTimeouts) {
	if st == nil {
		return
	}

	if st.AccessTokenLifespan != nil {
		realm.AccessTokenLifespan = st.AccessTokenLifespan
	}
	if st.SsoSessionIdleTimeout != nil {
		realm.SSOSessionIdleTimeout = st.SsoSessionIdleTimeout
	}
	if st.SsoSessionMaxLifespan != nil {
		realm.SSOSessionMaxLifespan = st.SsoSessionMaxLifespan
	}
	if st.OfflineSessionIdleTimeout != nil {
		realm.OfflineSessionIdleTimeout = st.OfflineSessionIdleTimeout
	}
}

// applyRealmSmtpServer builds the SMTP server settings map from the SmtpServer
// spec section, resolving the password from a Secret reference if needed,
// and assigns it onto the given representation.
func applyRealmSmtpServer(ctx context.Context, cl client.Client, realm *gocloak.RealmRepresentation, smtp *v1alpha1.RealmSmtpServer, ns string) error {
	if smtp == nil {
		return nil
	}

	smtpMap := map[string]string{
		"host": smtp.Host,
		"from": smtp.From,
	}

	if smtp.Port != nil {
		smtpMap["port"] = fmt.Sprintf("%d", *smtp.Port)
	}
	if smtp.FromDisplayName != nil {
		smtpMap["fromDisplayName"] = *smtp.FromDisplayName
	}
	if smtp.ReplyTo != nil {
		smtpMap["replyTo"] = *smtp.ReplyTo
	}
	if smtp.ReplyToDisplayName != nil {
		smtpMap["replyToDisplayName"] = *smtp.ReplyToDisplayName
	}
	if smtp.EnvelopeFrom != nil {
		smtpMap["envelopeFrom"] = *smtp.EnvelopeFrom
	}
	if smtp.Ssl != nil {
		smtpMap["ssl"] = strconv.FormatBool(*smtp.Ssl)
	}
	if smtp.StartTls != nil {
		smtpMap["starttls"] = strconv.FormatBool(*smtp.StartTls)
	}
	if smtp.Auth != nil {
		smtpMap["auth"] = strconv.FormatBool(*smtp.Auth)
	}
	if smtp.User != nil {
		smtpMap["user"] = *smtp.User
	}
	if smtp.Password != nil {
		password, err := ResolveSecretKeyRefOrVal(ctx, cl, smtp.Password, ns)
		if err != nil {
			return fmt.Errorf("failed to resolve SMTP password: %w", err)
		}
		smtpMap["password"] = password
	}

	realm.SMTPServer = smtpMap
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RealmReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := RegisterRefIndex(mgr, &v1alpha1.Realm{}, realmSecretIndexField, realmSecretRefKeys); err != nil {
		return err
	}

	newList := func() client.ObjectList { return &v1alpha1.RealmList{} }

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Realm{}).
		Watches(&corev1.Secret{}, RefEventHandler(mgr.GetClient(), newList, realmSecretIndexField)).
		Named("realm").
		Complete(r)
}

const realmSecretIndexField = ".spec.secretRefs"

// realmSecretRefKeys returns the RefIndexKey values for every Secret the
// given Realm instance may reference (currently the SMTP password).
func realmSecretRefKeys(obj client.Object) []string {
	realm, ok := obj.(*v1alpha1.Realm)
	if !ok || realm.Spec.SmtpServer == nil || realm.Spec.SmtpServer.Password == nil {
		return nil
	}

	ref := realm.Spec.SmtpServer.Password.SecretKeyRef
	if ref == nil {
		return nil
	}

	ns := ref.Namespace
	if ns == "" {
		ns = realm.Namespace
	}

	return []string{RefIndexKey(ns, ref.Name)}
}
