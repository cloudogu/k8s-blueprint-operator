# Ein Blueprint anwenden oder aktualisieren

Sie können ein Blueprint anwenden, indem Sie eine `Blueprint`-Ressource im Cluster erstellen. Der `k8s-blueprint-operator` wendet die Änderungen dann automatisch auf Ihr Cloudogu EcoSystem an.

## Schritte

1.  **Erstellen Sie eine YAML-Datei**, die Ihre Blueprint-Definition enthält. Ein Beispiel finden Sie unten.
2.  **Wenden Sie sie mit `kubectl` auf den Cluster an**.

    ```bash
    kubectl apply -f ihre-blueprint-datei.yaml
    ```

**Hinweis:** Pro Namespace ist nur ein Blueprint erlaubt. Wenn Sie ein Blueprint mit demselben Namen wie ein bereits vorhandenes anwenden (`apply`), aktualisiert der Operator das Deployment entsprechend der neuen Definition.

## Vollständiges Blueprint-Beispiel

Hier ist ein Beispiel für eine `Blueprint`-Ressource, das mehrere Funktionen demonstriert, einschließlich der Definition von Dogus, der Konfigurationseinstellung und der Verwendung einer Blueprint-Maske.

```yaml
apiVersion: k8s.cloudogu.com/v2
kind: Blueprint
metadata:
  labels:
    app: ces
    app.kubernetes.io/name: k8s-blueprint-lib
  name: blueprint-beispiel
spec:
  displayName: "Blueprint Beispiel v6.834"
  blueprint:
    config:
      dogus:
        mysql:
          - key: "logging/root"
            value: "WARN"
          - key: "sa-ldap/password"
            sensitive: true
            secretRef:
              name: "ldap-sa-secret"
              key: "password"
          - key: "to/be/deleted"
            absent: true
      global:
        - key: "my/global/key"
          value: "myValue"
        - absent: true
          key: "global/to/be/delete"
    dogus:
      - name: "official/mysql"
        version: "8.4.6-1"
  blueprintMask:
    dogus:
      - absent: true
        name: "official/mysql"
```

Dieses Beispiel bewirkt Folgendes:

*   Definiert das `official/mysql`-Dogu in der Version `8.4.6-1`.
*   Setzt und entfernt mehrere Konfigurationswerte für das `mysql`-Dogu und global.
*   Verwendet eine `secretRef`, um einen sensiblen Wert aus einem Kubernetes-Secret zu beziehen.
*   Verwendet eine `blueprintMask`, um das `mysql`-Dogu als `absent` zu markieren, was dessen Installation durch dieses Blueprint effektiv verhindert.

Eine detaillierte Aufschlüsselung aller möglichen Felder finden Sie in der offiziellen [Blueprint-Format-Dokumentation](https://github.com/cloudogu/k8s-blueprint-lib/blob/develop/docs/operations/blueprintV2_format_de.md).