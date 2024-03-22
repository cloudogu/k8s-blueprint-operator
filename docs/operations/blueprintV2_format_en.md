# Blueprint format

The blueprint offers the option of adding or removing Dogus, Components and configurations.
Special configurations can also be defined for Dogus and components.
These are written to the corresponding Dogu and component CRs and are not saved in the EcoSystem registry.

All fields of the blueprint are described below and illustrated with examples.

## BlueprintApi

* Not optional
* Data type: string
* Content: The `blueprintApi` field specifies the API version of the blueprint.
* Example: `"blueprintApi": "v2"`

## Dogus

* Not optional
* Data type: Array<Dogu>
* Content: The `dogus` field is a list of Dogus and describes the status of the Dogus in the system.
* Example:
```
"dogus": [
  {
    "name": "official/postgresql",
    "targetState": "present",
    "version": "12.15-2"
  }
]
```

### Dogu

A Dogu can contain the following fields:

#### Name

* Not optional
* Data type: string
* Content: Specifies the name including namespace of the Dogu.
* Example: `"name": "official/cas"`

#### TargetState

* Not optional
* Data type: string
* Content: Specifies whether a Dogu should be present or not.
* Example: `"targetState": "present"` or `"targetState": "absent"`

#### Version

* Optional for `targetState=absent`. Not optional for `targetState=present`.
* Data type: string
* Content: Specifies the version of the Dogu.
* Example: `"version": "12.15-2"`

#### PlatformConfig

The `platformConfig` field offers the option of transferring specific configurations for the execution platform (e.g. Kubernetes).
This configuration can be used to define resources and reverse proxy configurations.

##### Resource.minVolumeSize

* Optional
* Data type: string
* Content: Specifies the minimum volume size of a Dogu. If the current volume is smaller, a volume increase is performed. The unit must be specified with the binary prefix (e.g. `Mi` or `Gi`).
* Example:
```
"dogus": [
  {
    "name": "official/nexus",
    "targetState": "present",
    "version":"3.59.0-2",
    "platformConfig": {
      "resource": {
        "minVolumeSize": "5Gi"
      }
    }
  }
]
```

> The Dogu-Operator creates Dogus with 2Gi volumes. The Nexus Dogu requires a larger volume and must be configured via this entry.

##### ReverseProxy.maxBodySize

* Optional
* Data type: string
* Content: Specifies the maximum file size for HTTP request (`0`=unlimited). The unit must be specified with the decimal prefix (e.g. `M` or `G`).
* Example:
```
"dogus": [
  {
    "name": "official/nexus",
    "targetState": "present",
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
* Data type: string
* Content: Defines a rewrite target.
* Example:
```
"dogus": [
  {
    "name": "official/postgresql",
    "targetState": "present",
    "version": "12.15-2",
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
* Data type: string
* Content: Adds any additional proxy configuration.
* Example:
```
"dogus": [
  {
    "name": "official/postgresql",
    "targetState": "present",
    "version": "12.15-2",
    "platformConfig": {
      "reverseProxy": {
        "additionalConfig": "<config>"
      }
    }
  }
]
```

## Components

* Not optional
* Data type: Array
* Contents: The `components` field is a list of components and describes the status of the components in the system.
* Example:
```
"components": [
  {
    "name": "k8s/k8s-dogu-operator",
    "targetState": "present",
    "version": "1.0.1"
  },
  {
    "name": "k8s/k8s-dogu-operator-crd",
    "targetState": "present",
    "version": "1.0.1"
  }
}
```

### Component

A component can contain the following fields:

#### Name

* Not optional
* Data type: string
* Content: Specifies the name including namespace of the component.
* Example: `"name": "k8s/k8s-dogu-operator"`

#### TargetState

* Not optional
* Data type: string
* Content: Specifies whether a component should be present or not.
* Example: `"targetState": "present"` or `"targetState": "absent"`

#### Version

* Optional for `targetState=absent`. Not optional for `targetState=present`.
* Data type: string
* Content: Specifies the version of the component.
* Example: `"version": "12.15-2"`

#### DeployConfig

The `deployConfig` field offers the option of transferring specific configurations for the deployment of a component.
This configuration can be used, for example, to define the component CR or any helm values.

##### deployConfig.deployNamespace

* Optional
* Data type: string
* Content: Specifies the namespace in which the component is to be installed. This configuration is currently only required for the component `k8s/longhorn`.
* Example:
```
"components": [
  {
    "name": "k8s/k8s-longhorn",
    "targetState": "present",
    "version": "1.5.1-4",
    "deployConfig":{
      "deployNamespace": "longhorn-system"
    }
  }
]
```

##### deployConfig.overwriteConfig

* Not optional
* Data type: string
* Content: Defines additional configurations (Helm values) for the component.
* Example:
```
"components": [
  {
    "name": "k8s/k8s-longhorn",
    "targetState": "present",
    "version": "1.5.1-4",
    "deployConfig":{
      "overwriteConfig":{
        "longhorn":{
          "defaultSettings":{
            "backupTarget": "s3://longhorn@dummyregion/",
            "backupTargetCredentialSecret": "longhorn-backup-target"
          }
        }
      }
    }
  }
]
```

## Config

The `config` field can be used to edit global and Dogu-specific configurations of the EcoSystem registry.
It is also possible to save configurations for Dogus in encrypted form.

### global.present

* Optional
* Data type: map[string]string
* Content: Sets global configurations.
* Example:
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
* Data type: Array<string>
* Content: Removes global configurations.
* Example:
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
* Data type: map[string]string
* Content: Sets configurations for Dogus.
* Example:
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
* Data type: Array<string>
* Content: Removes configurations from Dogus.
* Example:
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
* Data type: map[string]string
* Content: Sets encrypted configurations for Dogus.
* Example:
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
* Data type: Array<string>
* Content: Removes encrypted configurations from Dogus.
* Example:
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
