# Installation

Um den `k8s-blueprint-operator` als Komponente zu installieren,
wenden Sie die folgende Komponente auf den Cluster-Namespace an,
in dem das Cloudogu MultiNode EcoSystem läuft:

```yaml
apiVersion: k8s.cloudogu.com/v1
kind: Component
metadata:
  name: k8s-blueprint-operator
spec:
  name: k8s-blueprint-operator
  namespace: k8s
  # version: <gewünschte-Version> # Sie können hier eine Version angeben, ansonsten wird die neueste verwendet.
```