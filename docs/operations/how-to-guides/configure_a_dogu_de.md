# Ein Dogu konfigurieren

Die Konfiguration für Dogus und das gesamte Cloudogu EcoSystem kann deklarativ innerhalb der `Blueprint`-Ressource im Abschnitt `spec.blueprint.config` verwaltet werden. Dies ermöglicht es Ihnen, Ihre Konfiguration zusammen mit Ihren Anwendungsdefinitionen in der Versionskontrolle zu verwalten.

## Konfigurationsbereiche

Es gibt zwei Bereiche für die Konfiguration:

1.  **Globale Konfiguration (`config.global`):** Diese Schlüssel-Wert-Paare gelten für das gesamte EcoSystem und sind für alle Dogus zugänglich.
2.  **Dogu-spezifische Konfiguration (`config.dogus`):** Diese Schlüssel-Wert-Paare sind auf ein einzelnes Dogu ausgerichtet.

```yaml
# blueprint.yaml
spec:
  blueprint:
    config:
      global:
        # Globale Schlüssel-Wert-Paare
      dogus:
        # Dogu-spezifische Schlüssel-Wert-Paare
        <dogu-name>:
          # ...
```

---

## Klartext-Konfiguration setzen

Für normale Konfigurationen können Sie Werte direkt im Blueprint angeben. Der Operator speichert diese in einer Kubernetes `ConfigMap` und mounted sie in das entsprechende Dogu.

### Beispiel

```yaml
config:
  global:
    - key: "admin/mail"
      value: "admin@my-ces.com"
  dogus:
    redmine:
      - key: "logging/root"
        value: "DEBUG"
```

- Der globale Schlüssel `admin/mail` wird auf `admin@my-ces.com` gesetzt.
- Für das `redmine`-Dogu wird der Schlüssel `logging/root` auf `DEBUG` gesetzt.

---

## Sensible Konfiguration setzen

Bei sensiblen Daten wie Passwörtern oder API-Token dürfen Sie den Wert nicht direkt im Blueprint platzieren. Stattdessen beziehen Sie ihn aus einem vorhandenen Kubernetes `Secret`.

Dazu markieren Sie die Konfiguration als `sensitive: true` und geben eine `secretRef` an, die auf den Namen des `Secret` und den `key` innerhalb dieses Secrets verweist.

**Voraussetzung:** Das Kubernetes `Secret`, auf das Sie verweisen, muss existieren, bevor Sie das Blueprint anwenden.

### Beispiel

Stellen Sie zunächst sicher, dass ein Secret existiert:
```yaml
# mein-ldap-secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: ldap-credentials
stringData:
  password: "super-geheimes-passwort"
```

Referenzieren Sie es dann in Ihrem Blueprint:
```yaml
config:
  dogus:
    usermgmt:
      - key: "ldap/password"
        sensitive: true
        secretRef:
          name: "ldap-credentials" # Name des Secrets
          key: "password"         # Schlüssel innerhalb des Secrets
```

---

## Einen Konfigurationsschlüssel löschen

Um einen Konfigurationsschlüssel aus dem EcoSystem zu entfernen, markieren Sie ihn als `absent: true`.

### Einen Klartext-Schlüssel löschen

Bei normalen, nicht-sensiblen Schlüsseln müssen Sie nur den `key` und `absent: true` angeben.

```yaml
config:
  global:
    # Dies entfernt den 'logging/root'-Schlüssel aus der globalen Konfiguration
    - key: "logging/root"
      absent: true
```

### Einen sensiblen Schlüssel löschen

Um einen sensiblen Konfigurationsschlüssel zu entfernen, müssen Sie **sowohl** `absent: true` als auch `sensitive: true` angeben. Dies weist den Operator an, den Schlüssel im Secret-Speicher des Dogus anstelle seiner öffentlichen ConfigMap zu suchen.

```yaml
config:
  dogus:
    usermgmt:
      # Dies entfernt den sensiblen 'ldap/password'-Schlüssel
      - key: "ldap/password"
        absent: true
        sensitive: true
```

---

## Konfigurationsregeln und Validierungen

Der Blueprint-Operator erzwingt mehrere Regeln, um sicherzustellen, dass die Konfiguration gültig ist. Ein ungültiges Blueprint schlägt bei der Validierung fehl und wird nicht angewendet. Die folgenden Kombinationen sind **nicht erlaubt**:

- Ein Konfigurationseintrag kann **nicht** gleichzeitig einen `value` haben und `absent: true` sein.
- Ein Konfigurationseintrag kann **nicht** gleichzeitig eine `secretRef` haben und `absent: true` sein.
- Ein Konfigurationseintrag kann **nicht** gleichzeitig einen `value` und eine `secretRef` haben.
- Ein Konfigurationseintrag mit einer `secretRef` **muss** auch `sensitive: true` sein.
- Ein Konfigurationseintrag mit `sensitive: true` **muss** `secretRef` verwenden und darf keinen Klartext-`value` haben.

---

## Vollständiges Beispiel

Hier ist ein Blueprint, das alle Konfigurationstypen demonstriert:

```yaml
apiVersion: k8s.cloudogu.com/v2
kind: Blueprint
metadata:
  name: blueprint-mit-config
spec:
  displayName: "Konfigurationsbeispiel"
  blueprint:
    dogus:
      - name: "official/usermgmt"
        version: "2.8.1-1"
      - name: "official/redmine"
        version: "5.1.2-3"
    config:
      global:
        # Einen globalen Klartextwert setzen
        - key: "admin/mail"
          value: "admin@my-ces.com"
        # Einen globalen Schlüssel entfernen
        - key: "old/global/key"
          absent: true
      dogus:
        usermgmt:
          # Einen sensiblen Wert aus einem Secret setzen
          - key: "ldap/password"
            sensitive: true
            secretRef:
              name: "ldap-credentials"
              key: "password"
          # Einen sensiblen Schlüssel entfernen
          - key: "old/ldap/password"
            absent: true
            sensitive: true
        redmine:
          # Einen dogu-spezifischen Klartextwert setzen
          - key: "logging/root"
            value: "DEBUG"
```