# Blueprint-Format

Das Blueprint bietet die Möglichkeit Dogus, Components und Konfigurationen hinzuzufügen oder zu entfernen.
Für Dogus und Components können ebenfalls spezielle Konfigurationen definiert werden.
Diese werden an die entsprechenden Dogu- und Component-CRs geschrieben und werden nicht in der EcoSystem-Registry gespeichert.

Folgend werden alle Felder des Blueprint beschrieben und mit Beispielen veranschaulicht.

## BlueprintApi

* Nicht Optional
* Datentyp: string
* Inhalt: Das Feld `blueprintApi` gibt die API-Version des Blueprints an.
* Beispiel: `"blueprintApi": "v2"`

## Dogus

* Nicht Optional
* Datentyp: Array<dogu>
* Inhalt: Das Feld `dogus` ist eine Liste von Dogus und beschreibt den Zustand der Dogus im System.
* Beispiel: 
```
"dogus": [
  {
    "name":"official/postgresql",
    "targetState":"present",
    "version":"12.15-2"
  }
]
```

### Dogu

Ein Dogu kann folgende Felder beinhaltet:

#### Name

* Nicht Optional
* Datentyp: string
* Inhalt: Gibt den Namen inklusive Namespace des Dogus an.
* Beispiel: `"name": "official/cas"`

#### TargetState

* Nicht Optional
* Datentyp: string
* Inhalt: Gibt an, ob ein Dogu vorhanden oder nicht vorhanden sein soll.
* Beispiel: `"targetState": "present"` oder `"targetState": "absent"`

#### Version

* Bei `targetState=absent` optional. Bei `targetState=present` nicht optional.
* Datentyp: string
* Inhalt: Gibt die Version des Dogus an.
* Beispiel: `"version": "12.15-2"`

#### PlatformConfig

Das Feld `platformConfig` bietet die Möglichkeit für die Ausführungsplattform (z.B. Kubernetes) spezifische Konfigurationen zu übergeben.
Mit dieser Konfiguration können Ressourcen und Reverse-Proxy-Konfigurationen definiert werden.

##### Resource.minVolumeSize

* Optional
* Datentyp: string
* Inhalt: Gibt die minimale Volume-Size eines Dogus an. Falls das aktuelle Volume kleiner ist, wird eine Volume-Vergrößerung durchgeführt. Die Einheit muss mit dem Binär-Prefix angegeben werden (z.B. `Mi` oder `Gi`).
* Beispiel:
```
"dogus": [
  {
    "name":"official/nexus",
    "targetState":"present",
    "version":"3.59.0-2",
    "platformConfig": {
      "resource": {
        "minVolumeSize": "5Gi"
      }
    }
  }
]
```

> Der Dogu-Operator erstellt Dogus mit 2Gi Volumes. Das Nexus-Dogu benötigt ein größeres Volume und muss
> über diesen Eintrag konfiguriert werden.

##### ReverseProxy.maxBodySize

* Optional
* Datentyp: string
* Inhalt: Gibt die maximale Dateigröße für HTTP-Request an (`0`=unbegrenzt). Die Einheit muss mit dem Dezimal-Prefix angegeben werden (z.B. `M` oder `G`).
* Beispiel:
```
"dogus": [
  {
    "name":"official/nexus",
    "targetState":"present",
    "version":"3.59.0-2",
    "platformConfig": {
      "reverseProxy": {
        "maxBodySize": "1G"
      }
    }
  }
]
```

##### ReverseProxy.rewriteTarget

* Optional
* Datentyp: string
* Inhalt: Definiert ein Rewrite-Target.
* Beispiel:
```
"dogus": [
  {
    "name":"official/postgresql",
    "targetState":"present",
    "version":"12.15-2",
    "platformConfig": {
      "reverseProxy": {
        "rewriteTarget": "/"
      }
    }
  }
]
```

##### ReverseProxy.additionalConfig

* Optional
* Datentyp: string
* Inhalt: Fügt beliebige zusätzliche Proxy-Konfiguration hinzu.
* Beispiel:
```
"dogus": [
  {
    "name":"official/postgresql",
    "targetState":"present",
    "version":"12.15-2",
    "platformConfig": {
      "reverseProxy": {
        "additionalConfig": "<config>"
      }
    }
  }
]
```

## Components

* Nicht Optional
* Datentyp: Array
* Inhalt: Das Feld `components` ist eine Liste von Components und beschreibt den Zustand der Components im System.
* Beispiel:
```
"components": [
  {
    "name":"k8s/k8s-dogu-operator",
    "targetState":"present",
    "version":"1.0.1"
  },
  {
    "name":"k8s/k8s-dogu-operator-crd",
    "targetState":"present",
    "version":"1.0.1"
  }
}
```

### Component

Eine Component kann folgende Felder beinhaltet:

#### Name

* Nicht Optional
* Datentyp: string
* Inhalt: Gibt den Namen inklusive Namespace der Component an.
* Beispiel: `"name": "k8s/k8s-dogu-operator"`

#### TargetState

* Nicht Optional
* Datentyp: string
* Inhalt: Gibt an, ob eine Component vorhanden oder nicht vorhanden sein soll.
* Beispiel: `"targetState": "present"` oder `"targetState": "absent"`

#### Version

* Bei `targetState=absent` optional. Bei `targetState=present` nicht optional.
* Datentyp: string
* Inhalt: Gibt die Version der Component an.
* Beispiel: `"version": "12.15-2"`

#### DeployConfig

Das Feld `platformConfig` bietet die Möglichkeit für das Deployment einer Component spezifische Konfigurationen zu übergeben.
Mit dieser Konfiguration können zum Beispiel die Component-CR oder beliebige Helm-Values definiert werden.

##### deployConfig.deployNamespace

* Optional
* Datentyp: string
* Inhalt: Gibt den Namespace an in den die Component installiert werden soll. Diese Konfiguration wird bisher nur für die Component `k8s/longhorn` benötigt.
* Beispiel: 
```
"components": [
  {
    "name":"k8s/k8s-longhorn",
    "targetState":"present",
    "version":"1.5.1-4",
    "deployConfig":{
      "deployNamespace":"longhorn-system"
    }
  }
]
```

##### deployConfig.overwriteConfig

* Nicht Optional
* Datentyp: string
* Inhalt: Definiert zusätzliche Konfigurationen (Helm-Values) für die Component.
* Beispiel:
```
"components": [
  {
    "name":"k8s/k8s-longhorn",
    "targetState":"present",
    "version":"1.5.1-4",
    "deployConfig":{
      "overwriteConfig":{
        "longhorn":{
          "defaultSettings":{
            "backupTarget":"s3://longhorn@dummyregion/",
            "backupTargetCredentialSecret":"longhorn-backup-target"
          }
        }
      }
    }
  }
]
```

## Config

Mit dem Feld `config` können globale und dogu-spezifische Konfigurationen der EcoSystem-Registry bearbeitet werden.
Außerdem ist es möglich Konfigurationen für Dogus verschlüsselt zu speichern.

### global.present

* Optional
* Datentyp: map[string]string
* Inhalt: Setzt globale Konfigurationen.
* Beispiel:
```
"config": {
  "global": {
    "present": {
      "global_key1": "global_value1",
      "global_key2": "global_value2"
    }
  }
}
```

### global.absent

* Optional
* Datentyp: Array<string>
* Inhalt: Entfernt globale Konfigurationen.
* Beispiel:
```
"config": {
  "global": {
    "absent": [
      "global_key1", "global_key2"
    ]
  }
}
```

### dogus.config.present

* Optional
* Datentyp: map[string]string
* Inhalt: Setzt Konfigurationen für Dogus.
* Beispiel:
```
"config": {
  "dogus": {
    "postgresql": {
      "config": {
        "present": {
          "key1": "value1",
          "key2": "value2"
        }
      }
    }
  }
}
```

### dogus.config.absent

* Optional
* Datentyp: Array<string>
* Inhalt: Entfernt Konfigurationen von Dogus.
* Beispiel:
```
"config": {
  "dogus": {
    "postgresql": {
      "config": {
        "absent": [
          "key1", "key2"
        ]
      }
    }
  }
}
```

### dogus.sensitiveConfig.present

* Optional
* Datentyp: map[string]string
* Inhalt: Setzt verschlüsselte Konfigurationen für Dogus.
* Beispiel:
```
"config": {
  "dogus": {
    "postgresql": {
      "sensitiveConfig": {
        "present": {
          "key1": "value1",
          "key2": "value2"
        }
      }
    }
  }
}
```

### dogus.sensitiveConfig.absent

* Optional
* Datentyp: Array<string>
* Inhalt: Entfernt verschlüsselte Konfigurationen von Dogus.
* Beispiel:
```
"config": {
  "dogus": {
    "postgresql": {
      "sensitiveConfig": {
        "absent": [
          "key1", "key2"
        ]
      }
    }
  }
}
```
