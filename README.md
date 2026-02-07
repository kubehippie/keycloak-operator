# keycloak-opereator

[![GitHub Repo](https://img.shields.io/badge/github-repo-yellowgreen)](https://github.com/kubehippie/keycloak-operator) [![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/kubehippie)](https://artifacthub.io/packages/helm/kubehippie/keycloak-operator)

This controller can configure an existing Keycloak instance. We don't wanted to
use Terraform or some JSON document for this anymore and since there haven't
been good controllers out there we built our own version of it.

## Instructions

Generally you should install this project via [Helm][helm], the other options
are not covered by this document as the chart deployment is the preferred way:

```sh
cat << EOF > values.yaml
fullnameOverride: keycloak-operartor
EOF

helm install keycloak-operator oci://ghcr.io/kubehippie/charts/keycloak-operator --values values.yaml
```

## Development

If you are not familiar with [Nix][nix] it is up to you to have a working
environment for Go (>= 1.24.0) as the setup won't be covered within this guide.
Please follow the official install instructions for [Go][golang] and. Beside
that we are using `make` to define all commands to build this project.

```console
git clone https://github.com/kubehippie/keycloak-operator.git
cd keycloak-operator

make build
./bin/manager -h
```

If you got [Nix][nix] and [Direnv][direnv] configured you can simply execute
the following commands to get all dependencies including `make` and the required
runtimes installed:

```console
cat << EOF > .envrc
use flake . --impure
EOF

direnv allow
```

To easily work on the operator we suggest to use [Tilt][tilt] for the local
development, this work pretty good in combination with Kind to get features like
hot reloading:

```console
kind create cluster \
    --name keycloak-operator

helm upgrade cert-manager cert-manager \
    --install --hide-notes --wait \
    --repo https://charts.jetstack.io \
    --set fullnameOverride=cert-manager \
    --set crds.enabled=true \
    --namespace cert-manager \
    --create-namespace

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
[nix]: https://nixos.org/
[golang]: http://golang.org/doc/install.html
[direnv]: https://direnv.net/
[tilt]: https://tilt.dev/
