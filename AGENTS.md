# keycloak-operator — Agent Guide

This document provides guidance for AI coding agents (GitHub Copilot, Codex,
Claude, etc.) working in this repository.

---

## Git Workflow — MANDATORY

- **Never commit changes** unless explicitly instructed to do so.
- **Never create a branch** unless explicitly instructed to do so.
- **Never open a pull request** unless explicitly instructed to do so.
- Leave all changes as unstaged working-tree modifications by default.

---

## Project Overview

`keycloak-operator` is a **Kubernetes operator** written in Go using
[controller-runtime](https://github.com/kubernetes-sigs/controller-runtime)
(kubebuilder pattern). It manages Keycloak realms, clients, users, groups, and
identity provider configuration declaratively via Kubernetes custom resources.

The operator connects to one or more Keycloak instances (referenced via
`Keycloak` resources) and reconciles the desired state expressed in CRDs into
the live Keycloak API using
[gocloak](https://github.com/Nerzal/gocloak).

---

## Repository Layout

```text
cmd/
  main.go                        # Entrypoint — sets up the controller-runtime manager

api/
  common/
    ref.go                       # Shared types: KeycloakRef, SecretKeyRefOrVal, SourceRef, …
  v1alpha1/
    keycloak_types.go            # Keycloak (namespaced connection resource)
    realm_types.go               # Realm
    user_types.go                # User
    group_types.go               # Group
    groupversion_info.go
    zz_generated.deepcopy.go     # Auto-generated — do not edit
  identity/v1alpha1/
    oidcidentityprovider_types.go
    attributeimportermapper_types.go
    oidcadvancedgroupmapper_types.go
    hardcodedrolemapper_types.go
    groupversion_info.go
    zz_generated.deepcopy.go
  openid/v1alpha1/
    openidclient_types.go
    clientscope_types.go
    defaultscopes_types.go
    groupmembershipprotocolmapper_types.go
    userpropertyprotocolmapper_types.go
    userattributeprotocolmapper_types.go
    fullnameprotocolmapper_types.go
    audienceprotocolmapper_types.go
    groupversion_info.go
    zz_generated.deepcopy.go

internal/
  controller/
    keycloak_controller.go       # Verifies connectivity and updates .status.connected
    realm_controller.go
    user_controller.go
    group_controller.go
    controller_permissions.go    # Shared RBAC markers for secrets/configmaps
    identity/                    # Controllers for identity provider resources
    openid/                      # Controllers for OpenID client resources
  webhook/
    v1alpha1/                    # Defaulting + validating webhooks for core types
    identity/v1alpha1/           # Webhooks for identity types
    openid/v1alpha1/             # Webhooks for openid types

config/
  crd/bases/                     # Generated CRD YAML — do not edit manually
  rbac/                          # Generated RBAC roles — do not edit manually
  webhook/                       # Webhook manifests — do not edit manually
  samples/                       # Example CRs for each type
  default/                       # Kustomize overlay used for deployment
  manager/

chart/                           # Helm chart (generated via kubebuilder helm plugin)

test/
  e2e/                           # End-to-end tests using Ginkgo + Kind
  utils/                         # Shared test helpers

grafana/                         # Grafana dashboard definitions

dist/                            # Generated dist/install.yaml (kustomize build output)

mise.toml                         # mise tool version definitions (dev shell)
Makefile                         # All build, test, lint, and deploy targets
PROJECT                          # kubebuilder project metadata
```

---

## Shell Environment

This repository uses [mise](https://mise.jdx.dev/) (`mise.toml`) to manage
tool versions for the dev shell (Go, Helm, kubebuilder, kind, kubectl, etc.).

To activate the tools:

```bash
mise install
mise use
```

---

## API Groups & CRD Types

### `keycloak-operator.webhippie.de/v1alpha1` (core)

| Kind | Scope | Description |
|------|-------|-------------|
| `Keycloak` | Namespaced | Connection to a Keycloak instance (URL, credentials) |
| `Realm` | Namespaced | A Keycloak realm managed by this operator |
| `User` | Namespaced | A user within a realm |
| `Group` | Namespaced | A group within a realm |

### `identity.keycloak-operator.webhippie.de/v1alpha1`

| Kind | Scope | Description |
|------|-------|-------------|
| `OIDCIdentityProvider` | Namespaced | OIDC identity provider on a realm |
| `AttributeImporterMapper` | Namespaced | Attribute importer mapper for an IdP |
| `OIDCAdvancedGroupMapper` | Namespaced | Advanced group mapper for an IdP |
| `HardcodedRoleMapper` | Namespaced | Hardcoded role mapper for an IdP |

### `openid.keycloak-operator.webhippie.de/v1alpha1`

| Kind | Scope | Description |
|------|-------|-------------|
| `OpenIDClient` | Namespaced | An OpenID Connect client on a realm |
| `ClientScope` | Namespaced | A client scope |
| `DefaultScopes` | Namespaced | Default client scopes assignment |
| `GroupMembershipProtocolMapper` | Namespaced | Group membership protocol mapper |
| `UserPropertyProtocolMapper` | Namespaced | User property protocol mapper |
| `UserAttributeProtocolMapper` | Namespaced | User attribute protocol mapper |
| `FullNameProtocolMapper` | Namespaced | Full name protocol mapper |
| `AudienceProtocolMapper` | Namespaced | Audience protocol mapper |

---

## Common Types (`api/common`)

| Type | Description |
|------|-------------|
| `KeycloakRef` | Reference to a `Keycloak` resource (`kind`, `name`, `namespace`) |
| `SecretKeyRefOrVal` | Secret key reference or inline value (username/password fields) |
| `SourceRef` | ConfigMap or Secret key reference (used for `caCert`) |
| `SourceRefOrVal` | `SourceRef` + inline value |

`KeycloakRef.Kind` is always `Keycloak` (enum-validated by kubebuilder).

---

## Build & Development Commands

All commands are defined in `Makefile`. Key targets:

| Target | Description |
|--------|-------------|
| `make build` | Build the `bin/manager` binary |
| `make run` | Run the controller locally (requires kubeconfig) |
| `make test` | Run unit/integration tests via envtest |
| `make test-e2e` | Run e2e tests against a temporary Kind cluster |
| `make lint` | Run `golangci-lint` |
| `make lint-fix` | Run linter with auto-fix |
| `make manifests` | Regenerate RBAC/CRD/webhook manifests via `controller-gen` |
| `make generate` | Regenerate DeepCopy methods via `controller-gen` |
| `make fmt` | Run `go fmt ./...` |
| `make vet` | Run `go vet ./...` |
| `make docker-build` | Build the container image (`ghcr.io/kubehippie/keycloak-operator`) |
| `make build-installer` | Generate `dist/install.yaml` via kustomize |

All tooling (`controller-gen`, `kustomize`, `golangci-lint`, `setup-envtest`)
is downloaded locally into `bin/` by the Makefile — no manual installation
needed.

After changing any `_types.go` file always run:

```bash
make generate manifests
```

---

## Generated Files — Do Not Edit

The following files are fully generated and must not be edited by hand:

- `api/*/zz_generated.deepcopy.go` — regenerated by `make generate`
- `config/crd/bases/*.yaml` — regenerated by `make manifests`
- `config/rbac/role.yaml` and the per-resource `*_admin/editor/viewer_role.yaml` files — regenerated by `make manifests`
- `config/webhook/manifests.yaml` — regenerated by `make manifests`

---

## Testing

### Unit / integration tests

Uses [Ginkgo](https://onsi.github.io/ginkgo/) +
[Gomega](https://onsi.github.io/gomega/) with `controller-runtime`'s
`envtest` (real API server, no full cluster).

```bash
make test
```

### End-to-end tests

Uses a [Kind](https://kind.sigs.k8s.io/) cluster spun up automatically.

```bash
make test-e2e
```

The Kind cluster is named `keycloak-operator-test-e2e` and is torn down after
the test run (`make cleanup-test-e2e`).

---

## CI Workflows

| Workflow | File | Purpose |
|----------|------|---------|
| General | `.github/workflows/general.yml` | Build, test, lint |
| Docker | `.github/workflows/docker.yml` | Build and push container image |
| Release | `.github/workflows/release.yml` | Semantic release and changelog |
| Helm docs | `.github/workflows/helmdocs.yml` | Regenerate chart documentation |
| Automerge | `.github/workflows/automerge.yml` | Renovate automation |

Required status checks on `master`: `lint`, `chart`, `tests`, `docker`.

---

## Deployment

The preferred installation method is the Helm chart:

```bash
helm install keycloak-operator \
  oci://ghcr.io/kubehippie/charts/keycloak-operator \
  --values values.yaml
```

Raw Kustomize manifests:

```bash
make build-installer   # outputs dist/install.yaml
```

---

## Contribution Conventions

- Use pull requests for all changes; squash or rebase merge only.
- Run `make fmt vet lint` before pushing.
- After any `_types.go` change, run `make generate manifests` and commit the
  updated generated files together with the type change.
- Keep the Helm chart in sync with CRD/config changes (`make chart`).
- For security issues, contact `thomas@webhippie.de`.
