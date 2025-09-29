# Anwenden von Blueprints

Sie können einen Blueprint anwenden, indem Sie eine `Blueprint`-Ressource auf den Cluster-Namespace anwenden, in dem das Cloudogu MultiNode EcoSystem läuft:

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
  # fügen Sie die blueprint hier ein
  blueprint:
    dogus: ...
    config:
      global: ...
      dogus: ...
  # fügen Sie hier die blueprint-mask ein
  blueprintMask:
    dogus: ...
```

Das Dokument [Blueprint-Format](https://github.com/cloudogu/k8s-blueprint-lib/blob/develop/docs/operations/blueprintV2_format_de.md) beschreibt die Struktur des Blueprint im Detail.
Blueprint-CR-Beispiele können dem [Sample-Repository](https://github.com/cloudogu/k8s-ecosystem-samples/tree/main/blueprints) entnommen werden. Wenn `k8s-blueprint-operator` korrekt installiert wurde, lässt sich dies z. B. so auf den Cluster anwenden:

```bash
kubectl apply -n ecosystem -f k8s_v2_blueprint.yaml
```

**Hinweis:** Pro Namespace ist nur ein Blueprint zulässig. Ändern Sie entweder das vorhandene Blueprint oder wenden Sie erneut ein `kubectl apply` mit demselben Blueprint-Namen an, um es zu aktualisieren.
