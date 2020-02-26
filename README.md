# argocdprogs-controller

This repository implements a simple controller for scheduling ArgoCD app sync, then moving to another app after a period if everything went ok and so on for multiple hops. The idea can be used mainly for microservices using async communication, where a service mesh is not useful. This way we can still provide some sort of progressive delivery.

## Details

The argocdprog-controller uses [client-go library](https://github.com/kubernetes/client-go/tree/master/tools/cache) extensively and a CRD called ArgoCDProg. In the CR we can define the first app to sync at a specific time and then a time period to wait until getting to the next app to sync.

For now it has support just for one ArgoCD API Server, with the server/sync token passed as an environment variable. The plan is to create a secret with entries for each api server and server token, pass the api server for each app sync in the CR and read the token from the secret.
Even more down the road I want to be able to read metrics from prometheus related to specific app and if they are under a specific threshold the sync chain can continue. Other monitoring platforms will also be considered.

## Purpose

This can be used either for starting a sync from dev > staging > sandbox > production (and maybe other environments), or it can also be used for propagating a change to multiple customers, when each has specific environment.

## Running

**Prerequisite**: to be identified

```sh
# assumes you have a working kubeconfig, not required if operating in-cluster
go build -o argocdprog .
./argocdprog -kubeconfig=$HOME/.kube/config

# create the ArgoCDProg CRD
kubectl create -f manifests/crd/argocdprog.yaml

# create an instance of type ArgoCDProg
kubectl create -f manifests/crd/example-argocdprod.yaml

# check ArgoCDProg resources for their status
kubectl get argocdprogs
```

## Cleanup

You can clean up the created CustomResourceDefinition with:

    kubectl delete crd argocdprogs.knapps.lcostea.io
