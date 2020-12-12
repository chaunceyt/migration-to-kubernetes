# Kubernetes deployment notes

The `kind: Deployment` is an abstraction over `kind: ReplicaSet` 

- [Deployments](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/)
- [kubernetes/kubernetes/pkg/controller/deployment](https://github.com/kubernetes/kubernetes/blob/master/pkg/controller/deployment/deployment_controller.go)

```
kubectl explain deployment
kubectl explain deployment.spec.strategy

```

```
# Create a namespace
kubectl create ns production

# Deploy nginx 1.18 to the production namespace
kubectl create deployment webapp --image=nginx:1.18-alpine -n production

# List the objects created
kubectl get deploy,rs,po -l app=webapp -n production

# Describe the replicaset and pod noting what controlls them
kubectl describe rs -l app=webapp -n production | grep "Controlled By:"
kubectl describe po -l app=webapp -n production | grep "Controlled By:"

# Scale the deploy and note the strategy
kubectl scale deploy webapp --replicas=3 -n production

# in a different terminal watch
kubectl get po -n production

# Update the deployment and not the rollout strategy
kubectl set image deploy webapp nginx=nginx:1.19-alpine -n production

# Review the rollout history for the deployment
kubectl rollout history deploy webapp -n production
kubectl rollout history deploy webapp -n production --revision=1

# Undo a deployment
kubectl rollout undo deploy webapp --to-revision=1 -n production
```


## Verbose output
```
kubectl create ns stage
kubectl create deployment webapp2 --image=nginx:1.18-alpine -n stage -v10
```

```
I1206 10:36:42.956801   48740 loader.go:375] Config loaded from file:  /Users/cthorn/.kube/config
I1206 10:36:42.957704   48740 round_trippers.go:423] curl -k -v -XGET  -H "Accept: application/json, */*" -H "User-Agent: kubectl/v1.18.6 (darwin/amd64) kubernetes/dff82dc" 'https://127.0.0.1:50682/apis/apps/v1?timeout=32s'
I1206 10:36:42.972574   48740 round_trippers.go:443] GET https://127.0.0.1:50682/apis/apps/v1?timeout=32s 200 OK in 14 milliseconds
I1206 10:36:42.972588   48740 round_trippers.go:449] Response Headers:
I1206 10:36:42.972592   48740 round_trippers.go:452]     Date: Sun, 06 Dec 2020 15:36:42 GMT
I1206 10:36:42.972595   48740 round_trippers.go:452]     Cache-Control: no-cache, private
I1206 10:36:42.972597   48740 round_trippers.go:452]     Content-Type: application/json
I1206 10:36:42.972600   48740 round_trippers.go:452]     Content-Length: 2196
I1206 10:36:42.972768   48740 request.go:1068] Response Body: {"kind":"APIResourceList","apiVersion":"v1","groupVersion":"apps/v1","resources":[{"name":"controllerrevisions","singularName":"","namespaced":true,"kind":"ControllerRevision","verbs":["create","delete","deletecollection","get","list","patch","update","watch"],"storageVersionHash":"85nkx63pcBU="},{"name":"daemonsets","singularName":"","namespaced":true,"kind":"DaemonSet","verbs":["create","delete","deletecollection","get","list","patch","update","watch"],"shortNames":["ds"],"categories":["all"],"storageVersionHash":"dd7pWHUlMKQ="},{"name":"daemonsets/status","singularName":"","namespaced":true,"kind":"DaemonSet","verbs":["get","patch","update"]},{"name":"deployments","singularName":"","namespaced":true,"kind":"Deployment","verbs":["create","delete","deletecollection","get","list","patch","update","watch"],"shortNames":["deploy"],"categories":["all"],"storageVersionHash":"8aSe+NMegvE="},{"name":"deployments/scale","singularName":"","namespaced":true,"group":"autoscaling","version":"v1","kind":"Scale","verbs":["get","patch","update"]},{"name":"deployments/status","singularName":"","namespaced":true,"kind":"Deployment","verbs":["get","patch","update"]},{"name":"replicasets","singularName":"","namespaced":true,"kind":"ReplicaSet","verbs":["create","delete","deletecollection","get","list","patch","update","watch"],"shortNames":["rs"],"categories":["all"],"storageVersionHash":"P1RzHs8/mWQ="},{"name":"replicasets/scale","singularName":"","namespaced":true,"group":"autoscaling","version":"v1","kind":"Scale","verbs":["get","patch","update"]},{"name":"replicasets/status","singularName":"","namespaced":true,"kind":"ReplicaSet","verbs":["get","patch","update"]},{"name":"statefulsets","singularName":"","namespaced":true,"kind":"StatefulSet","verbs":["create","delete","deletecollection","get","list","patch","update","watch"],"shortNames":["sts"],"categories":["all"],"storageVersionHash":"H+vl74LkKdo="},{"name":"statefulsets/scale","singularName":"","namespaced":true,"group":"autoscaling","version":"v1","kind":"Scale","verbs":["get","patch","update"]},{"name":"statefulsets/status","singularName":"","namespaced":true,"kind":"StatefulSet","verbs":["get","patch","update"]}]}
I1206 10:36:42.974148   48740 cached_discovery.go:114] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/servergroups.json
I1206 10:36:42.974294   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/discovery.k8s.io/v1beta1/serverresources.json
I1206 10:36:42.974335   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/authentication.k8s.io/v1beta1/serverresources.json
I1206 10:36:42.974383   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/rbac.authorization.k8s.io/v1beta1/serverresources.json
I1206 10:36:42.974456   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/rbac.authorization.k8s.io/v1/serverresources.json
I1206 10:36:42.974461   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/storage.k8s.io/v1/serverresources.json
I1206 10:36:42.974459   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/apps/v1/serverresources.json
I1206 10:36:42.974462   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/extensions/v1beta1/serverresources.json
I1206 10:36:42.974536   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/events.k8s.io/v1/serverresources.json
I1206 10:36:42.974591   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/events.k8s.io/v1beta1/serverresources.json
I1206 10:36:42.974581   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/authorization.k8s.io/v1/serverresources.json
I1206 10:36:42.974632   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/v1/serverresources.json
I1206 10:36:42.974647   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/authentication.k8s.io/v1/serverresources.json
I1206 10:36:42.974842   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/admissionregistration.k8s.io/v1/serverresources.json
I1206 10:36:42.974850   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/autoscaling/v2beta2/serverresources.json
I1206 10:36:42.974704   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/apiregistration.k8s.io/v1/serverresources.json
I1206 10:36:42.974722   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/storage.k8s.io/v1beta1/serverresources.json
I1206 10:36:42.974724   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/autoscaling/v1/serverresources.json
I1206 10:36:42.975007   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/apiextensions.k8s.io/v1beta1/serverresources.json
I1206 10:36:42.974761   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/apiregistration.k8s.io/v1beta1/serverresources.json
I1206 10:36:42.974786   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/autoscaling/v2beta1/serverresources.json
I1206 10:36:42.974837   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/scheduling.k8s.io/v1/serverresources.json
I1206 10:36:42.974658   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/authorization.k8s.io/v1beta1/serverresources.json
I1206 10:36:42.974896   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/admissionregistration.k8s.io/v1beta1/serverresources.json
I1206 10:36:42.974906   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/batch/v1/serverresources.json
I1206 10:36:42.974949   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/apiextensions.k8s.io/v1/serverresources.json
I1206 10:36:42.974964   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/certificates.k8s.io/v1beta1/serverresources.json
I1206 10:36:42.974967   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/networking.k8s.io/v1/serverresources.json
I1206 10:36:42.975028   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/networking.k8s.io/v1beta1/serverresources.json
I1206 10:36:42.975042   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/batch/v1beta1/serverresources.json
I1206 10:36:42.975078   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/scheduling.k8s.io/v1beta1/serverresources.json
I1206 10:36:42.975086   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/policy/v1beta1/serverresources.json
I1206 10:36:42.975111   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/certificates.k8s.io/v1/serverresources.json
I1206 10:36:42.975146   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/coordination.k8s.io/v1beta1/serverresources.json
I1206 10:36:42.975167   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/coordination.k8s.io/v1/serverresources.json
I1206 10:36:42.975169   48740 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/node.k8s.io/v1beta1/serverresources.json
I1206 10:36:42.978428   48740 request.go:1068] Request Body: {"apiVersion":"apps/v1","kind":"Deployment","metadata":{"creationTimestamp":null,"labels":{"app":"webapp2"},"name":"webapp2"},"spec":{"replicas":1,"selector":{"matchLabels":{"app":"webapp2"}},"strategy":{},"template":{"metadata":{"creationTimestamp":null,"labels":{"app":"webapp2"}},"spec":{"containers":[{"image":"nginx:1.18-alpine","name":"nginx","resources":{}}]}}},"status":{}}
I1206 10:36:42.978479   48740 round_trippers.go:423] curl -k -v -XPOST  -H "Accept: application/json" -H "User-Agent: kubectl/v1.18.6 (darwin/amd64) kubernetes/dff82dc" 'https://127.0.0.1:50682/apis/apps/v1/namespaces/stage/deployments'
I1206 10:36:42.988584   48740 round_trippers.go:443] POST https://127.0.0.1:50682/apis/apps/v1/namespaces/stage/deployments 201 Created in 10 milliseconds
I1206 10:36:42.988598   48740 round_trippers.go:449] Response Headers:
I1206 10:36:42.988602   48740 round_trippers.go:452]     Cache-Control: no-cache, private
I1206 10:36:42.988604   48740 round_trippers.go:452]     Content-Type: application/json
I1206 10:36:42.988607   48740 round_trippers.go:452]     Content-Length: 1761
I1206 10:36:42.988609   48740 round_trippers.go:452]     Date: Sun, 06 Dec 2020 15:36:43 GMT
I1206 10:36:42.988979   48740 request.go:1068] Response Body: {"kind":"Deployment","apiVersion":"apps/v1","metadata":{"name":"webapp2","namespace":"stage","selfLink":"/apis/apps/v1/namespaces/stage/deployments/webapp2","uid":"b45c3f85-7ae7-4331-b3d8-47361c18aa22","resourceVersion":"563522","generation":1,"creationTimestamp":"2020-12-06T15:36:43Z","labels":{"app":"webapp2"},"managedFields":[{"manager":"kubectl","operation":"Update","apiVersion":"apps/v1","time":"2020-12-06T15:36:43Z","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:labels":{".":{},"f:app":{}}},"f:spec":{"f:progressDeadlineSeconds":{},"f:replicas":{},"f:revisionHistoryLimit":{},"f:selector":{"f:matchLabels":{".":{},"f:app":{}}},"f:strategy":{"f:rollingUpdate":{".":{},"f:maxSurge":{},"f:maxUnavailable":{}},"f:type":{}},"f:template":{"f:metadata":{"f:labels":{".":{},"f:app":{}}},"f:spec":{"f:containers":{"k:{\"name\":\"nginx\"}":{".":{},"f:image":{},"f:imagePullPolicy":{},"f:name":{},"f:resources":{},"f:terminationMessagePath":{},"f:terminationMessagePolicy":{}}},"f:dnsPolicy":{},"f:restartPolicy":{},"f:schedulerName":{},"f:securityContext":{},"f:terminationGracePeriodSeconds":{}}}}}}]},"spec":{"replicas":1,"selector":{"matchLabels":{"app":"webapp2"}},"template":{"metadata":{"creationTimestamp":null,"labels":{"app":"webapp2"}},"spec":{"containers":[{"name":"nginx","image":"nginx:1.18-alpine","resources":{},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"IfNotPresent"}],"restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"ClusterFirst","securityContext":{},"schedulerName":"default-scheduler"}},"strategy":{"type":"RollingUpdate","rollingUpdate":{"maxUnavailable":"25%","maxSurge":"25%"}},"revisionHistoryLimit":10,"progressDeadlineSeconds":600},"status":{}}
deployment.apps/webapp2 created
```

`kubectl set image deploy webapp2 nginx=nginx:1.19-alpine -n stage -v10`

```
I1206 10:41:24.194800   48758 loader.go:375] Config loaded from file:  /Users/cthorn/.kube/config
I1206 10:41:24.196108   48758 cached_discovery.go:114] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/servergroups.json
I1206 10:41:24.196567   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/discovery.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.196575   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/apiregistration.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.196588   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/events.k8s.io/v1/serverresources.json
I1206 10:41:24.196599   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/extensions/v1beta1/serverresources.json
I1206 10:41:24.196602   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/apps/v1/serverresources.json
I1206 10:41:24.196639   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/events.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.196667   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/apiregistration.k8s.io/v1/serverresources.json
I1206 10:41:24.196686   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/policy/v1beta1/serverresources.json
I1206 10:41:24.196735   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/autoscaling/v2beta2/serverresources.json
I1206 10:41:24.196753   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/v1/serverresources.json
I1206 10:41:24.196785   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/authentication.k8s.io/v1/serverresources.json
I1206 10:41:24.196796   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/batch/v1/serverresources.json
I1206 10:41:24.196837   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/authorization.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.196859   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/autoscaling/v1/serverresources.json
I1206 10:41:24.196864   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/batch/v1beta1/serverresources.json
I1206 10:41:24.197273   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/coordination.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.196895   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/authentication.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.196917   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/autoscaling/v2beta1/serverresources.json
I1206 10:41:24.196934   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/certificates.k8s.io/v1/serverresources.json
I1206 10:41:24.196962   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/authorization.k8s.io/v1/serverresources.json
I1206 10:41:24.196973   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/apiextensions.k8s.io/v1/serverresources.json
I1206 10:41:24.197004   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/certificates.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.197021   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/scheduling.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.197053   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/rbac.authorization.k8s.io/v1/serverresources.json
I1206 10:41:24.197073   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/networking.k8s.io/v1/serverresources.json
I1206 10:41:24.197076   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/apiextensions.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.197122   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/rbac.authorization.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.197127   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/networking.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.197135   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/scheduling.k8s.io/v1/serverresources.json
I1206 10:41:24.197192   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/storage.k8s.io/v1/serverresources.json
I1206 10:41:24.197211   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/storage.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.197215   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/admissionregistration.k8s.io/v1/serverresources.json
I1206 10:41:24.197249   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/admissionregistration.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.197281   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/coordination.k8s.io/v1/serverresources.json
I1206 10:41:24.197305   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/node.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.197750   48758 cached_discovery.go:114] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/servergroups.json
I1206 10:41:24.197867   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/events.k8s.io/v1/serverresources.json
I1206 10:41:24.197870   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/certificates.k8s.io/v1/serverresources.json
I1206 10:41:24.197890   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/discovery.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.197919   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/authorization.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.197923   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/autoscaling/v1/serverresources.json
I1206 10:41:24.197945   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/certificates.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.197974   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/events.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.197978   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/autoscaling/v2beta1/serverresources.json
I1206 10:41:24.198167   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/batch/v1beta1/serverresources.json
I1206 10:41:24.198028   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/v1/serverresources.json
I1206 10:41:24.198038   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/authentication.k8s.io/v1/serverresources.json
I1206 10:41:24.198040   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/autoscaling/v2beta2/serverresources.json
I1206 10:41:24.198295   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/apps/v1/serverresources.json
I1206 10:41:24.198085   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/apiregistration.k8s.io/v1/serverresources.json
I1206 10:41:24.198099   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/batch/v1/serverresources.json
I1206 10:41:24.198135   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/authentication.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.198145   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/apiregistration.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.198201   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/extensions/v1beta1/serverresources.json
I1206 10:41:24.198212   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/authorization.k8s.io/v1/serverresources.json
I1206 10:41:24.198271   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/networking.k8s.io/v1/serverresources.json
I1206 10:41:24.198280   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/admissionregistration.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.198334   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/apiextensions.k8s.io/v1/serverresources.json
I1206 10:41:24.198615   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/node.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.198342   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/networking.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.198391   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/apiextensions.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.198410   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/storage.k8s.io/v1/serverresources.json
I1206 10:41:24.198383   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/rbac.authorization.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.198450   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/scheduling.k8s.io/v1/serverresources.json
I1206 10:41:24.198473   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/storage.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.198494   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/policy/v1beta1/serverresources.json
I1206 10:41:24.198499   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/scheduling.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.198534   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/admissionregistration.k8s.io/v1/serverresources.json
I1206 10:41:24.198562   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/coordination.k8s.io/v1/serverresources.json
I1206 10:41:24.198566   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/rbac.authorization.k8s.io/v1/serverresources.json
I1206 10:41:24.198590   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/coordination.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.198844   48758 cached_discovery.go:114] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/servergroups.json
I1206 10:41:24.198915   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/discovery.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.198976   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/certificates.k8s.io/v1/serverresources.json
I1206 10:41:24.198977   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/certificates.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.199014   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/admissionregistration.k8s.io/v1/serverresources.json
I1206 10:41:24.199028   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/admissionregistration.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.199090   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/networking.k8s.io/v1/serverresources.json
I1206 10:41:24.199151   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/v1/serverresources.json
I1206 10:41:24.199155   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/networking.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.199180   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/rbac.authorization.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.199212   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/apiextensions.k8s.io/v1/serverresources.json
I1206 10:41:24.199218   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/apiregistration.k8s.io/v1/serverresources.json
I1206 10:41:24.199228   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/storage.k8s.io/v1/serverresources.json
I1206 10:41:24.199259   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/policy/v1beta1/serverresources.json
I1206 10:41:24.199281   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/apiextensions.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.199289   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/storage.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.199325   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/rbac.authorization.k8s.io/v1/serverresources.json
I1206 10:41:24.199345   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/coordination.k8s.io/v1/serverresources.json
I1206 10:41:24.199382   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/scheduling.k8s.io/v1/serverresources.json
I1206 10:41:24.199387   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/scheduling.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.199437   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/coordination.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.199452   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/authentication.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.199481   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/apps/v1/serverresources.json
I1206 10:41:24.199494   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/node.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.199510   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/events.k8s.io/v1/serverresources.json
I1206 10:41:24.199546   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/apiregistration.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.199567   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/autoscaling/v2beta1/serverresources.json
I1206 10:41:24.199570   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/events.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.199603   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/extensions/v1beta1/serverresources.json
I1206 10:41:24.199627   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/authorization.k8s.io/v1/serverresources.json
I1206 10:41:24.199635   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/authentication.k8s.io/v1/serverresources.json
I1206 10:41:24.199666   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/batch/v1/serverresources.json
I1206 10:41:24.199742   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/authorization.k8s.io/v1beta1/serverresources.json
I1206 10:41:24.199743   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/autoscaling/v1/serverresources.json
I1206 10:41:24.199744   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/autoscaling/v2beta2/serverresources.json
I1206 10:41:24.199829   48758 cached_discovery.go:71] returning cached discovery info from /Users/cthorn/.kube/cache/discovery/127.0.0.1_50682/batch/v1beta1/serverresources.json
I1206 10:41:24.201904   48758 round_trippers.go:423] curl -k -v -XGET  -H "Accept: application/json, */*" -H "User-Agent: kubectl/v1.18.6 (darwin/amd64) kubernetes/dff82dc" 'https://127.0.0.1:50682/apis/apps/v1/namespaces/stage/deployments/webapp2'
I1206 10:41:24.219117   48758 round_trippers.go:443] GET https://127.0.0.1:50682/apis/apps/v1/namespaces/stage/deployments/webapp2 200 OK in 17 milliseconds
I1206 10:41:24.219143   48758 round_trippers.go:449] Response Headers:
I1206 10:41:24.219146   48758 round_trippers.go:452]     Cache-Control: no-cache, private
I1206 10:41:24.219150   48758 round_trippers.go:452]     Content-Type: application/json
I1206 10:41:24.219167   48758 round_trippers.go:452]     Content-Length: 3026
I1206 10:41:24.219169   48758 round_trippers.go:452]     Date: Sun, 06 Dec 2020 15:41:24 GMT
I1206 10:41:24.219230   48758 request.go:1068] Response Body: {"kind":"Deployment","apiVersion":"apps/v1","metadata":{"name":"webapp2","namespace":"stage","selfLink":"/apis/apps/v1/namespaces/stage/deployments/webapp2","uid":"b45c3f85-7ae7-4331-b3d8-47361c18aa22","resourceVersion":"563544","generation":1,"creationTimestamp":"2020-12-06T15:36:43Z","labels":{"app":"webapp2"},"annotations":{"deployment.kubernetes.io/revision":"1"},"managedFields":[{"manager":"kubectl","operation":"Update","apiVersion":"apps/v1","time":"2020-12-06T15:36:43Z","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:labels":{".":{},"f:app":{}}},"f:spec":{"f:progressDeadlineSeconds":{},"f:replicas":{},"f:revisionHistoryLimit":{},"f:selector":{"f:matchLabels":{".":{},"f:app":{}}},"f:strategy":{"f:rollingUpdate":{".":{},"f:maxSurge":{},"f:maxUnavailable":{}},"f:type":{}},"f:template":{"f:metadata":{"f:labels":{".":{},"f:app":{}}},"f:spec":{"f:containers":{"k:{\"name\":\"nginx\"}":{".":{},"f:image":{},"f:imagePullPolicy":{},"f:name":{},"f:resources":{},"f:terminationMessagePath":{},"f:terminationMessagePolicy":{}}},"f:dnsPolicy":{},"f:restartPolicy":{},"f:schedulerName":{},"f:securityContext":{},"f:terminationGracePeriodSeconds":{}}}}}},{"manager":"kube-controller-manager","operation":"Update","apiVersion":"apps/v1","time":"2020-12-06T15:36:44Z","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:deployment.kubernetes.io/revision":{}}},"f:status":{"f:availableReplicas":{},"f:conditions":{".":{},"k:{\"type\":\"Available\"}":{".":{},"f:lastTransitionTime":{},"f:lastUpdateTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}},"k:{\"type\":\"Progressing\"}":{".":{},"f:lastTransitionTime":{},"f:lastUpdateTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}}},"f:observedGeneration":{},"f:readyReplicas":{},"f:replicas":{},"f:updatedReplicas":{}}}}]},"spec":{"replicas":1,"selector":{"matchLabels":{"app":"webapp2"}},"template":{"metadata":{"creationTimestamp":null,"labels":{"app":"webapp2"}},"spec":{"containers":[{"name":"nginx","image":"nginx:1.18-alpine","resources":{},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"IfNotPresent"}],"restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"ClusterFirst","securityContext":{},"schedulerName":"default-scheduler"}},"strategy":{"type":"RollingUpdate","rollingUpdate":{"maxUnavailable":"25%","maxSurge":"25%"}},"revisionHistoryLimit":10,"progressDeadlineSeconds":600},"status":{"observedGeneration":1,"replicas":1,"updatedReplicas":1,"readyReplicas":1,"availableReplicas":1,"conditions":[{"type":"Available","status":"True","lastUpdateTime":"2020-12-06T15:36:44Z","lastTransitionTime":"2020-12-06T15:36:44Z","reason":"MinimumReplicasAvailable","message":"Deployment has minimum availability."},{"type":"Progressing","status":"True","lastUpdateTime":"2020-12-06T15:36:44Z","lastTransitionTime":"2020-12-06T15:36:43Z","reason":"NewReplicaSetAvailable","message":"ReplicaSet \"webapp2-88f64ffd9\" has successfully progressed."}]}}
I1206 10:41:24.224411   48758 request.go:1068] Request Body: {"spec":{"template":{"spec":{"$setElementOrder/containers":[{"name":"nginx"}],"containers":[{"image":"nginx:1.19-alpine","name":"nginx"}]}}}}
I1206 10:41:24.224455   48758 round_trippers.go:423] curl -k -v -XPATCH  -H "Accept: application/json, */*" -H "Content-Type: application/strategic-merge-patch+json" -H "User-Agent: kubectl/v1.18.6 (darwin/amd64) kubernetes/dff82dc" 'https://127.0.0.1:50682/apis/apps/v1/namespaces/stage/deployments/webapp2'
I1206 10:41:24.236315   48758 round_trippers.go:443] PATCH https://127.0.0.1:50682/apis/apps/v1/namespaces/stage/deployments/webapp2 200 OK in 11 milliseconds
I1206 10:41:24.236328   48758 round_trippers.go:449] Response Headers:
I1206 10:41:24.236332   48758 round_trippers.go:452]     Cache-Control: no-cache, private
I1206 10:41:24.236334   48758 round_trippers.go:452]     Content-Type: application/json
I1206 10:41:24.236337   48758 round_trippers.go:452]     Content-Length: 3026
I1206 10:41:24.236339   48758 round_trippers.go:452]     Date: Sun, 06 Dec 2020 15:41:24 GMT
I1206 10:41:24.236623   48758 request.go:1068] Response Body: {"kind":"Deployment","apiVersion":"apps/v1","metadata":{"name":"webapp2","namespace":"stage","selfLink":"/apis/apps/v1/namespaces/stage/deployments/webapp2","uid":"b45c3f85-7ae7-4331-b3d8-47361c18aa22","resourceVersion":"564357","generation":2,"creationTimestamp":"2020-12-06T15:36:43Z","labels":{"app":"webapp2"},"annotations":{"deployment.kubernetes.io/revision":"1"},"managedFields":[{"manager":"kube-controller-manager","operation":"Update","apiVersion":"apps/v1","time":"2020-12-06T15:36:44Z","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:deployment.kubernetes.io/revision":{}}},"f:status":{"f:availableReplicas":{},"f:conditions":{".":{},"k:{\"type\":\"Available\"}":{".":{},"f:lastTransitionTime":{},"f:lastUpdateTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}},"k:{\"type\":\"Progressing\"}":{".":{},"f:lastTransitionTime":{},"f:lastUpdateTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}}},"f:observedGeneration":{},"f:readyReplicas":{},"f:replicas":{},"f:updatedReplicas":{}}}},{"manager":"kubectl","operation":"Update","apiVersion":"apps/v1","time":"2020-12-06T15:41:24Z","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:labels":{".":{},"f:app":{}}},"f:spec":{"f:progressDeadlineSeconds":{},"f:replicas":{},"f:revisionHistoryLimit":{},"f:selector":{"f:matchLabels":{".":{},"f:app":{}}},"f:strategy":{"f:rollingUpdate":{".":{},"f:maxSurge":{},"f:maxUnavailable":{}},"f:type":{}},"f:template":{"f:metadata":{"f:labels":{".":{},"f:app":{}}},"f:spec":{"f:containers":{"k:{\"name\":\"nginx\"}":{".":{},"f:image":{},"f:imagePullPolicy":{},"f:name":{},"f:resources":{},"f:terminationMessagePath":{},"f:terminationMessagePolicy":{}}},"f:dnsPolicy":{},"f:restartPolicy":{},"f:schedulerName":{},"f:securityContext":{},"f:terminationGracePeriodSeconds":{}}}}}}]},"spec":{"replicas":1,"selector":{"matchLabels":{"app":"webapp2"}},"template":{"metadata":{"creationTimestamp":null,"labels":{"app":"webapp2"}},"spec":{"containers":[{"name":"nginx","image":"nginx:1.19-alpine","resources":{},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"IfNotPresent"}],"restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"ClusterFirst","securityContext":{},"schedulerName":"default-scheduler"}},"strategy":{"type":"RollingUpdate","rollingUpdate":{"maxUnavailable":"25%","maxSurge":"25%"}},"revisionHistoryLimit":10,"progressDeadlineSeconds":600},"status":{"observedGeneration":1,"replicas":1,"updatedReplicas":1,"readyReplicas":1,"availableReplicas":1,"conditions":[{"type":"Available","status":"True","lastUpdateTime":"2020-12-06T15:36:44Z","lastTransitionTime":"2020-12-06T15:36:44Z","reason":"MinimumReplicasAvailable","message":"Deployment has minimum availability."},{"type":"Progressing","status":"True","lastUpdateTime":"2020-12-06T15:36:44Z","lastTransitionTime":"2020-12-06T15:36:43Z","reason":"NewReplicaSetAvailable","message":"ReplicaSet \"webapp2-88f64ffd9\" has successfully progressed."}]}}
deployment.apps/webapp2 image updated
```
