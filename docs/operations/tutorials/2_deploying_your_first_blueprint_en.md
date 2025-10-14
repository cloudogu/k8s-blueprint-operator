# Deploying Your First Blueprint

This tutorial will guide you through deploying your first application (a `dogu`) using a simple Blueprint. A Blueprint is a custom resource that declaratively defines the desired state of applications and their configurations in your Cloudogu EcoSystem.

## Goal

By the end of this tutorial, you will have deployed the `mysql` dogu using a Blueprint.

## 1. Create the Blueprint file

First, create a new YAML file named `my-first-blueprint.yaml` and add the following content:

```yaml
apiVersion: k8s.cloudogu.com/v2
kind: Blueprint
metadata:
  name: my-first-blueprint
spec:
  displayName: "My First Blueprint"
  blueprint:
    dogus:
      - name: "official/mysql"
        version: "8.4.5-4"
```

### What does this file do?

*   `kind: Blueprint`: This tells Kubernetes that the resource is a Blueprint, which the `k8s-blueprint-operator` knows how to handle.
*   `metadata.name`: This gives our Blueprint a unique name, `my-first-blueprint`.
*   `spec.blueprint.dogus`: This is the list of applications we want to be present in the system. Here, we are defining one dogu:
    *   `name: "official/mysql"`: Specifies the `mysql` dogu from the `official` namespace.
    *   `version: "8.4.5-4"`: Specifies the exact version we want to install.

## 2. Apply the Blueprint

Now, use `kubectl` to apply this resource to your cluster. This will trigger the `k8s-blueprint-operator` to start reconciling the state.

```bash
kubectl apply -f my-first-blueprint.yaml -n ecosystem
```

## 3. Check the Status

The operator will now work to achieve the state defined in the Blueprint. You can check the progress by describing the Blueprint resource:

```bash
kubectl describe blueprint my-first-blueprint
```

Look at the `Status` and `Events` sections in the output. You should see events indicating that the operator is validating the blueprint and applying the changes. Once finished, the status will show that the blueprint is `Completed`.

## Conclusion

Congratulations! You have successfully used the `k8s-blueprint-operator` to deploy your first dogu. You can now build on this by adding more dogus or providing configuration as shown in the other guides.
