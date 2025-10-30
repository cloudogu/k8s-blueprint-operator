# Upgrading a Dogu with a Blueprint

In the previous tutorial, you deployed your first application (`mysql`) using a blueprint. This guide will show you how to upgrade that application to a new version simply by changing your blueprint definition.

This declarative approach is one of the core strengths of the operator. You define the desired state, and the operator handles the steps required to get there.

### Prerequisites

This tutorial assumes you have completed the "Deploying Your First Blueprint" tutorial and have the `my-first-blueprint.yaml` file available.

## Step 1: Modify the Blueprint Version

To upgrade an application, you just need to specify the new version you want in your blueprint file.

Open your `my-first-blueprint.yaml` file and change the version of the `mysql` dogu from `8.4.5-4` to a newer version, for example, `8.4.6-1`.

```diff
# my-first-blueprint.yaml

apiVersion: k8s.cloudogu.com/v2
kind: Blueprint
metadata:
  name: my-first-blueprint
spec:
  displayName: "My First Blueprint"
  blueprint:
    dogus:
      - name: "official/mysql"
-       version: "8.4.5-4"
+       version: "8.4.6-1"
```

## Step 2: Apply the Updated Blueprint

Now, apply the modified file to your cluster using the same command as before. The operator will detect the change to the `Blueprint` resource and start a new reconciliation loop.

```bash
kubectl apply -f my-first-blueprint.yaml -n ecosystem
```

The operator compares the new desired state from your file with the current state of the cluster and sees that the version of `mysql` is different. It will automatically determine that an upgrade is required.

## Step 3: Observe the Upgrade

If you are fast enough, you can see what the operator has planned by inspecting the `StateDiff` in the blueprint's status. This provides a transparent view of the actions the operator will take.

```bash
kubectl get blueprint my-first-blueprint -n ecosystem -o yaml
```

In the `status.stateDiff.doguDiffs` section, you might now see an entry for `mysql` with `upgrade` listed under `neededActions` (only until the operator has finished the upgrade).

To watch the upgrade happen in real-time, you can watch the `dogu` resource itself:

```bash
kubectl get dogu mysql --watch -n ecosystem
```

You will see the status of the dogu change as the operator performs the upgrade to the new version.

## Alternative: Using `kubectl edit`

For quick, interactive changes, you can also edit the blueprint directly in the cluster without modifying a local file:

```bash
kubectl edit blueprint my-first-blueprint -n ecosystem
```

This will open the blueprint resource in your default editor. You can change the version, save, and close the file. The operator will detect the change and start the upgrade immediately.

**Note:** While `kubectl edit` is convenient for development or testing, we strongly recommend keeping your blueprint YAML files in a version control system (like Git) for production environments. This `GitOps` approach allows you to track, review, and roll back changes to your application landscape safely.

## Conclusion

Congratulations! You have now seen how the `k8s-blueprint-operator` simplifies application lifecycle management. By declaratively managing your application versions in a single file, you can perform upgrades reliably and predictably.
