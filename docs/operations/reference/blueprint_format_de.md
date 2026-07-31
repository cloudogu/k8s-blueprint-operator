# Blueprint-Format

Das Blueprint bietet die Möglichkeit Dogus und Konfigurationen hinzuzufügen oder zu entfernen.
Für Dogus können ebenfalls spezielle Konfigurationen definiert werden.
Diese werden an die entsprechenden Dogu-CRs geschrieben und werden nicht in der EcoSystem-Registry gespeichert.

Folgend werden alle Felder des Blueprint beschrieben und mit Beispielen veranschaulicht.

## BlueprintApi

* Pflichtfeld
* Datentyp: string
* Inhalt: Das Feld `blueprintApi` gibt die API-Version des Blueprints an.
* Beispiel: `blueprintApi: "v3"`

## Dogus

* Pflichtfeld
* Datentyp: Array<dogu>
* Inhalt: Das Feld `dogus` ist eine Liste von Dogus und beschreibt den Zustand der Dogus im System.
* Beispiel:
```
dogus:
  - name: official/mysql
    version: 8.4.6-1
```

### Dogu

Ein Dogu kann folgende Felder beinhalten:

#### Name

* Pflichtfeld
* Datentyp: string
* Inhalt: Gibt den Namen inklusive Namespace des Dogus an.
* Beispiel: `name: "official/cas"`

#### Absent

* Datentyp: boolean
* Inhalt: Gibt an, ob ein Dogu entfernt werden soll.
* Beispiel: `absent: true`

#### Version

* Bei `absent=true` optional. Sonst is die Version ein Pflichtfeld.
* Datentyp: string
* Inhalt: Gibt die Version des Dogus an.
* Beispiel: `version: "12.15-2"`


#### PlatformConfig

Das Feld `platformConfig` bietet die Möglichkeit für die Ausführungsplattform (z.B. Kubernetes) spezifische Konfigurationen zu übergeben.
Mit dieser Konfiguration können Ressourcen und Reverse-Proxy-Konfigurationen definiert werden.

##### Resource.minVolumeSize

* Optional
* Datentyp: string
* Inhalt: Gibt die minimale Volume-Size eines Dogus an. Falls das aktuelle Volume kleiner ist, wird eine Volume-Vergrößerung durchgeführt. Die Einheit muss mit dem Binär-Prefix angegeben werden (z.B. `Mi` oder `Gi`).
* Beispiel:
```
dogus: 
  name: "official/nexus"
  version: "3.59.0-2"
  platformConfig: 
    resource: 
      minVolumeSize: "5Gi"
```

> Der Dogu-Operator erstellt Dogus mit 2Gi Volumes. Das Nexus-Dogu benötigt ein größeres Volume und muss
> über diesen Eintrag konfiguriert werden.

> Das Verkleinern von Volumes wird nicht unterstützt.

##### ReverseProxy.maxBodySize

* Optional
* Datentyp: string
* Inhalt: Gibt die maximale Dateigröße für HTTP-Request an (`0`=unbegrenzt). Die Einheit muss mit dem Dezimal-Prefix angegeben werden (z.B. `M` oder `G`).
* Beispiel:
```
dogus:
  name: "official/nexus"
  version: "3.59.0-2"
  platformConfig: 
    reverseProxy: 
      maxBodySize: "1G"
```

##### ReverseProxy.rewriteTarget

* Optional
* Datentyp: string
* Inhalt: Definiert ein Rewrite-Target.
* Beispiel:
```
dogus:
  name: "official/postgresql"
  version: "12.15-2"
  platformConfig: 
    reverseProxy: 
      rewriteTarget: "/"
```

##### ReverseProxy.additionalConfig

* Optional
* Datentyp: string
* Inhalt: Fügt beliebige zusätzliche Proxy-Konfiguration hinzu.
* Beispiel:
```
dogus:
  name: "official/postgresql"
  version: "12.15-2"
  platformConfig: 
    reverseProxy: 
      additionalConfig: "<config>"
```

##### AdditionalMounts

* Optional
* Datentyp: Array
* Inhalt: Das Feld `additionalDataMounts` ist eine List die Elemente von Typ `AdditionalMount` enthält
* Beispiel:
```
dogus:
  name: "official/nginx-static"
  version: "1.26.3-2"
  platformConfig: 
    additionalMounts: 
      - sourceType: "ConfigMap"
        name: "my-html-configmap"
        volume: "customhtml"
        subfolder: "about"
      - sourceType: "Secret" 
        name: "my-ssh-key-secret"
        volume: "private_data"
```

###### AdditionalMount

* Datentyp: AdditionalMount
* Inhalt: Ein `AdditionalMount` definiert eine Quelle und ein Zielvolume, um Files in ein Dogu zu mounten.

###### AdditionalMount.SourceType

* Pflichtfeld
* Datentyp: Enum <ConfigMap; Secret>
* Inhalt: SourceType legt den Typ der Quelle fest.
  Gültige Optionen sind:
    - ConfigMap - Daten, die in einer kubernetes ConfigMap gespeichert sind.
    - Secret - Daten, die in einem kubernetes Secret gespeichert sind.
* Beispiel: `sourceType: ConfigMap`

###### AdditionalMount.Name

* Pflichtfeld
* Datentyp: String
* Inhalt: Name ist der Name der Datenquelle.
* Beispiel: `name: my-configmap`

###### AdditionalMount.Volume

* Pflichtfeld
* Datentyp: String
* Inhalt: Volume ist der Name des Volumes, in das die Daten gemountet werden sollen. Dieses wird in der jeweiligen
  dogu.json definiert.
* Beispiel: `volume: importHistory`

###### AdditionalMount.Subfolder

* Optional
* Datentyp: String
* Inhalt: Subfolder definiert einen Unterordner, in dem die Daten innerhalb des Volumes abgelegt werden sollen.
* Beispiel: `subfolder: "my-configmap-subfolder"`

## Config

Mit dem Feld `config` können globale und dogu-spezifische Konfigurationen der EcoSystem-Registry bearbeitet werden.
Außerdem ist es möglich Konfigurationen für Dogus verschlüsselt zu speichern.

### global

* Optional
* Datentyp: Array<configEntry>
* Inhalt: Setzt globale Konfigurationen.
* Beispiel:
```
config: 
  global: 
    - key: "gloval_key1"
      value: "global_value1"
    - key: "global_key2"
      value: "global_value2"
```

### dogus

* Optional
* Datentyp: map[string]Array<configEntry>
* Inhalt: Setzt Konfigurationen für Dogus.
* Beispiel:
```
config:
  dogus:
    postgresql:
      - key: "key1"
        value: "value1"
      - key: "key2"
        value: "value2"
```

### configEntry

#### key

* Datentyp: string
* Inhalt: Gibt den Schlüssel der Konfiguration an.
* Beispiel: `key: "key1"`

#### value

* Optional
* Datentyp: string
* Inhalt: Gibt den Wert der Konfiguration an.
* Beispiel: `value: "value1"`

#### absent

* Optional
* Datentyp: boolean
* Inhalt: Gibt an, ob der Konfigurationsschlüssel gelöscht werden soll. Wird dieser Wert gesetzt, muss lediglich der `key` mit angegeben werden.
* Beispiel: `absent: true`

#### sensitive

* Optional
* Datentyp: boolean
* Inhalt: Gibt an, ob die Konfiguration sensibel ist (und entsprechend gespeichert werden soll).
* Beispiel: `sensitive: true`.

#### secretRef

Die `SecretReference` kann genutzt werden, um statt eines `value` einen `key` aus einem Secret zu verwenden.

##### name

* Datentyp: string
* Inhalt: Gibt den Namen des Secrets an.
* Beispiel: `name: "my-secret"`

##### key

* Datentyp: string
* Inhalt: Gibt den Namen des Schlüssels aus dem Secret an.
* Beispiel: `key: "my-secret-key"`