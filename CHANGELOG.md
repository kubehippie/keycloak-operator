# Changelog

## [1.6.0](https://github.com/kubehippie/keycloak-operator/compare/v1.5.0...v1.6.0) (2026-09-10)

### Features

* integrate client and realm roles for service accounts ([8fb6cae](https://github.com/kubehippie/keycloak-operator/commit/8fb6caeffb1e8f8206c47b351fcf0e4de000c396))

### Bugfixes

* **deps:** update golang:1.27.1 docker digest to f44f6e8 ([#110](https://github.com/kubehippie/keycloak-operator/issues/110)) ([1e4beaf](https://github.com/kubehippie/keycloak-operator/commit/1e4beafd6094316f7e1f609bf0d39c6130c3facb))
* resolve typos or wrong wording ([84ca5c5](https://github.com/kubehippie/keycloak-operator/commit/84ca5c56dea501fd05e07d375db18634999e99b3))

### Dependencies

* **patch:** update module github.com/onsi/ginkgo/v2 to v2.32.2 ([#108](https://github.com/kubehippie/keycloak-operator/issues/108)) ([5184fff](https://github.com/kubehippie/keycloak-operator/commit/5184fffe8ba93ebac414bccf91cb136cac0640be))

## [1.5.0](https://github.com/kubehippie/keycloak-operator/compare/v1.4.3...v1.5.0) (2026-09-09)

### Bugfixes

* correctly handle reconciles and keep values intact like login theme ([4616b47](https://github.com/kubehippie/keycloak-operator/commit/4616b47494fa75fa5fc43a66123bae3e27d8c499))
* **deps:** update gcr.io/distroless/static:nonroot docker digest to 1c2c046 ([#88](https://github.com/kubehippie/keycloak-operator/issues/88)) ([1ba638b](https://github.com/kubehippie/keycloak-operator/commit/1ba638b8e4ab60b54c01d96ab5636de2b203ad4a))
* **deps:** update golang:1.27.0 docker digest to 0ecdc2a ([#91](https://github.com/kubehippie/keycloak-operator/issues/91)) ([f31762b](https://github.com/kubehippie/keycloak-operator/commit/f31762b017d42dea3892c1da97de01cec8630103))
* **deps:** update golang:1.27.0 docker digest to 4013ae0 ([#98](https://github.com/kubehippie/keycloak-operator/issues/98)) ([204bdd9](https://github.com/kubehippie/keycloak-operator/commit/204bdd994d17901fe19cc44df6bcca787683bcde))
* **deps:** update golang:1.27.0 docker digest to 713ecb4 ([#89](https://github.com/kubehippie/keycloak-operator/issues/89)) ([64afbc4](https://github.com/kubehippie/keycloak-operator/commit/64afbc4141dcc374f9627e04e5a1990cd86a4a50))
* **deps:** update golang:1.27.0 docker digest to f42f854 ([#90](https://github.com/kubehippie/keycloak-operator/issues/90)) ([5072e26](https://github.com/kubehippie/keycloak-operator/commit/5072e26edb0d243eee58200243bd1ac36945f8a5))
* **deps:** update golang:1.27.1 docker digest to 512690a ([#105](https://github.com/kubehippie/keycloak-operator/issues/105)) ([ded5129](https://github.com/kubehippie/keycloak-operator/commit/ded5129fe62881e081d268da6ef105b64e202072))
* replace deprecated Requeue by RequeueAfter ([4019a74](https://github.com/kubehippie/keycloak-operator/commit/4019a74a6f47950fbf0a43df0b905458711a950d))
* split realm controller to spread the complexity ([e9f3e43](https://github.com/kubehippie/keycloak-operator/commit/e9f3e43a0fc9d09aa49cca1326054c36f7a1b9a7))

### Dependencies

* **minor:** update dependency go to v1.27.0 ([#80](https://github.com/kubehippie/keycloak-operator/issues/80)) ([5c45a81](https://github.com/kubehippie/keycloak-operator/commit/5c45a81b51d73f61d67299a4d7c94806371d6f0d))
* **minor:** update dependency golangci/golangci-lint to v2.13.0 ([#83](https://github.com/kubehippie/keycloak-operator/issues/83)) ([98ace22](https://github.com/kubehippie/keycloak-operator/commit/98ace22f5dbd060825f8d582bd3fc5bcf04ced01))
* **minor:** update dependency kind to v0.33.0 ([#92](https://github.com/kubehippie/keycloak-operator/issues/92)) ([3f664fd](https://github.com/kubehippie/keycloak-operator/commit/3f664fd7084fcb29acef812fa63187f1dabe9fa2))
* **minor:** update dependency kubectl to v1.37.0 ([#93](https://github.com/kubehippie/keycloak-operator/issues/93)) ([b4348ec](https://github.com/kubehippie/keycloak-operator/commit/b4348ec20fef676778a65374676b84d87be76796))
* **minor:** update dependency kubernetes-sigs/controller-tools to v0.22.0 ([#104](https://github.com/kubehippie/keycloak-operator/issues/104)) ([f163255](https://github.com/kubehippie/keycloak-operator/commit/f163255a19f836b0cdbb19967e74927f477529c4))
* **minor:** update dependency pipx to v1.17.1 ([#99](https://github.com/kubehippie/keycloak-operator/issues/99)) ([3b2e503](https://github.com/kubehippie/keycloak-operator/commit/3b2e5037fbe4fd25deca33c4295261901d6cd3bb))
* **minor:** update golang docker tag to v1.27.0 ([#81](https://github.com/kubehippie/keycloak-operator/issues/81)) ([73fd884](https://github.com/kubehippie/keycloak-operator/commit/73fd884a3289b24f30d89bf1a74e07b718e73c01))
* **minor:** update kubernetes monorepo to v0.37.0 ([#94](https://github.com/kubehippie/keycloak-operator/issues/94)) ([ea89592](https://github.com/kubehippie/keycloak-operator/commit/ea8959201c388fcf20b211945897b8aee8b613ed))
* **minor:** update module github.com/onsi/gomega to v1.43.0 ([#95](https://github.com/kubehippie/keycloak-operator/issues/95)) ([763f4f4](https://github.com/kubehippie/keycloak-operator/commit/763f4f4b64250a46bbf60e3e6aff061cd02c2921))
* **minor:** update module sigs.k8s.io/controller-runtime to v0.25.0 ([#106](https://github.com/kubehippie/keycloak-operator/issues/106)) ([b35f8a4](https://github.com/kubehippie/keycloak-operator/commit/b35f8a4c1fe4e12e4ba2ccf4a313e3c45119df84))
* **patch:** update dependency go to v1.26.7 ([#79](https://github.com/kubehippie/keycloak-operator/issues/79)) ([26ff5cd](https://github.com/kubehippie/keycloak-operator/commit/26ff5cdae776283c56f2c099bcd2bb68e13e5bee))
* **patch:** update dependency go to v1.27.1 ([#100](https://github.com/kubehippie/keycloak-operator/issues/100)) ([a1ccad1](https://github.com/kubehippie/keycloak-operator/commit/a1ccad172250e0913ee6664b65d071e64f650a91))
* **patch:** update dependency golangci/golangci-lint to v2.13.1 ([#86](https://github.com/kubehippie/keycloak-operator/issues/86)) ([8465b43](https://github.com/kubehippie/keycloak-operator/commit/8465b43388b300922ebb58bf82409a57d546740c))
* **patch:** update dependency golangci/golangci-lint to v2.13.2 ([#97](https://github.com/kubehippie/keycloak-operator/issues/97)) ([85a0c67](https://github.com/kubehippie/keycloak-operator/commit/85a0c67c8c700ecbbabd5b6651cd2d6ae09e2dfe))
* **patch:** update dependency kubectl to v1.36.4 ([#85](https://github.com/kubehippie/keycloak-operator/issues/85)) ([8cbbc40](https://github.com/kubehippie/keycloak-operator/commit/8cbbc40b3fdab6b7670325782b7d325ede1fdea2))
* **patch:** update dependency pipx to v1.17.2 ([#101](https://github.com/kubehippie/keycloak-operator/issues/101)) ([063aa9e](https://github.com/kubehippie/keycloak-operator/commit/063aa9e64ad61e42dfeb598d465884de51813a40))
* **patch:** update dependency yq to v4.53.4 ([#78](https://github.com/kubehippie/keycloak-operator/issues/78)) ([7459bbe](https://github.com/kubehippie/keycloak-operator/commit/7459bbebc351dd28456946c37be2096843f3498c))
* **patch:** update dependency yq to v4.53.6 ([#84](https://github.com/kubehippie/keycloak-operator/issues/84)) ([d3134d3](https://github.com/kubehippie/keycloak-operator/commit/d3134d30066f8e0b0936f0b47e77047c0ab21753))
* **patch:** update golang docker tag to v1.27.1 ([#102](https://github.com/kubehippie/keycloak-operator/issues/102)) ([7983b0e](https://github.com/kubehippie/keycloak-operator/commit/7983b0efeccd6a81d87a9478fd4d07ccbf5854e8))
* **patch:** update kubernetes monorepo to v0.36.4 ([#87](https://github.com/kubehippie/keycloak-operator/issues/87)) ([4d140d1](https://github.com/kubehippie/keycloak-operator/commit/4d140d1d1852084fdf6226c99f5b20873dca606b))

## [1.4.3](https://github.com/kubehippie/keycloak-operator/compare/v1.4.2...v1.4.3) (2026-08-17)

## [1.4.2](https://github.com/kubehippie/keycloak-operator/compare/v1.4.1...v1.4.2) (2026-08-10)

## [1.4.1](https://github.com/kubehippie/keycloak-operator/compare/v1.4.0...v1.4.1) (2026-08-03)

## [1.4.0](https://github.com/kubehippie/keycloak-operator/compare/v1.3.1...v1.4.0) (2026-07-27)

## [1.3.1](https://github.com/kubehippie/keycloak-operator/compare/v1.3.0...v1.3.1) (2026-07-21)

## [1.3.0](https://github.com/kubehippie/keycloak-operator/compare/v1.2.1...v1.3.0) (2026-07-21)

## [1.2.1](https://github.com/kubehippie/keycloak-operator/compare/v1.2.0...v1.2.1) (2026-07-20)

## [1.2.0](https://github.com/kubehippie/keycloak-operator/compare/v1.1.0...v1.2.0) (2026-07-16)

## [1.1.0](https://github.com/kubehippie/keycloak-operator/compare/v1.0.2...v1.1.0) (2026-07-15)

## [1.0.2](https://github.com/kubehippie/keycloak-operator/compare/v1.0.1...v1.0.2) (2026-07-14)

## [1.0.1](https://github.com/kubehippie/keycloak-operator/compare/v1.0.0...v1.0.1) (2026-07-14)

## 1.0.0 (2026-07-14)
