# Fehlerbehebung bei einem fehlerhaften Blueprint

Wenn sich ein Blueprint nicht wie erwartet verhält, ist die `Blueprint`-Custom-Resource selbst der beste Ausgangspunkt für Ihre Untersuchung. Sie enthält detaillierte Statusinformationen, einen Aktionsplan und ist mit Events verknüpft, die ein klares Bild der Operator-Aktivität liefern.

Dieser Leitfaden führt Sie durch die wichtigsten Bereiche, die bei der Fehlerbehebung zu überprüfen sind.

## 1. Überprüfen Sie die Statusbedingungen des Blueprints

Die unmittelbarste Informationsquelle ist das Feld `status.conditions` Ihrer Blueprint-Ressource. Sie können dies mit `kubectl describe` oder `kubectl get blueprint -o yaml` anzeigen.

```bash
kubectl describe blueprint <your-blueprint-name> -n <your-namespace>
```

Suchen Sie nach dem Abschnitt `Status`. Hier sind die wichtigsten Bedingungen und ihre Bedeutung:

- **`Valid`**: Wenn diese Bedingung `False` ist, bedeutet dies, dass ein struktureller oder logischer Fehler in Ihrer Blueprint-Definition vorliegt. Der Grund und die Meldung sagen Ihnen oft genau, was falsch ist (z. B. ein Syntaxfehler oder eine fehlende Abhängigkeit für ein Dogu).

- **`Executable`**: Dies ist `False`, wenn die berechneten Änderungen nicht zulässig sind. Der häufigste Grund ist ein versuchtes Dogu-Downgrade, das standardmäßig blockiert ist. Die mit dieser Bedingung verbundene Meldung erklärt die problematische Änderung.

- **`EcosystemHealthy`**: Dies zeigt an, ob der Operator darauf wartet, dass das Ecosystem healthy wird, bevor Änderungen angewendet werden. Wenn es `False` ist, bedeutet dies, dass ein oder mehrere Dogus nicht in einem bereiten Zustand sind.

- **`Completed`**: Dies zeigt an, ob das Blueprint vollständig angewendet wurde. Wenn es lange nach dem Anwenden `False` ist, bedeutet dies, dass der Operator noch arbeitet oder feststeckt.

- **`LastApplySucceeded`**: Dies ist eine kritische Bedingung für die Fehlerbehebung. Wenn ein Vorgang fehlschlägt (z. B. das Anwenden einer ConfigMap oder die Installation eines Dogus), wird diese Bedingung `False`. **Entscheidend ist, dass sie die letzte Fehlermeldung enthält** und über mehrere Reconciliation-Loops hinweg bestehen bleibt, bis das Blueprint erfolgreich abgeschlossen ist. Dies ermöglicht es Ihnen, die Grundursache eines Fehlers zu sehen, selbst wenn der Operator es erneut versucht.

Beginnen Sie damit, nach einer Bedingung zu suchen, die `False` ist, und lesen Sie die zugehörige `message` für Details.

## 2. Analysieren Sie den StateDiff

Wenn das Blueprint gültig und das Ecosystem healthy ist, aber Änderungen nicht wie erwartet angewendet werden, ist das Feld `status.stateDiff` Ihre nächste Anlaufstelle. Dieses Feld zeigt den genauen Plan, den der Operator durch den Vergleich des gewünschten Zustands (Ihres Blueprints) mit dem tatsächlichen Zustand des Clusters berechnet hat.

Sie können es mit `kubectl get` anzeigen:

```bash
kubectl get blueprint <your-blueprint-name> -n <your-namespace> -o yaml
```

Der `stateDiff` zeigt Ihnen genau, welche Dogus und Konfigurationen der Operator hinzufügen, aktualisieren, entfernen oder ändern möchte. Dies ist nützlich, um Folgendes zu erkennen:
- **Unerwartete Änderungen**: Enthält der Diff Änderungen, die Sie nicht erwartet haben? Dies könnte auf ein Problem mit Ihrem `blueprint` oder `blueprintMask` hinweisen.
- **Problematische Operationen**: Der Diff könnte explizit ein geplantes Downgrade zeigen, was erklären würde, warum die Bedingung `Executable` `False` ist.

## 3. Events überprüfen

Kubernetes-Events bieten ein chronologisches Protokoll dessen, was der Operator getan hat. Wenn Sie `kubectl describe blueprint <your-blueprint-name>` ausführen, erhalten Sie auch eine Liste der zugehörigen Events am Ende.

Diese Events zeigen Ihnen unter anderem:
- Wann ein Reconciliation-Loop gestartet wurde.
- Das Ergebnis von Validierungsprüfungen.
- Den Beginn und das Ende der Anwendungsphase.
- Alle Fehler, die bei der Interaktion mit anderen Ressourcen aufgetreten sind.

## 4. Überprüfen Sie die Operator-Logs

Für die detailliertesten Informationen müssen Sie die Logs des `k8s-blueprint-operator`-Pods selbst überprüfen. Hier finden Sie detaillierte Fehlermeldungen und Stack-Traces, die die genaue Codezeile lokalisieren können, in der ein Fehler aufgetreten ist.

1.  **Finden Sie den Operator-Pod:**
    ```bash
    kubectl get pods -n <operator-namespace> -l app.kubernetes.io/name=k8s-blueprint-operator
    ```
    *(Der Namespace ist typischerweise `ecosystem` oder dort, wo Sie die Komponente installiert haben).

2.  **Streamen Sie die Logs:**
    ```bash
    kubectl logs -f <operator-pod-name> -n <operator-namespace>
    ```

Wenn die Standard-Logs auf `info`-Ebene nicht ausreichen, können Sie das Loglevel auf `debug` oder `trace` erhöhen. Dies geschieht durch Ändern der `Component`-Ressource für den Operator. Detaillierte Anweisungen finden Sie in der [Operator-Konfigurationsdokumentation](../reference/operator_configuration_de.md).