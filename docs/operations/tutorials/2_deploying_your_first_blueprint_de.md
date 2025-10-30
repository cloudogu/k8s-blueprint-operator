# Dein erstes Blueprint deployen

Dieses Tutorial leitet dich an, deine erste Anwendung (ein `dogu`) mit einem einfachen Blueprint zu deployen. Ein Blueprint ist eine Custom Resource, die den gewünschten Zustand von Anwendungen und deren Konfigurationen in deinem Cloudogu EcoSystem deklarativ definiert.

## Ziel

Am Ende dieses Tutorials wirst du das `mysql`-Dogu mit einem Blueprint deployed haben.

## 1. Die Blueprint-Datei erstellen

Erstelle zuerst eine neue YAML-Datei mit dem Namen `my-first-blueprint.yaml` und füge den folgenden Inhalt hinzu:

```yaml
apiVersion: k8s.cloudogu.com/v2
kind: Blueprint
metadata:
  name: my-first-blueprint
spec:
  displayName: "Mein erstes Blueprint"
  blueprint:
    dogus:
      - name: "official/mysql"
        version: "8.4.5-4"
```

### Was macht diese Datei?

*   `kind: Blueprint`: Dies teilt Kubernetes mit, dass die Ressource ein Blueprint ist, mit dem der `k8s-blueprint-operator` umgehen kann.
*   `metadata.name`: Dies gibt unserem Blueprint einen eindeutigen Namen, `my-first-blueprint`.
*   `spec.blueprint.dogus`: Dies ist die Liste der Anwendungen, die im System vorhanden sein sollen. Hier definieren wir ein Dogu:
    *   `name: "official/mysql"`: Gibt das `mysql`-Dogu aus dem `official`-Namespace an.
    *   `version: "8.4.5-4"`: Gibt die genaue Version an, die wir installieren möchten.

## 2. Das Blueprint anwenden

Verwende nun `kubectl`, um diese Ressource auf deinen Cluster anzuwenden. Dies veranlasst den `k8s-blueprint-operator`, mit dem Abgleich des Zustands zu beginnen.

```bash
kubectl apply -f my-first-blueprint.yaml -n ecosystem
```

## 3. Den Status überprüfen

Der Operator wird nun daran arbeiten, den im Blueprint definierten Zustand zu erreichen. Du kannst den Fortschritt überprüfen, indem du die Beschreibung der Blueprint-Ressource aufrufst:

```bash
kubectl describe blueprint my-first-blueprint
```

Schaue dir die Abschnitte `Status` und `Events` in der Ausgabe an. Du solltest Events sehen, die darauf hinweisen, dass der Operator das Blueprint validiert und die Änderungen anwendet. Sobald dies abgeschlossen ist, zeigt der Status an, dass das Blueprint `Completed` ist.

## Fazit

Glückwunsch! Du hast erfolgreich den `k8s-blueprint-operator` verwendet, um dein erstes Dogu zu deployen. Darauf aufbauend kannst du nun weitere Dogus hinzufügen oder Konfigurationen bereitstellen, wie in den anderen Anleitungen gezeigt.
