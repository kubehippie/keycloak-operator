# keycloak-operator

[![GitHub Repo](https://img.shields.io/badge/github-repo-yellowgreen)](https://github.com/kubehippie/keycloak-operator) [![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/kubehippie)](https://artifacthub.io/packages/helm/kubehippie/keycloak-operator)

> [!WARNING]
> **This project is in early development.** The builtin resources and their APIs
> are **not yet stable** and can introduce **breaking changes** at any time,
> including changes to the CRD schemas, defaulting/validation behavior, and
> reconciliation logic. Pin your version, review the changelog before upgrading,
> and do not use this operator for production workloads until it reaches a
> stable release.

This controller can configure an existing Keycloak instance. We don't wanted to
use Terraform or some JSON document for this anymore and since there haven't
been good controllers out there we built our own version of it.

## Instructions

Generally you should install this project via [Helm][helm], the other options
are not covered by this document as the chart deployment is the preferred way:

```sh
cat << EOF | helm install keycloak-operator oci://ghcr.io/kubehippie/charts/keycloak-operator --values -
fullnameOverride: keycloak-operator
EOF
```

## Development

We are using [Mise][mise] to install all required tools with fixed versions to
keep everything as far as possible compatible. If you don't want to use
[Mise][mise] it is up to you to install the required tools like Go. Beside that
we are using `make` to define all commands to build this project.

```console
git clone https://github.com/kubehippie/keycloak-operator.git
cd keycloak-operator

mise trust
mise install

make build
./bin/manager -h
```

To easily work on the operator we suggest to use [Tilt][tilt] for the local
development, this work pretty good in combination with Kind to get features like
hot reloading:

```console
kind create cluster \
    --name keycloak-operator

tilt up

kind delete cluster \
    --name keycloak-operator
```

## Security

If you find a security issue please contact
[thomas@webhippie.de](mailto:thomas@webhippie.de) first.

## Contributing

Fork -> Patch -> Push -> Pull Request

## Authors

-   [Thomas Boerger](https://github.com/tboerger)

## License

Apache-2.0

## Copyright

```console
Copyright (c) 2025 Thomas Boerger <thomas@webhippie.de>
```

[helm]: https://helm.sh/
[mise]: https://mise.jdx.dev/getting-started.html
[tilt]: https://tilt.dev/
