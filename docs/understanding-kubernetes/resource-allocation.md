# Notes on resource allocation

## Allocatable CPU & Memory

- Capacity determined by the eviction threshold to prevent system OOMs
- Resources reserved for the underlying VM such as operating system, system daemons like sshd, etc
- Resources needed to run Kubernetes such as  kubelet, container runtime, kube-proxy
- Resources for other Kubernetes-related add-ons such as  monitoring agents, and CNI plugins, etc
- Resources available for running applications

## Resource Ranges & Quotas

- Requests: lower bound on resource usage per workload
- Limits: upper bound on resource usage per workload

- Guaranteed = the memory and cpu request and limits matches

```
    resources:
      limits:
        cpu: 700m
        memory: 200Mi
      requests:
        cpu: 700m
        memory: 200Mi
```

- Burstable = requests < limit

```
    resources:
      limits:
        memory: "200Mi"
      requests:
        memory: "100Mi"
```


- Best Effort = pod doesn't have memory or cpu limits or requests

```
 resources: {}

```
## Resources

- https://kubernetes.io/docs/tasks/administer-cluster/reserve-compute-resources/#node-allocatable
- https://cloud.google.com/kubernetes-engine/docs/concepts/cluster-architecture#node_allocatable
- https://kubernetes.io/docs/tasks/configure-pod-container/quality-service-pod/
- https://docs.microsoft.com/en-us/azure/aks/concepts-clusters-workloads#resource-reservations
- https://github.com/kubernetes-sigs/descheduler
- https://github.com/kubernetes/autoscaler/tree/master/vertical-pod-autoscaler
- https://github.com/FairwindsOps/goldilocks/

