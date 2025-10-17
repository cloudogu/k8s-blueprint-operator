# How to Configure a Dogu

Configuration for dogus and the entire Cloudogu EcoSystem can be managed declaratively within the `Blueprint` resource, under the `spec.blueprint.config` section. This allows you to version-control your configuration alongside your application definitions.

## Configuration Scopes

There are two scopes for configuration:

1.  **Global Configuration (`config.global`):** These key-value pairs apply to the entire EcoSystem and are accessible to all dogus.
2.  **Dogu-Specific Configuration (`config.dogus`):** These key-value pairs are targeted at a single dogu.

```yaml
# blueprint.yaml
spec:
  blueprint:
    config:
      global:
        # Global key-values go here
      dogus:
        # Dogu-specific key-values go here
        <dogu-name>:
          # ...
```

---

## Setting Plaintext Configuration

For standard configuration, you can provide values directly in the blueprint. The operator will store these in a Kubernetes `ConfigMap` and mount them into the appropriate dogu.

### Example

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

- The global key `admin/mail` is set to `admin@my-ces.com`.
- For the `redmine` dogu, the key `logging/root` is set to `DEBUG`.

---

## Setting Sensitive Configuration

For sensitive data like passwords or API tokens, you must not place the value directly in the blueprint. Instead, you source it from an existing Kubernetes `Secret`.

To do this, you mark the configuration as `sensitive: true` and provide a `secretRef` pointing to the name of the `Secret` and the `key` within that secret.

**Prerequisite:** The Kubernetes `Secret` you reference must exist before you apply the blueprint.

### Example

First, ensure a secret exists:
```yaml
# my-ldap-secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: ldap-credentials
stringData:
  password: "super-secret-password"
```

Then, reference it in your blueprint:
```yaml
config:
  dogus:
    usermgmt:
      - key: "ldap/password"
        sensitive: true
        secretRef:
          name: "ldap-credentials" # Name of the Secret
          key: "password"         # Key within the Secret
```

---

## Deleting a Configuration Key

To remove a configuration key from the EcoSystem, you mark it as `absent: true`.

### Deleting a Plaintext Key

For standard, non-sensitive keys, you only need to provide the `key` and `absent: true`.

```yaml
config:
  global:
    # This will remove the 'logging/root' key from the global config
    - key: "logging/root"
      absent: true
```

### Deleting a Sensitive Key

To remove a sensitive configuration key, you must specify **both** `absent: true` and `sensitive: true`. This tells the operator to look for the key in the dogu's secret storage instead of its public configmap.

```yaml
config:
  dogus:
    usermgmt:
      # This will remove the sensitive 'ldap/password' key
      - key: "ldap/password"
        absent: true
        sensitive: true
```

---

## Configuration Rules and Validations

The blueprint operator enforces several rules to ensure configuration is valid. An invalid blueprint will fail validation and will not be applied. The following combinations are **not allowed**:

- A configuration entry **cannot** have both a `value` and be `absent: true`.
- A configuration entry **cannot** have both a `secretRef` and be `absent: true`.
- A configuration entry **cannot** have both a `value` and a `secretRef`.
- A configuration entry with a `secretRef` **must** also have `sensitive: true`.
- A configuration entry with `sensitive: true` **must** use `secretRef` and cannot have a plaintext `value`.

---

## Full Example

Here is a blueprint that demonstrates all configuration types:

```yaml
apiVersion: k8s.cloudogu.com/v2
kind: Blueprint
metadata:
  name: blueprint-with-config
spec:
  displayName: "Configuration Example"
  blueprint:
    dogus:
      - name: "official/usermgmt"
        version: "2.8.1-1"
      - name: "official/redmine"
        version: "5.1.2-3"
    config:
      global:
        # Set a global plaintext value
        - key: "admin/mail"
          value: "admin@my-ces.com"
        # Remove a global key
        - key: "old/global/key"
          absent: true
      dogus:
        usermgmt:
          # Set a sensitive value from a secret
          - key: "ldap/password"
            sensitive: true
            secretRef:
              name: "ldap-credentials"
              key: "password"
          # Remove a sensitive key
          - key: "old/ldap/password"
            absent: true
            sensitive: true
        redmine:
          # Set a dogu-specific plaintext value
          - key: "logging/root"
            value: "DEBUG"
```
