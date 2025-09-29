# Applying Blueprints

You can apply a blueprint by applying a `Blueprint` resource to the cluster namespace where the Cloudogu MultiNode EcoSystem is running in:

```yaml
apiVersion: k8s.cloudogu.com/v2
kind: Blueprint
metadata:
  labels:
    app: ces
    app.kubernetes.io/name: k8s-blueprint-lib
  name: my-blueprint
spec:
  displayName: "Blueprint Sample v6.834"
  # put your blueprint here
  blueprint: 
    dogus: ...
    config:
      global: ...
      dogus: ...
  # put your blueprint-mask here
  blueprintMask:
      dogus: ...
```
The document [blueprint format](https://github.com/cloudogu/k8s-blueprint-lib/blob/develop/docs/operations/blueprintV2_format_en.md) describes the structure of the Blueprint in detail.
You may see examples of Blueprint-CRs in the [sample repository](https://github.com/cloudogu/k8s-ecosystem-samples/tree/main/blueprints). With `k8s-blueprint-operator` properly being installed, you can apply it to the cluster like this:

```bash
kubectl apply -n ecosystem -f k8s_v2_blueprint.yaml
```

**Note:** Only one blueprint is permitted per namespace. Either change the existing one or apply with the same name to update it.
