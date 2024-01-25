# k8s-blueprint-operator Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- [#1] Initially set up operator and Blueprint CRD
- [#4] Set up domain model
- [#4] Add static validation of blueprint specs
- [#4] Add dynamic validation of blueprints via dogu specification
- [#4] Process Blueprint CRs in cluster
- [#4] Calculate effective blueprint
- [#4] introduce flag `allowDoguNamespaceSwitch` for dogu namespace switch
- [#7] Create diff between effective blueprint and cluster state
- [#9] Check dogu health
- [#11] apply new dogu states based on blueprint
- [#12] implement maintenance mode
- [#15] Check if required components are installed
- [#15] Check component health
- [#17] add health checks before and after applying the blueprint
- [#20] Add exception for `nginx` dependency validation. Map this dependency to `nginx-ingress` and `nginx-static`.
