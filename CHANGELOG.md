# k8s-blueprint-operator Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Fixed
- [#117] configuring operator log-level

### Added
- [#117] add metadata mapping for logLevel

## [v2.6.0] - 2025-06-05

### Added
- [#113] blueprints now support additionalMounts for configMaps and secrets in dogus
- [#113] additionalMounts now get validated with explanations for possible problems

### Changed
- [#114] Dependency: minimum version of `k8s-dogu-operator-crd` is now 2.8.0
  - because we now set the new `minDataVolumeSize` field instead of `dataVolumeSize`
- [#113] Dependency: minimum version of `k8s-Blueprint-operator-crd` is now 1.3.0
- because we now need the additional mounts fields for dogus
- [#113] update to go 1.24.3
- [#113] extract Blueprint-CRD into bluerpint-lib
- [#113] extract Blueprint-k8s rest client into blueprint-lib
- [#113] update k8s.io/client-go to v0.33.0
- [#113] Moved blueprint format docs to [crd repository](https://github.com/cloudogu/k8s-blueprint-lib)

## [v2.5.0] - 2025-04-22
### Changed
- [#110] Set resource requests and limits to sensible defaults

## [v2.4.1] - 2025-04-17
### Changed
- [#109] Use retry for upfront health checks
   - blueprints will now not fail immediately if a dogu or component is not yet healthy, no workaround needed anymore

## [v2.4.0] - 2025-03-31
### Added
- [#107] Add additional print columns and aliases to CRDs

## [v2.3.1] - 2025-02-26
### Changed
- [#103] Extract blueprint structs and associated functionality and move it into blueprint-lib 
  - only that code was extracted that handles JSON parsing and superficial conversion in order to enable other modules to parse different blueprint versions.

### Fixed
- Fix autogeneration boilerplate text to match with the current make target `generate-deepcopy`

### Security
- Fix CVE-2024-45338 from golang.org/x/net

## [v2.3.0] - 2025-01-27
### Added
- [#105] Proxy support for the dogu registry. The proxy will be used from the secret `ces-proxy` which will be created by the setup.

## [v2.2.2] - 2024-12-19
### Fixed
- [#101] Fix CVE-2024-45337

## [v2.2.1] - 2024-12-18
### Fixed
- [#99] Service account creation fails because of dogu restarts
  - before restarting, we now wait for all dogus to get healthy

## [v2.2.0] - 2024-12-05
### Added
- [#97] Add a `deny-all` network-policy, to block all incoming traffic

### Removed
- [#97] Remove RBAC-Proxy along with k8s-metrics-service, because metrics are currently no used and all incoming traffic is blocked by the network-policy
- [#97] Remove unused WebHookServer

## [v2.1.1] - 2024-11-28
### Fixed
- [#95] Fix a bug of the dogu config state diff where multiple dogus replaced the whole diff.

## [v2.1.0] - 2024-11-22
### Changed
- [#87] Use ces-commons-lib for common errors and common types
- [#87] Use remote-dogu-descriptor-lib
- [#87] Use retry-lib
- [#93] deactivate operator leader election

### Security
- [#93] Remove RBAC permissions that seem unnecessary for the execution of the operator
  - this is an operational security measure

## [v2.0.1] - 2024-11-06
### Fixed
- [#81] Forbid component downgrades because the component operator can't handle this operation.
- [#81] Remove dogu configuration from removed dogus via the blueprint mask.
- [#81] Refactor DoguConfigDiff and remove code duplication for sensitive dogu config

## [v2.0.0] - 2024-10-29
### Changed
- Update module to v2
- [#85] Make imagePullSecrets configurable via helm values and use `ces-container-registries` as default.
- [#81] migrate etcd access to ecosystem config to k8s-config
- [#81] create configmaps and secrets for dogu config if the dogu is not yet installed
- [#81] give operator permissions to see configmaps and secrets
- [#81] use maintenance mode implementation from k8s-registry-lib
- [#81] use dogu v2 implementation
- [#81] small refactorings on configDiff implementation
- [#81] update various dependencies
- [#81] use go 1.23

### Fixed
- [#81] fix go-linter to support go 1.23
- [#81] fix superfluous response headers in tests

### Removed
- [#81] remove encryption for sensitive dogu config
- [#81] remove etcd from default list of required components

## [v1.1.0] - 2024-09-18
### Changed
- [#79] Relicense to AGPL-3.0-only

## [v1.0.0] - 2024-03-22
### Added
- [#59] Restart dogus if needed by config changes.

## [v0.2.1] - 2024-03-20
### Fixed
- [#74] Fix ldap-mapper dependency check by ignoring registrator Dogu.

## [v0.2.0] - 2024-03-20
### Added
- [#71] Add optional volume mount for self-signed certificate for the dogu registry.

## [v0.1.0] - 2024-03-20
### Added
General:
- [#1] Initially set up operator and Blueprint CRD
- [#4] Set up domain model
- [#4] Add static validation of blueprint specs
- [#4] Add dynamic validation of blueprints via dogu specification
- [#4] Process Blueprint CRs in cluster
- [#4] Calculate effective blueprint
- [#7] Create diff between effective blueprint and cluster state
- [#12] implement maintenance mode
- [#15] Check if required components are installed
- [#15] Check component health
- [#17] add health checks before and after applying the blueprint
- [#22] Add `dryRun` option. If `dryRun` is active the blueprint procedure stops before applying resources to the cluster and remains in the actual state. One can set the option to false and continue at this state.
- [#66] Write Event to Blueprint CR if parsing the Blueprint or Blueprint Mask fails

Dogu-specific:
- [#4] introduce flag `allowDoguNamespaceSwitch` for dogu namespace switch
- [#9] Check dogu health
- [#11] apply new dogu states based on blueprint
- [#20] Add exception for `nginx` dogu dependency validation. Map this dependency to `nginx-ingress` and `nginx-static`.
- [#30] Implement dogu namespace switch.
- [#55] Add dogu specific config for the volume size and the reverse proxy.

Component-specific:
- [#19] Create component differences between effective blueprint and cluster state
- [#14] Apply new component states based on blueprint.
- [#61] Add process for safe self upgrades of the blueprint operator
- [#62] Add component specific property map to configure attributes like deployNamespace or helm values in k8s.

Config-specific:
- [#42] Implement config repositories
- [#48] Save config diff in cluster
- [#39] Encrypt sensitive data
- [#36] Set registry configuration for dogu and global config.
- [#38] Censor all sensitive configuration data after applying the blueprint
- [#45] Set registry configuration for encrypted values.

### Fixed
- [68] Make dogu config or sensitiveConfig not required if one of them is specified.

