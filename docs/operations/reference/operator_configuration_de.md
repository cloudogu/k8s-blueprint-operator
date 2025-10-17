# Operator-Konfiguration

Der `k8s-blueprint-operator` wird als Komponente bereitgestellt und kann während der Installation konfiguriert werden.
Die Standardwerte sind im Helm-Chart der Komponente definiert, können aber über das Feld `spec.valuesYamlOverwrite` in
der `Component`-Custom-Resource überschrieben werden.

## Überschreiben von Standardwerten

Um Ihre eigene Konfiguration bereitzustellen, fügen Sie das Feld `valuesYamlOverwrite` zu Ihrer `Component`-Definition
hinzu. Dieses Feld akzeptiert einen mehrzeiligen YAML-String, der die Struktur der internen `values.yaml`-Datei
widerspiegelt.

### Beispiel

Hier ist ein Beispiel, wie Sie das Memory-Ressourcenlimit auf 150M erhöhen und das Log-Level auf `debug` setzen:

```yaml
apiVersion: k8s.cloudogu.com/v1
kind: Component
metadata:
  name: k8s-blueprint-operator
  # Die Komponente sollte sich im selben Namespace wie die anderen EcoSystem-Komponenten befinden
  namespace: ecosystem
spec:
  name: k8s-blueprint-operator
  # Der Namespace, in den der Operator installiert werden soll
  namespace: k8s
  version: 2.8.0 # Verwenden Sie Ihre gewünschte Version
  valuesYamlOverwrite: |
    manager:
      resourceLimits:
        memory: 150M
      env:
        logLevel: debug
```

---

## Konfigurationsparameter

Die folgenden Abschnitte beschreiben die verfügbaren Konfigurationsparameter.

### `global`

Globale Einstellungen, die mehrere Teile der Bereitstellung beeinflussen können.

| Parameter          | Beschreibung                                                                                    | Standardwert                               |
|:-------------------|:------------------------------------------------------------------------------------------------|:-------------------------------------------|
| `imagePullSecrets` | Eine Liste von Kubernetes-Secrets, die zum Pullen von Container-Images verwendet werden sollen. | `[ { name: "ces-container-registries" } ]` |

### `manager`

Konfiguration für den Controller-Manager-Pod des Operators.

| Parameter                   | Beschreibung                                                                                                                                                                                  | Standardwert                      |
|:----------------------------|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|:----------------------------------|
| `replicas`                  | Die Anzahl der auszuführenden Operator-Pods.                                                                                                                                                  | `1`                               |
| `image.registry`            | Die Container-Registry für das Operator-Image.                                                                                                                                                | `docker.io`                       |
| `image.repository`          | Das Repository für das Operator-Image.                                                                                                                                                        | `cloudogu/k8s-blueprint-operator` |
| `image.tag`                 | Der Tag des bereitzustellenden Operator-Images.                                                                                                                                               | `2.8.0`                           |
| `imagePullPolicy`           | Die Kubernetes Image-Pull-Policy.                                                                                                                                                             | `IfNotPresent`                    |
| `env.logLevel`              | Die Ausführlichkeit der Protokollierung. Kann `info`, `debug` oder `trace` sein.                                                                                                              | `info`                            |
| `env.stage`                 | Die Bereitstellungsphase. Kann `production` oder `development` sein.                                                                                                                          | `production`                      |
| `resourceLimits.memory`     | Das Speicherlimit für den Operator-Container.                                                                                                                                                 | `105M`                            |
| `resourceRequests.cpu`      | Die CPU-Anforderung für den Operator-Container.                                                                                                                                               | `15m`                             |
| `resourceRequests.memory`   | Die Speicheranforderung für den Operator-Container.                                                                                                                                           | `105M`                            |
| `networkPolicies.enabled`   | Wenn `true`, werden `NetworkPolicy`-Ressourcen erstellt, um den Datenverkehr einzuschränken.                                                                                                  | `true`                            |
| `reconciler.debounceWindow` | Das Zeitfenster, in dem auf weitere Cluster-Ereignisse (z. B. ConfigMap-Änderungen) gewartet wird, bevor eine neue Reconciliation gestartet wird. Dies verhindert übermäßige Reconciliations. | `10s`                             |

### `doguRegistry`

Konfiguration bezüglich der Dogu-Registry, aus der Dogu-Spezifikationen abgerufen werden.

| Parameter            | Beschreibung                                                                                           | Standardwert         |
|:---------------------|:-------------------------------------------------------------------------------------------------------|:---------------------|
| `certificate.secret` | Der Name des Kubernetes-Secrets, das das TLS-Zertifikat für den Zugriff auf die Dogu-Registry enthält. | `dogu-registry-cert` |
