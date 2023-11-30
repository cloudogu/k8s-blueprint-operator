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