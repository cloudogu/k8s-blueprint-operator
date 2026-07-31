# Blueprint format

The blueprint offers the option of adding or removing Dogus and configurations.
Special configurations can also be defined for Dogus.
These are written to the corresponding Dogu CRs and are not saved in the EcoSystem registry.

All fields of the blueprint are described below and illustrated with examples.

## BlueprintApi

* Required
* Data type: string
* Content: The `blueprintApi` field specifies the API version of the blueprint.
* Example: `blueprintApi: "v3"`

## Dogus

* Required
* Data type: Array<Dogu>
* Content: The `dogus` field is a list of Dogus and describes the status of the Dogus in the system.
* Example:
```
dogus:
  - name: official/mysql
    version: 8.4.6-1
```

### Dogu

A Dogu can contain the following fields:

#### Name

* Required
* Data type: string
* Content: Specifies the name including the namespace of the Dogu.
* Example: `name: "official/cas"`

#### Absent

* Data type: boolean
* Content: Specifies whether a Dogu should be ommitted.
* Example: `absent: true`

#### Version

* Optional for `absent=true`. Otherwise `version` is required.
* Data type: string
* Content: Specifies the version of the Dogu.
* Example: `version: "12.15-2"`

#### PlatformConfig

The `platformConfig` field offers the option of transferring specific configurations for the execution platform (e.g. Kubernetes).
This configuration can be used to define resources and reverse proxy configurations.

##### Resource.minVolumeSize

* Optional
* Data type: string
* Content: Specifies the minimum volume size of a Dogu. If the current volume is smaller, a volume increase is performed. The unit must be specified with the binary prefix (e.g. `Mi` or `Gi`).
* Example:
```
dogus: 
  name: "official/nexus"
  version: "3.59.0-2"
  platformConfig: 
    resource: 
      minVolumeSize: "5Gi"
```

> The Dogu-Operator creates Dogus with 2Gi volumes. The Nexus Dogu requires a larger volume and must be configured via this entry.

> Shrinking volumes is not supported.

##### ReverseProxy.maxBodySize

* Optional
* Data type: string
* Content: Specifies the maximum file size for HTTP request (`0`=unlimited). The unit must be specified with the decimal prefix (e.g. `M` or `G`).
* Example:
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
* Data type: string
* Content: Defines a rewrite target.
* Example:
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
* Data type: string
* Content: Adds any additional proxy configuration.
* Example:
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
* Data type: Array
* Contents: The field `additionalDataMounts` is a list that contains elements of type `AdditionalMount`.
* Example:
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

* Data type: AdditionalMount
* Contents: An `AdditionalMount` defines a source and a target volume to mount files in a dogu.

###### AdditionalMount.SourceType

* Mandatory field
* Data type: Enum <ConfigMap; Secret>
* Content: SourceType specifies the type of source.
  Valid options are:
    - ConfigMap - Data stored in a kubernetes ConfigMap.
    - Secret - Data stored in a kubernetes Secret.
* Example: `sourceType: ConfigMap`

###### AdditionalMount.Name

* Mandatory field
* Data type: String
* Content: Name is the name of the data source.
* Example: `name: my-configmap`

###### AdditionalMount.Volume

* Mandatory field
* Data type: String
* Content: Volume is the name of the volume in which the data should be mounted. This is defined in the respective dogu.json file.
* Example: `volume: importHistory`

###### AdditionalMount.Subfolder

* Optional
* Data type: String
* Content: Subfolder defines a subfolder in which the data is to be stored within the volume.
* Example: `subfolder: "my-configmap-subfolder"`

## Config

The `config` field can be used to edit global and Dogu-specific configurations of the EcoSystem registry.
It is also possible to save configurations for Dogus in encrypted form.

### global

* Optional
* Data type: Array<configEntry>
* Content: Sets global configurations.
* Example:
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
* Data type: map[string]Array<configEntry>
* Content: Sets configurations for Dogus.
* Example:
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

* Data type: string
* Content: Specifies the configuration key.
* Example: `key: "key1"`

#### value

* Optional
* Data type: string
* Content: Specifies the configuration value.
* Example: `value: "value1"`

#### absent

* Optional
* Data type: boolean
* Content: Specifies whether the configuration key should be deleted. If this value is set, only the `key` needs to be specified.
* Example: `absent: true`

#### sensitive

* Optional
* Data type: boolean
* Content: Specifies whether the configuration is sensitive (and should be stored accordingly).
* Example: `sensitive: true`.

#### secretRef

The `SecretReference` can be used to use a `key` from a secret instead of a `value`.

##### name

* Data type: string
* Content: Specifies the secret name.
* Example: `name: "my-secret"`

##### key

* Data type: string
* Content: Specifies the name of the key from the secret.
* Example: `key: "my-secret-key"`