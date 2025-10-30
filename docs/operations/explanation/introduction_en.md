# Introduction

The `k8s-blueprint-operator` is a Kubernetes operator designed to manage complex application landscapes within the Cloudogu EcoSystem.

## The Problem

Managing a suite of interconnected applications (called `dogus`) and their configurations can be complex. Ensuring that the correct versions of each application are deployed and that their configurations are in sync is critical for a stable system. This process is especially challenging during initial setup and subsequent upgrades.

## The Solution: Blueprints

This operator introduces a Custom Resource called a `Blueprint`. A Blueprint is a single, declarative YAML file where you define the entire desired state of your application landscape:

*   The specific versions of all `dogus` that should be installed.
*   The configuration for each `dogu`.
*   Global configuration that applies to the entire system.

By applying a single `Blueprint` resource to your Kubernetes cluster, you trigger the `k8s-blueprint-operator`, which then takes on the work of making the cluster's state match the Blueprint's definition. It acts as a controller, continuously working to install, upgrade, and configure your applications until the desired state is met.

This approach allows you to package and version-control your entire application setup, treating your ecosystem's state as code. It is particularly useful for managing tested software packages, ensuring that you can deploy a known, tested combination of applications and configurations reliably.
