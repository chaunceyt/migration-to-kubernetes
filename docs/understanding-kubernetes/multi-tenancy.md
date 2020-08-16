# Notes on Kubernetes multi tenancy best practices

- [Cluster multi-tenancy](https://cloud.google.com/kubernetes-engine/docs/concepts/multitenancy-overview)
- [Best practices for enterprise multi-tenancy](https://cloud.google.com/kubernetes-engine/docs/best-practices/enterprise-multitenancy)


> A multi-tenant cluster is shared by multiple users and/or workloads which are referred to as "tenants". 

The Kubernetes cluster that I am currently the administrator of, is a multi tenant cluster. We run all of our project's development environments. The patterns we're using for isolation is:

- Namespaces (Project isolation)
- Taints/Tolerations (Node isolation)
- RBAC
- Probes 

Planned:

- Pod Anti-affinity/Affinity
- ResourceQuota (resources: pods, loadbalancers, and nodeports)
- resources: requests

Understand

- Network isolation ensuring ingress and monitoring of namespaces via prometheus is not regressed
- Pod security policies
- [Eviction threshold](https://cloud.google.com/kubernetes-engine/docs/concepts/cluster-architecture#eviction_threshold) To determine how much memory is available for Pods.
- [Allocatable memory and CPU resources](https://cloud.google.com/kubernetes-engine/docs/concepts/cluster-architecture#memory_cpu) `Allocatable = Capacity - Reserved - Eviction Threshold`

Example `ResourceQuota` using promql we can query `kube_resourcequota{namespace=""}`

```
apiVersion: v1
kind: ResourceQuota
metadata:
  name: object-counts
spec:
  hard:
    pods: "35"
    services.loadbalancers: "0"
    services.nodeports: "0"
```

Example `NetworkPolicy` for network isolation.
Note: Add label to ingress-nginx, monitoring, and kube-system namespaces to allow communication with isolated project namespaces.

```
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: deny-all-namespace-network-policy
  namespace: [namespace]
spec:
  podSelector:
    matchLabels:
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          infra-plane: 'true' # allow communication from ingress-nginx, monitoring and kube-system.
  egress:
  - {}  
  policyTypes:
  - Egress
  - Ingress
```

# Tools

- [The Hierarchical Namespace Controller (HNC)](https://github.com/kubernetes-sigs/multi-tenancy/tree/master/incubator/hnc) [Demos](https://docs.google.com/document/d/1tKQgtMSf0wfT3NOGQx9ExUQ-B8UkkdVZB6m4o3Zqn64/edit#)

## Playing around with hnc

Using a Ubuntu EC2 instance in AWS with kind and docker installed and configured.

```
kind create cluster

HNC_VERSION=v0.5.1
kubectl apply -f https://github.com/kubernetes-sigs/multi-tenancy/releases/download/hnc-${HNC_VERSION}/hnc-manager.yaml
HNC_VERSION=v0.5.1
curl -L https://github.com/kubernetes-sigs/multi-tenancy/releases/download/hnc-${HNC_VERSION}/kubectl-hns -o ./kubectl-hns
chmod +x ./kubectl-hns
sudo mv kubectl-hns /usr/local/bin/

kubectl create ns team-c
kubectl create sa ns-admin -n team-c
kubectl create rolebinding ns-admin --clusterrole=admin --serviceaccount=team-c:ns-admin -n team-c
kubectl get rolebinding -n team-c
kubectl describe rolebinding ns-admin -n team-c
kubectl auth can-i --list --as=system:sericeaccount:team-c:ns-admin

kubectl hns create project-ci -n team-c
kubectl hns create design -n project-ci
kubectl hns create demos  -n team-c

kubectl hns tree team-c

# Example output of tree command.
team-c
├── demos
└── project-ci
    └── design

kubectl delete ns team-c

# Resulting error attempting to delete toplevel namespace

Error from server (Forbidden): admission webhook "namespaces.hnc.x-k8s.io" denied the request: Please set allowCascadingDelete first either in the parent namespace or in all the subnamespaces.
  Subnamespace(s) without allowCascadingDelete set: [demos project-ci].

# If you understand the risk set allowCascadingDelete
kubectl hns set team-c --allowCascadingDelete

# Now delete the namespace.
kubectl delete ns team-c
```

The Hierarchical Namespace Controller (HNC) doesn't address a need in the current environment I help manage.
