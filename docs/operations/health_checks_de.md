# Health-Checks

Vor und nach dem Anwenden des Blueprints werden Health-Checks für die Dogus als auch die Komponenten ausgeführt.

Bei den Checks nach dem Anwenden des Blueprints wird gewartet, bis alle Health-Checks in Ordnung sind.
Timeout und Check-Interval lassen sich dafür in der [Health-Config](#health-config) festlegen.

## Dogus

Es wird die Health aller installierten Dogus geprüft.

Wird dies nicht gewünscht, können die Dogu-Health-Checks in der Blueprint deaktiviert werden,
indem `spec.ignoreDoguHealth` auf `true` gesetzt wird.

## Komponenten

Zunächst wird überprüft, dass benötigte Komponenten installiert sind.
Dann wird die Health aller installierten Komponenten geprüft.

Für die Konfiguration benötigter Komponenten, siehe [Health-Config](#health-config).

Sollen die Komponenten-Health-Checks nicht ausgeführt werden, können sie in der Blueprint deaktiviert werden,
indem `spec.ignoreComponentHealth` auf `true` gesetzt wird.

## Health-Config

Die Health-Konfiguration kann im Feld `valuesYamlOverwrite` der Komponenten-CR des Blueprint-Operators überschrieben werden.
Folgendes Beispiel zeigt die möglichen Einstellungen mit ihrer Default-Konfiguration:

```yaml
valuesYamlOverwrite: |
  healthConfig:
    components:
      required: # These components are required for health checks to succeed.
      - name: k8s-etcd
      - name: k8s-dogu-operator
      - name: k8s-service-discovery
      - name: k8s-component-operator
    wait: # Define timeout and check-interval for the ecosystem to become healthy.
      timeout: 10m
      interval: 10s
```