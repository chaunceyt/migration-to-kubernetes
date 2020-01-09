# RBAC-Manager

## Create ClusterRole

```
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: web-developer
rules:
- apiGroups: [""]
  resources: ["pods", "pods/log"]
  verbs: ["get", "list"]
- apiGroups: [""]
  resources: ["pods/exec", "pods/portforward"]
  verbs: ["create"]
```

## Install

```
git clone https://github.com/reactiveops/rbac-manager.git
cd rbac-manager/
kubectl apply -f deploy/
kubectl get all -n rbac-manager
```

## Create RBACDefinition

```
apiVersion: rbacmanager.reactiveops.io/v1beta1
kind: RBACDefinition
metadata:
  name: web-developer-access
rbacBindings:
  - name: web-developer-rb
    subjects:
      - kind: Group
        name: developers@domain.tld
    roleBindings:
      - clusterRole: web-developer
        namespaceSelector:
          matchLabels:
            purpose: web-project
```

Now when we create a namespace with the label `purpose: web-project` The rolebinding responsible for group permission is applied to namespace by RBAC-Manager controller.

We found a what seems to be a better solution [here](https://raw.githubusercontent.com/kubernetes/k8s.io/master/infra/gcp/namespaces/namespace-user-role.yml).

```
---
kind: Role
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: namespace-user
  namespace: {{namespace}}
rules:
  - apiGroups: [""]
    resources: ["configmaps", "endpoints", "persistentvolumeclaims", "pods", "resourcequotas", "services"]
    verbs: ["*"]
  - apiGroups: ["certmanager.k8s.io"]
    resources: ["certificates"]
    verbs: ["*"]
  - apiGroups: ["coordination.k8s.io"]
    resources: ["leases"]
    verbs: ["*"]
  - apiGroups: ["batch"]
    resources: ["cronjobs", "jobs"]
    verbs: ["*"]
  - apiGroups: ["autoscaling"]
    resources: ["horizontalpodautoscalers"]
    verbs: ["*"]
  - apiGroups: ["apps"]
    resources: ["deployments"]
    verbs: ["*"]
  - apiGroups: ["extensions"]
    resources: ["deployments", "ingresses", "networkpolicies"]
    verbs: ["*"]
  - apiGroups: ["networking.k8s.io"]
    resources: ["networkpolicies"]
    verbs: ["*"]
  - apiGroups: [""]
    resources: ["secrets"]
    verbs: ["list"]
  - apiGroups: ["storage.k8s.io"]
    resources: ["storageclasses"]
    verbs: ["list"]
  - apiGroups: ["scheduling.k8s.io"]
    resources: ["priorityclasses"]
    verbs: ["list"]
  - apiGroups: ["rbac.authorization.k8s.io"]
    resources: ["clusterrolebindings", "clusterroles", "rolebindings", "roles"]
    verbs: ["get", "list"]
  - apiGroups: [""]
    resources: ["events"]
    verbs: ["get", "list"]
 ```