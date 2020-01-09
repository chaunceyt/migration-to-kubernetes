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