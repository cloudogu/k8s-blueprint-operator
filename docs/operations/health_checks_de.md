# Health-Checks

Vor und nach dem Anwenden des Blueprints wird gewartet, dass das Ecosystem healthy ist.
Dabei wird folgendes geprüft:
- Health aller Dogus anhand der Dogu-CRs
- Health aller Components anhand der Component-CRs
- Überprüfung, ob alle notwendigen Components installiert sind, die für das Blueprint gebraucht werden

Die Health-Checks verwenden einen eingebauten Retry. 
Timeout und Check-Interval lassen sich dafür in der [Health-Config](#health-config) festlegen.

## Health ignorieren

Die Health-Checks vor der Ausführung des Blueprints können deaktiviert werden:
- für Dogus, wenn `spec.ignoreDoguHealth` auf `true` gesetzt wird,
- für Components, wenn `spec.ignoreComponentHealth` auf `true` gesetzt wird.

So ist es möglich, per Blueprint Fehler an Dogus und Komponenten zu beheben.
Für ein Dogu-Upgrade muss ein Dogu allerdings healthy sein, um Pre-Upgrade-Skripte ausführen zu können.
Das Ignorieren der Dogu-Health kann also zu Folgefehlern während der Ausführung des Blueprints führen.

## Health-Config

Die Health-Konfiguration kann im Feld `valuesYamlOverwrite` der Komponenten-CR des Blueprint-Operators überschrieben werden.
Folgendes Beispiel zeigt die möglichen Einstellungen mit ihrer Default-Konfiguration:

```yaml
valuesYamlOverwrite: |
  healthConfig:
    components:
      required: # These components are required for health checks to succeed.
      - name: k8s-dogu-operator
      - name: k8s-service-discovery
      - name: k8s-component-operator
    wait: # Define timeout and check-interval for the ecosystem to become healthy.
      timeout: 10m
      interval: 10s
```