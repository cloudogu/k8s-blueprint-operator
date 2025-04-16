# Health checks

Before and after applying the blueprint, the ecosystem is checked to ensure that it is healthy.
The following is checked:
- Health of all Dogus based on the Dogu-CRs
- Health of all components based on the component CRs
- Check whether all necessary components required for the blueprint are installed

The health checks use a built-in retry.
The timeout and check interval can be defined in the [Health-Config](#health-config).

## Ignoring health

Upfront health checks can be deactivated:
- for Dogus, if `spec.ignoreDoguHealth` is set to `true`,
- for components, if `spec.ignoreComponentHealth` is set to `true`.

This makes it possible to fix errors on Dogus and components via Blueprint.
For a Dogu upgrade, however, a Dogu must be healthy in order to be able to execute pre-upgrade scripts.
Ignoring the dogu health can therefore lead to subsequent errors during the execution of the blueprint.

## Health-Config

The health configuration can be overwritten in the `valuesYamlOverwrite` field of the component CR of the blueprint operator.
The following example shows the possible settings with their default configuration:

```yaml
valuesYamlOverwrite: |
  healthConfig:
    components:
      required: # These components are required for health checks to succeed.
      - name: k8s-dogu-operator
      - name: k8s-service-discovery
      - name: k8s-component-operator
    wait: # Define timeout and check-interval for the ecosystem to become healthy.
      timeout: 10m
      interval: 10s
```