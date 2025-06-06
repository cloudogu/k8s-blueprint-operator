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
      "dogus": [ ... ],
      "components": [ ... ],
      "config": {
        "global": { ... },
        "dogus": { ... }
      }
    }
  # put your blueprint-mask.json here
  blueprintMask: |
    {
      "blueprintMaskApi": "v1",
      "blueprintMaskId": "my-blueprint-mask",
      "dogus": [ ... ]
    }
```
The document [blueprint format](https://github.com/cloudogu/k8s-blueprint-lib/blob/develop/docs/operations/blueprintV2_format_en.md) describes the structure of the Blueprint in detail.
You may see examples of Blueprint-CRs in the [sample repository](https://github.com/cloudogu/k8s-ecosystem-samples/tree/main/blueprints). With `k8s-blueprint-operator` properly being installed, you can apply it to the cluster like this:

```bash
kubectl apply -n ecosystem -f k8s_v1_blueprint.yaml
```
