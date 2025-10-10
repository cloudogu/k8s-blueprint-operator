# How to Apply or Update a Blueprint

You can apply a blueprint by creating a `Blueprint` resource in the cluster. The `k8s-blueprint-operator` will then automatically apply the changes to your Cloudogu EcoSystem.

## Steps

1.  **Create a YAML file** containing your Blueprint definition. See the example below for a reference.
2.  **Apply it to the cluster** using `kubectl`.

    ```bash
    kubectl apply -f your-blueprint-file.yaml
    ```

**Note:** Only one blueprint is permitted per namespace. If you `apply` a blueprint with the same name as an existing one, the operator will update the deployment to match the new definition.

## Full Blueprint Example

Here is an example of a `Blueprint` resource that demonstrates several features, including defining dogus, setting configuration, and using a blueprint mask.

```yaml
apiVersion: k8s.cloudogu.com/v2
kind: Blueprint
metadata:
  labels:
    app: ces
    app.kubernetes.io/name: k8s-blueprint-lib
  name: blueprint-sample
spec:
  displayName: "Blueprint Sample v6.834"
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

This example does the following:

*   Defines the `official/mysql` dogu at version `8.4.6-1`.
*   Sets and removes several configuration values for the `mysql` dogu and globally.
*   Uses a `secretRef` to source a sensitive value from a Kubernetes secret.
*   Uses a `blueprintMask` to mark the `mysql` dogu as `absent`, effectively preventing it from being installed by this blueprint.

For a detailed breakdown of all possible fields, see the official [Blueprint Format Documentation](https://github.com/cloudogu/k8s-blueprint-lib/blob/develop/docs/operations/blueprintV2_format_en.md).