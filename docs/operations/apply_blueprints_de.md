# Anwenden von Blueprints

Sie können einen Blueprint anwenden, indem Sie eine `Blueprint`-Ressource auf den Cluster-Namespace anwenden, in dem das Cloudogu MultiNode EcoSystem läuft:

```yaml
apiVersion: k8s.cloudogu.com/v1
kind: Blueprint
metadata:
  name: my-blueprint
spec:
  # fügen Sie die blueprint.json hier ein
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
  # fügen Sie hier die blueprint-mask.json ein
  blueprintMask: |
    {
      "blueprintMaskApi": "v1",
      "blueprintMaskId": "my-blueprint-mask",
      "dogus": [ ... ]
    }
```