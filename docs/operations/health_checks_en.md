# Health checks

Health checks are carried out for the Dogus and the components before and after the blueprint is applied.

The checks after applying the blueprint wait until all health checks are OK.
The timeout and check interval can be set in the [Health-Config](#health-config).

## Dogus

The health of all installed Dogus is checked.

If this is not desired, the Dogu health checks can be deactivated in the blueprint,
by setting `spec.ignoreDoguHealth` to `true`.

## Components

First, it is checked that the required components are installed.
Then the health of all installed components is checked.

For the configuration of required components, see [Health-Config](#health-config).

If the component health checks are not to be executed, they can be deactivated in the blueprint,
by setting `spec.ignoreComponentHealth` to `true`.

## Health-Config

The health configuration can be overwritten in the `valuesYamlOverwrite` field of the component CR of the blueprint operator.
The following example shows the possible settings with their default configuration:

```yaml
valuesYamlOverwrite: |
  healthConfig:
    components:
      required: # These components are required for health checks to succeed.
      - name: k8s-etcd
      - name: k8s-dogu-operator
      - name: k8s-service-discovery
      - name: k8s-component-operator
    wait: # Define timeout and check-interval for the ecosystem to become healthy.
      timeout: 10m
      interval: 10s
```