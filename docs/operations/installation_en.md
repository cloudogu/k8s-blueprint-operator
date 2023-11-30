# Installation

To install the `k8s-blueprint-operator` as a component,
apply the following component to the cluster namespace
the Cloudogu MultiNode EcoSystem is running in:

```yaml
apiVersion: k8s.cloudogu.com/v1
kind: Component
metadata:
  name: k8s-blueprint-operator
spec:
  name: k8s-blueprint-operator
  namespace: k8s
  # version: <desired-version> # You can specify a version here, otherwise, latest will be used.
```
