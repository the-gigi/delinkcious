# Global Kubernetes resources and plicies

This directory contains manifests for resources and policies that apply across the cluster

## NATS cluster

This is the only required service. Install it using the following commands:

```
$ kubectl apply -f https://github.com/nats-io/nats-operator/releases/download/v0.4.2/00-prereqs.yaml
$ kubectl apply -f https://github.com/nats-io/nats-operator/releases/download/v0.4.2/10-deployment.yaml
$ kubectl create -f nats_cluster.yaml
```


