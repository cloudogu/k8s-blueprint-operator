#!/bin/bash
etcdClient="$(kubectl get pods -l=app.kubernetes.io/name=etcd-client -o name 2>&1 | head -n 1)"

kubectl exec --namespace ecosystem "$etcdClient" -- etcdctl set "/config/_global/global_delete" "doesnotmatter" && \
kubectl exec --namespace ecosystem "$etcdClient" -- etcdctl set "/config/postgresql/example_key_to_remove" "doesnotmatter" && \
kubectl exec --namespace ecosystem "$etcdClient" -- etcdctl set "/config/postgresql/example_key_to_remove_encrypted" "doesnotmatter"