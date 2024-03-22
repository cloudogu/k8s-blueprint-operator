# k8s-blueprint-operator Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

