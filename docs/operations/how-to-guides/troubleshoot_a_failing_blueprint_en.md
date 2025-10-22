# How to Troubleshoot a Failing Blueprint

When a Blueprint doesn't behave as expected, the `Blueprint` custom resource itself is the best place to start your investigation. It contains detailed status information, a plan of action, and is linked to events that provide a clear picture of the operator's activity.

This guide will walk you through the key areas to inspect when troubleshooting.

## 1. Check the Blueprint's Status Conditions

The most immediate source of information is the `status.conditions` field of your Blueprint resource. You can view this by using `kubectl describe` or `kubectl get blueprint -o yaml`.

```bash
kubectl describe blueprint <your-blueprint-name> -n <your-namespace>
```

Look for the `Status` section. Here are the most important conditions and what they mean:

- **`Valid`**: If this condition is `False`, it means there is a structural or logical error in your blueprint definition. The reason and message will often tell you exactly what's wrong (e.g., a syntax error or a missing dependency for a dogu).

- **`Executable`**: This will be `False` if the calculated changes are not allowed. The most common reason is an attempted dogu downgrade, which is blocked by default. The message associated with this condition will explain the problematic change.

- **`EcosystemHealthy`**: This indicates whether the operator is waiting for the ecosystem to become healthy before applying changes. If it's `False`, it means one or more dogus are not in a ready state.

- **`Completed`**: This shows if the blueprint has been fully applied. If it's `False` long after you've applied it, it means the operator is still working or is stuck.

- **`LastApplySucceeded`**: This is a critical condition for troubleshooting. If an operation fails (like applying a configmap or installing a dogu), this condition will become `False`. **Crucially, it holds the last error message** and persists across multiple reconciliation loops until the blueprint is successfully completed. This allows you to see the root cause of a failure even if the operator is retrying.

Start by looking for any condition that is `False` and read its associated `message` for details.

## 2. Analyze the StateDiff

If the blueprint is valid and the ecosystem is healthy, but changes aren't being applied as you expect, the `status.stateDiff` field is your next stop. This field shows the exact plan the operator has calculated by comparing the desired state (your blueprint) with the actual state of the cluster.

You can view it with `kubectl get`:

```bash
kubectl get blueprint <your-blueprint-name> -n <your-namespace> -o yaml
```

The `stateDiff` will show you exactly what dogus and configurations the operator intends to add, upgrade, remove, or modify. This is useful for spotting:
- **Unintended Changes**: Does the diff include changes you didn't expect? This might point to an issue with your `blueprint` or `blueprintMask`.
- **Problematic Operations**: The diff might explicitly show a planned downgrade, which would explain why the `Executable` condition is `False`.

## 3. Inspect Events

Kubernetes events provide a chronological log of what the operator has been doing. When you run `kubectl describe blueprint <your-blueprint-name>`, you also get a list of associated events at the bottom.

These events will show you among other things:
- When a reconciliation loop started.
- The outcome of validation checks.
- The beginning and end of the apply phase.
- Any errors encountered while interacting with other resources.

## 4. Check the Operator Logs

For the most detailed information, you need to check the logs of the `k8s-blueprint-operator` pod itself. This is where you'll find detailed error messages and stack traces that can pinpoint the exact line of code where a failure occurred.

1.  **Find the operator pod:**
    ```bash
    kubectl get pods -n <operator-namespace> -l app.kubernetes.io/name=k8s-blueprint-operator
    ```
    *(The namespace is typically `ecosystem` or where you installed the component).*

2.  **Stream the logs:**
    ```bash
    kubectl logs -f <operator-pod-name> -n <operator-namespace>
    ```

If the default `info` level logs aren't sufficient, you can increase the verbosity to `debug` or `trace`. This is done by modifying the `Component` resource for the operator. For detailed instructions, see the [Operator Configuration documentation](../reference/operator_configuration_en.md).
