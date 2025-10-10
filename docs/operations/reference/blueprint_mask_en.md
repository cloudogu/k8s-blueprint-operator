# Blueprint Mask

The `blueprintMask` provides a powerful way to customize a blueprint for a specific environment without altering the original blueprint definition. It acts as a filter, allowing you to selectively disable certain dogus that are defined in the main `blueprint` section.

## Use Case: Standardized Blueprints, Customized Deployments

Imagine maintaining a single, comprehensive, and battle-tested blueprint that defines a full installation of your application suite. This "golden-master" blueprint is used across multiple teams or customers.

However, not every team needs every application. For example, one team might not need the `redmine` dogu.

Instead of creating and maintaining a separate, slightly different blueprint file for that team, you can use the `blueprintMask`. You apply the same complete blueprint to every cluster but use a specific mask for that team's cluster to prevent `redmine` from being installed.

This approach has several advantages:
- **Single Source of Truth**: You manage one master blueprint, reducing complexity and the risk of configuration drift.
- **Consistency**: All environments are based on the same tested foundation.
- **Flexibility**: You can easily enable or disable dogus for any given installation on the fly.

## How it Works

The `k8s-blueprint-operator` first reads the `blueprint` section and then applies the `blueprintMask` over it. If a dogu is listed in the mask with `absent: true`, it is removed from the final set of dogus to be installed or managed. The result is called the "effective blueprint."

## Example

Consider the following `Blueprint` resource. The `blueprint` section defines both `scm` and `redmine`.

```yaml
apiVersion: k8s.cloudogu.com/v2
kind: Blueprint
metadata:
  name: my-instance-blueprint
spec:
  # This is the "master" blueprint with all possible dogus
  blueprint:
    dogus:
      - name: "official/scm"
        version: "3.11.0-1"
      - name: "official/redmine"
        version: "6.0.6-2"

  # This mask customizes the blueprint for this specific instance
  blueprintMask:
    dogus:
      - name: "official/redmine"
        absent: true
```

### Outcome

When the operator processes this resource:
1. It sees `scm` and `redmine` in the `blueprint`.
2. It then applies the `blueprintMask`, which says `redmine` should be absent.
3. The resulting **effective blueprint** only contains the `scm` dogu.

As a result, the operator will only install `official/scm:3.11.0-1` and will completely ignore the `redmine` dogu. If `redmine` was already installed, it would be marked for uninstallation.

## Syntax

The structure of the `blueprintMask` mirrors the `blueprint` itself. To exclude a dogu, you only need to provide its `name` and the `absent: true` flag.

```yaml
blueprintMask:
  dogus:
    - name: "<namespace>/<dogu-name>"
      absent: true
    # Add other dogus to exclude here
```
