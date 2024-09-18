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
      "dogus": [ ... ],
      "components": [ ... ],
      "config": {
        "global": { ... },
        "dogus": { ... }
      }
    }
  # fügen Sie hier die blueprint-mask.json ein
  blueprintMask: |
    {
      "blueprintMaskApi": "v1",
      "blueprintMaskId": "my-blueprint-mask",
      "dogus": [ ... ]
    }
```

Das Dokument [Blueprint-Format](./blueprint_format_de.md) beschreibt die Struktur des Blueprint im Detail.
Blueprint-CR-[Beispiele](../../samples/k8s_v1_blueprint.yaml) können dem [Sample-Verzeichnis](../../samples/) entnommen werden. Wenn `k8s-blueprint-operator` korrekt installiert wurde, lässt sich dies z. B. so auf den Cluster anwenden:

```bash
kubectl -n ecosystem -f k8s_v1_blueprint.yaml
```

