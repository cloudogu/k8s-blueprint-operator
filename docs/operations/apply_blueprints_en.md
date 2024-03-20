# Applying Blueprints

You can apply a blueprint by applying a `Blueprint` resource to the cluster namespace where the Cloudogu MultiNode EcoSystem is running in:

```yaml
apiVersion: k8s.cloudogu.com/v1
kind: Blueprint
metadata:
  name: my-blueprint
spec:
  # put your blueprint.json here
  blueprint: |
    {
      "blueprintApi": "v2",
      "blueprintId": "my-blueprint",
      "dogus": [ ... ],
      "components": [ ... ],
      "registryConfig": { ... },
      "registryConfigEncrypted": { ... },
      "registryConfigAbsent": [ ... ]
    }
  # put your blueprint-mask.json here
  blueprintMask: |
    {
      "blueprintMaskApi": "v1",
      "blueprintMaskId": "my-blueprint-mask",
      "dogus": [ ... ]
    }
```

You may see [examples](../../samples/k8s_v1_blueprint.yaml) of Blueprint-CRs in the [Sample directory](../../samples/). With `k8s-blueprint-operator` properly being installed, you can apply it to the cluster like this:

```bash
kubectl -n ecosystem -f k8s_v1_blueprint.yaml
```
