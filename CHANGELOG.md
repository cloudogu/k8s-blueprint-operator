# k8s-blueprint-operator Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
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

Dogu-specific:
- [#4] introduce flag `allowDoguNamespaceSwitch` for dogu namespace switch
- [#9] Check dogu health
- [#11] apply new dogu states based on blueprint
- [#20] Add exception for `nginx` dogu dependency validation. Map this dependency to `nginx-ingress` and `nginx-static`.
- [#30] Implement dogu namespace switch.

Component-specific:
- [#19] Create component differences between effective blueprint and cluster state
- [#14] Apply new component states based on blueprint.

Config-specific:
- [#42] Implement config repositories
- [#48] Save config diff in cluster
- [#36] Set registry configuration for dogu and global config.
- [#45] Set registry configuration for encrypted values.