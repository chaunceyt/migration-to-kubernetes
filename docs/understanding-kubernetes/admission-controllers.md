# Admission controllers notes



- [Admission Webhooks: Configuration and Debugging Best Practices - Haowei Cai, Google](https://www.youtube.com/watch?v=r_v07P8Go6w)
- [TGI Kubernetes 112: Deep dive into Admission Controllers](https://www.youtube.com/watch?v=fEvOzL_eosg)
- [Using Admission Controllers](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/)
- [A Guide to Kubernetes Admission Controllers](https://kubernetes.io/blog/2019/03/21/a-guide-to-kubernetes-admission-controllers/)
- [Dynamic Admission Control](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/)
- [Customizing and Extending the Kubernetes API with Admission Controllers](https://www.youtube.com/watch?v=P7QAfjdbogY)

Two types [docs](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/)

- Validating 
- Mutating

### Validating

Policy - OPA + Gatekeeper

- [TGI Kubernetes 119: Gatekeeper and OPA](https://www.youtube.com/watch?v=ZJgaGJm9NJE)

### Mutating

Injestion of sidcars, VPA

## Admission plugins

- [NamespaceLifecycle](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#namespacelifecycle)
- [LimitRanger](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#limitranger)
- [ServiceAccount](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#serviceaccount)
- [TaintNodesByCondition](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#taintnodesbycondition)
- [Priority](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#priority)
- [DefaultTolerationSeconds](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#defaulttolerationseconds)
- [DefaultStorageClass](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#defaultstorageclass)
- [PersistentVolumeClaimResize](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#persistentvolumeclaimresize)
- [MutatingAdmissionWebhook](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#mutatingadmissionwebhook)
- [ValidatingAdmissionWebhook](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#validatingadmissionwebhook)
- [ResourceQuota](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#resourcequota)


## Using Kind as the cluster environment.

```
---
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
kubeadmConfigPatches:
- |
  apiVersion: kubeadm.k8s.io/v1beta2
  kind: ClusterConfiguration
  metadata:
    name: config
  apiServer:
    extraArgs:
      "enable-admission-plugins": "NamespaceLifecycle,LimitRanger,ServiceAccount,TaintNodesByCondition,Priority,DefaultTolerationSeconds,DefaultStorageClass,PersistentVolumeClaimResize,MutatingAdmissionWebhook,ValidatingAdmissionWebhook,ResourceQuota"
nodes:
- role: control-plane
- role: worker
- role: worker
