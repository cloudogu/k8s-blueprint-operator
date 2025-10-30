# Upgrade eines Dogus mit einem Blueprint

Im vorherigen Tutorial hast du deine erste Anwendung (`mysql`) mit einem Blueprint deployed. Diese Anleitung zeigt dir, wie du diese Anwendung auf eine neue Version upgradest, indem du einfach deine Blueprint-Definition änderst.

Dieser deklarative Ansatz ist eine der Kernstärken des Operators. Du definierst den gewünschten Zustand, und der Operator kümmert sich um die notwendigen Schritte, um dorthin zu gelangen.

### Voraussetzungen

Dieses Tutorial setzt voraus, dass du das Tutorial "Dein erstes Blueprint deployen" abgeschlossen hast und die Datei `my-first-blueprint.yaml` verfügbar ist.

## Schritt 1: Die Blueprint-Version ändern

Um eine Anwendung upzugraden, musst du nur die neue Version in deiner Blueprint-Datei angeben.

Öffne deine `my-first-blueprint.yaml`-Datei und ändere die Version des `mysql`-Dogus von `8.4.5-4` auf eine neuere Version, zum Beispiel `8.4.6-1`.

```diff
# my-first-blueprint.yaml

apiVersion: k8s.cloudogu.com/v2
kind: Blueprint
metadata:
  name: my-first-blueprint
spec:
  displayName: "Mein erstes Blueprint"
  blueprint:
    dogus:
      - name: "official/mysql"
-       version: "8.4.5-4"
+       version: "8.4.6-1"
```

## Schritt 2: Das aktualisierte Blueprint anwenden

Wende nun die geänderte Datei mit demselben Befehl wie zuvor auf deinen Cluster an. Der Operator erkennt die Änderung an der `Blueprint`-Ressource und startet eine neue Reconciliation Loop.

```bash
kubectl apply -f my-first-blueprint.yaml -n ecosystem
```

Der Operator vergleicht den neuen gewünschten Zustand aus deiner Datei mit dem aktuellen Zustand des Clusters und stellt fest, dass die Version von `mysql` unterschiedlich ist. Er wird automatisch feststellen, dass ein Upgrade erforderlich ist.

## Schritt 3: Das Upgrade beobachten

Wenn du schnell genug bist, kannst du sehen, was der Operator geplant hat, indem du die `StateDiff` im Status des Blueprints inspizierst. Dies bietet eine transparente Sicht auf die Aktionen, die der Operator ausführen wird.

```bash
kubectl get blueprint my-first-blueprint -n ecosystem -o yaml
```

Im Abschnitt `status.stateDiff.doguDiffs` siehst du möglicherweise einen Eintrag für `mysql` mit `upgrade` unter `neededActions` (nur bis der Operator das Upgrade abgeschlossen hat).

Um das Upgrade in Echtzeit zu beobachten, kannst du die `dogu`-Ressource selbst beobachten:

```bash
kubectl get dogu mysql --watch -n ecosystem
```

Du wirst sehen, wie sich der Status des Dogus ändert, während der Operator das Upgrade auf die neue Version durchführt.

## Alternative: `kubectl edit` verwenden

Für schnelle, interaktive Änderungen kannst du das Blueprint auch direkt im Cluster bearbeiten, ohne eine lokale Datei zu ändern:

```bash
kubectl edit blueprint my-first-blueprint -n ecosystem
```

Dadurch wird die Blueprint-Ressource in deinem Standardeditor geöffnet. Du kannst die Version ändern, speichern und die Datei schließen. Der Operator erkennt die Änderung und startet das Upgrade sofort.

**Hinweis:** Obwohl `kubectl edit` für Entwicklung oder Tests praktisch ist, empfehlen wir dringend, deine Blueprint-YAML-Dateien in einem Versionskontrollsystem (wie Git) für Produktionsumgebungen aufzubewahren. Dieser `GitOps`-Ansatz ermöglicht es dir, Änderungen an deiner Anwendungslandschaft sicher zu verfolgen, zu überprüfen und zurückzusetzen.

## Fazit

Glückwunsch! Du hast nun gesehen, wie der `k8s-blueprint-operator` das Application-Lifecycle-Management vereinfacht. Indem du deine Anwendungsversionen deklarativ in einer einzigen Datei verwaltest, kannst du Upgrades zuverlässig und vorhersagbar durchführen.
