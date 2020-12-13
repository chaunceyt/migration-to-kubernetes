# Kubernetes scheduler notes

Responsible for "making sure that Pods are matched to Nodes so that Kubelet can run them"

[Kubernetes Scheduler](https://kubernetes.io/docs/concepts/scheduling-eviction/kube-scheduler/)
[Scheduling policies](https://kubernetes.io/docs/reference/scheduling/policies/)
[Scheduling Framework](https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/)


**Scheduling stage**

- Queue
- Prefilter
- Filter
- PreScore
- Score
- Normalise score
- Notifier
- Binding policies

**Binding stage**

- WaitOnPermit
- PreBind
- Bind
- PostBind


Filter stage - filter out nodes

[Predicates](https://kubernetes.io/docs/reference/scheduling/policies/#predicates)

- Hard constraints (Memory requirement, nodeSelector)

[Priorities](https://kubernetes.io/docs/reference/scheduling/policies/#priorities)

- Soft constraints (spreading)

 
[Score stage](https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/#scoring) -  select Nnde with the highest score among the feasible ones to run the Pod. 


### Assigning pods to nodes using

- nodeName
- NodeSelector
- NodeAffinity
- PodAffinity
- Taints
- Tolerations

### Use kubectl explain to gain insight

- kubectl explain pod.spec.nodeName
- kubectl explain pod.spec.nodeSelector
- kubectl explain pod.spec.tolerations
- kubectl explain node.spec.taints
- kubectl explain pod.spec.affinity

### Create kind cluster 

Using a kind cluster allows us to have a multi-node cluster locally gives us very good insight into how the scheduler is executed and configured.
 
`vi kind-config.yaml`

```
kind: Cluster
apiVersion: kind.sigs.k8s.io/v1alpha4
nodes:
- role: control-plane
- role: worker
- role: worker
- role: worker
```

```BASH
kind create cluster --name devcloud --config kind-config.yaml --image kindest/node:v1.18.2
```

### Review RBAC permissions for scheduler

The scheduler has to interact with the kube-apiserver and needs the least privileges. Review those permissions. Compare them to the admin account `kubectl auth can-i --list --kubeconfig=/etc/kubernetes/admin.conf` or the kubelet `kubectl auth can-i --list --kubeconfig=/etc/kubernetes/kubelet.conf`

```BASH
docker exec -it devcloud-control-plane bash
root@devcloud-control-plane:/# kubectl auth can-i --list --kubeconfig=/etc/kubernetes/scheduler.conf

```

```
Resources                                       Non-Resource URLs   Resource Names     Verbs
events                                          []                  []                 [create patch update]
events.events.k8s.io                            []                  []                 [create patch update]
bindings                                        []                  []                 [create]
endpoints                                       []                  []                 [create]
pods/binding                                    []                  []                 [create]
tokenreviews.authentication.k8s.io              []                  []                 [create]
selfsubjectaccessreviews.authorization.k8s.io   []                  []                 [create]
selfsubjectrulesreviews.authorization.k8s.io    []                  []                 [create]
subjectaccessreviews.authorization.k8s.io       []                  []                 [create]
leases.coordination.k8s.io                      []                  []                 [create]
pods                                            []                  []                 [delete get list watch]
persistentvolumeclaims                          []                  []                 [get list watch get list patch update watch]
persistentvolumes                               []                  []                 [get list watch get list patch update watch]
nodes                                           []                  []                 [get list watch]
replicationcontrollers                          []                  []                 [get list watch]
services                                        []                  []                 [get list watch]
replicasets.apps                                []                  []                 [get list watch]
statefulsets.apps                               []                  []                 [get list watch]
replicasets.extensions                          []                  []                 [get list watch]
poddisruptionbudgets.policy                     []                  []                 [get list watch]
csinodes.storage.k8s.io                         []                  []                 [get list watch]
storageclasses.storage.k8s.io                   []                  []                 [get list watch]
endpoints                                       []                  [kube-scheduler]   [get update]
leases.coordination.k8s.io                      []                  [kube-scheduler]   [get update]
                                                [/api/*]            []                 [get]
                                                [/api]              []                 [get]
                                                [/apis/*]           []                 [get]
                                                [/apis]             []                 [get]
                                                [/healthz]          []                 [get]
                                                [/healthz]          []                 [get]
                                                [/livez]            []                 [get]
                                                [/livez]            []                 [get]
                                                [/openapi/*]        []                 [get]
                                                [/openapi]          []                 [get]
                                                [/readyz]           []                 [get]
                                                [/readyz]           []                 [get]
                                                [/version/]         []                 [get]
                                                [/version/]         []                 [get]
                                                [/version]          []                 [get]
                                                [/version]          []                 [get]
pods/status                                     []                  []                 [patch update]

```
### View the pod manifest for kube-scheduler

```
root@devcloud-control-plane:/# cat /etc/kubernetes/manifests/kube-scheduler.yaml

```

### Review the client certificate data

```
root@devcloud-control-plane:/# kubectl config view --flatten --kubeconfig=/etc/kubernetes/scheduler.conf  | grep "client-certificate-data:" | awk '{print $2}' | base64 -d | openssl x509 -text
```

### Add taint and label to nodes

```BASH
kubectl taint nodes devcloud-worker3 key=runner:NoSchedule
kubectl label no devcloud-worker disk=ssd
```

### Pod with tolerations

```
apiVersion: v1
kind: Pod
metadata:
  name: ci-runner
  labels:
    env: ci
    app: ci-runner
spec:
  containers:
  - name: ci-runner
    image: k8s.gcr.io/pause
    imagePullPolicy: IfNotPresent
tolerations:
- key: "key"
  operator: "Equal"
  value: "runner"
  effect: "NoSchedule"
```

### Pod using nodeSelector

```
apiVersion: v1
kind: Pod
metadata:
  name: redis
  labels:
    env: uat
    app: redis
spec:
  nodeSelector:
    disk: ssd
  containers:
  - name: cache
    image: redis
    imagePullPolicy: IfNotPresent
```

### Pod using nodeName

```
apiVersion: v1
kind: Pod
metadata:
  name: webapp
  labels:
    env: uat
    app: webapp
spec:
  nodeName: devcloud-worker2
  containers:
  - name: webapp
    image: nginx
    imagePullPolicy: IfNotPresent
```    

### Pod with AntiAffinity that prevents pod from co-locating on a single node

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis-cache
spec:
  selector:
    matchLabels:
      app: store
  replicas: 3
  template:
    metadata:
      labels:
        app: store
    spec:
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
          - labelSelector:
              matchExpressions:
              - key: app
                operator: In
                values:
                - store
            topologyKey: "kubernetes.io/hostname"
      containers:
      - name: redis-server
        image: redis:3.2-alpine
```

### pod with AntiAffinity and Affinity. Prevent web-store pod from being co-located on a single node. Ensure web-store is co-located on a single node with redis

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-server
spec:
  selector:
    matchLabels:
      app: web-store
  replicas: 3
  template:
    metadata:
      labels:
        app: web-store
    spec:
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
          - labelSelector:
              matchExpressions:
              - key: app
                operator: In
                values:
                - web-store
            topologyKey: "kubernetes.io/hostname"
        podAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
          - labelSelector:
              matchExpressions:
              - key: app
                operator: In
                values:
                - store
            topologyKey: "kubernetes.io/hostname"
      containers:
      - name: web-app
        image: nginx:1.16-alpine        
```

    
### Additional Sources

- https://kubernetes.io/docs/concepts/scheduling-eviction/kube-scheduler/
- https://kubernetes.io/docs/concepts/scheduling-eviction/scheduler-perf-tuning/
- https://www.alibabacloud.com/blog/a-brief-analysis-on-the-implementation-of-the-kubernetes-scheduler_595083
- https://github.com/jamiehannaford/what-happens-when-k8s
- https://jvns.ca/blog/2017/07/27/how-does-the-kubernetes-scheduler-work/
- https://www.magalix.com/blog/kubernetes-scheduler-101
- https://medium.com/@dominik.tornow/the-kubernetes-scheduler-cd429abac02f
- https://www.youtube.com/watch?v=rDCWxkvPlAw
- https://www.youtube.com/watch?v=eDkE4WNWKUc

