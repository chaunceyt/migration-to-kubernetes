# Notes on Kubernetes secrets

- [Design docs for secrets](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/auth/secrets.md)
- [Documentation](https://kubernetes.io/docs/concepts/configuration/secret/)
- [TGI Kubernetes 120: CSI and Secrets!](https://www.youtube.com/watch?v=IznsHhKL428)
- [TGI Kubernetes 113: Kubernetes Secrets Take 3](https://www.youtube.com/watch?v=an9D2FyFwR0) ([git repo](https://github.com/vmware-tanzu/tgik/tree/master/episodes/113))
- [Understand secrets management in Kubernetes](https://www.youtube.com/watch?v=KmhM33j5WYk)
- [Base64 is not encryption](https://www.youtube.com/watch?v=f4Ru6CPG1z4)
- [kubesec](https://github.com/shyiko/kubesec) Secure secret management for Kubernetes (with gpg, Google Cloud KMS and AWS KMS backends).
- [Kamus](https://github.com/Soluto/kamus) manages encrypted secrets than can be decrypted only by the application
- [sops](https://github.com/mozilla/sops) an editor of encrypted files that supports YAML, JSON, ENV, INI and BINARY formats and encrypts with AWS KMS, GCP KMS, Azure Key Vault and PGP

```
# Get secret
kubectl get secret <secret_name>

# Get secret in yaml format
kubectl get secret <secret_name> -o yaml

# Edit secret
kubectl edit secret <secret_name>

# Describe secret
kubectl describe secret <secret_name>

# Delete secret
kubectl delete secret <secret_name>

# Create secret(s)
kubectl create secret generic project-secret --dry-run=client \
	--from-literal=username=admin \
	--from-literal=password=4P@s5wOrd2N0

```

NOTE: By default secrets are stored in etcd in plaintext. 

To confirm use a kind cluster.

```
kind create cluster --name secrets --kubeconfig secrets-cluster
docker exec -it secrets-control-plane bash
cd /etc/kubernetes/manifests
curl -LO https://git.io/etcdclient.yaml
exit

# Confirm pod is running
kubectl get po etcdclient-secrets-control-plane -n kube-system --kubeconfig secrets-cluster

# Create a secret
kubectl create secret generic project-secret \
	--from-literal=username=admin \
	--from-literal=password=4P@s5wOrd2N0 \
	--kubeconfig secrets-cluster

# Get contents of secret object in the cluster
kubectl get secret --kubeconfig secrets-cluster project-secret -o yaml

# output...
apiVersion: v1
data:
  password: NFBAczV3T3JkMk4w <--- base64 encoded string
  username: YWRtaW4= <--- base64 encoded string
kind: Secret
metadata:
  creationTimestamp: "2020-09-20T15:17:24Z"
  managedFields:
  - apiVersion: v1
    fieldsType: FieldsV1
    fieldsV1:
      f:data:
        .: {}
        f:password: {}
        f:username: {}
      f:type: {}
    manager: kubectl
    operation: Update
    time: "2020-09-20T15:17:24Z"
  name: project-secret
  namespace: default
  resourceVersion: "1918"
  selfLink: /api/v1/namespaces/default/secrets/project-secret
  uid: 45a82982-a5b3-4477-9a2a-0226c27177e1
type: Opaque

# what does it look like in etcd? exec into the static pod we deployed above.
kubectl --kubeconfig secrets-cluster exec -it etcdclient-secrets-control-plane -n kube-system -- sh

ETCDCTL_API=3 etcdctl get /registry/secrets/default/project-secret  --prefix

# output...

2020-09-20 15:23:02.604097 W | pkg/flags: unrecognized environment variable ETCDCTL_CLUSTER=true
/registry/secrets/default/project-secret
k8s


v1Secret�
�
project-secretdefault"*$45a82982-a5b3-4477-9a2a-0226c27177e12���z�l
kubectlUpdatev���FieldsV1:A
?{"f:data":{".":{},"f:password":{},"f:username":{}},"f:type":{}}
password
        4P@s5wOrd2N0 <---- password
usernameadminOpaque" <---- username


```
## Encryption at rest

[Docs for encrypting secrets](https://kubernetes.io/docs/tasks/administer-cluster/encrypt-data/)

```
docker exec -it secrets-control-plane

head -c 32 /dev/urandom | base64

cat > /etc/kubernetes/pki/enc.conf
apiVersion: apiserver.config.k8s.io/v1
kind: EncryptionConfiguration
resources:
  - resources:
    - secrets
    providers:
    - aescbc:
        keys:
        - name: key1
          secret: SECRET_HASH
    - identity: {}
ctrl c

cd /kind/
cat kubeadm.conf
# Copy contents to another file.
apiServer:
...
  extraArgs:
    ...
    encryption-provider-config: /etc/kubernetes/pki/enc.conf

cat > /kind/kubeadm2.conf
ctrl v
enter
ctrl c

# Reconfigure the apiserver
kubeadm init phase control-plane apiserver --config /kind/kubeadm2.conf

# Encrypt all secrets cluster-wide
kubectl get secrets --all-namespaces -o json | kubectl replace -f -

kubectl --kubeconfig secrets-cluster exec -it etcdclient-secrets-control-plane -n kube-system -- sh

ETCDCTL_API=3 etcdctl get /registry/secrets/default/project-secret  --prefix

# content should be encrypted.

```

### Undo encryption at rest

```
cat > /etc/kubernetes/pki/enc.conf
apiVersion: apiserver.config.k8s.io/v1
kind: EncryptionConfiguration
resources:
  - resources:
    - secrets
    providers:
    - identity: {}
    - aescbc:
        keys:
        - name: key1
          secret: SECRET_HASH
ctrl c

# restart kube-apiserver
mv /etc/kubernetes/manifests/kube-apiserver.yaml /tmp/
mv /tmp/kube-apiserver.yaml /etc/kubernetes/manifests/
exit

kubectl get secrets --all-namespaces -o json | kubectl replace -f -

kubectl --kubeconfig secrets-cluster exec -it etcdclient-secrets-control-plane -n kube-system -- sh

ETCDCTL_API=3 etcdctl get /registry/secrets/default/project-secret  --prefix

# content should be encrypted.

```
## Managing secrets

### Mutating webhook
[Using Kubernetes Secrets in GitOps Workflows Securely](https://www.youtube.com/watch?v=-k6HEXaE75k) using Google cloud KMS with JWE mutating webhook. [slides](https://static.sched.com/hosted_files/kccnceu20/98/Using%20Kubernetes%20Secrets%20in%20GitOps%20Workflows%20Securely.pdf) 
[Mutating webhook git repo](https://github.com/immutableT/k8s-secrets-and-gitops)


### Using Sealed Secrets.

Sealed Secrets is different from encryption of secrets in the cluster at rest in etcd. Sealed secrets encrypt the secrets before they are submitted to cluster. Making it secure to store the secret in a git repo.

Sealed Secrets introduces a CRD and controller.

When the controller starts in the cluster,  a private and public keypair are created. As a result, one can encrypt (kubeseal) with the public key and decrypt (controller) with private key. The `kubeseal` client will not work without a cluster.

- ["Sealed Secrets" for Kubernetes](https://github.com/bitnami-labs/sealed-secrets)
- [Managing secrets deployment in Kubernetes using Sealed Secrets](https://aws.amazon.com/blogs/opensource/managing-secrets-deployment-in-kubernetes-using-sealed-secrets/)
- [TGI Kubernetes 132: Sealed Secrets!](https://www.youtube.com/watch?v=x-cDk8DIvwE) [hackmd](https://hackmd.io/q2lXm3UFTv26dSsvfG1VRA)
- [19th Mar 2020 Angus Lees Kubernetes Sealed Secrets #MeetupMadness](https://www.youtube.com/watch?v=kYn8WN5UDiM)

#### Install sealed secrets and "play" around with it.

```
helm fetch stable/sealed-secrets --untar
helm template --name sealed-secrets-controller --namespace kube-system sealed-secrets/ | kubectl -n kube-system apply -f -
wget https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.12.5/kubeseal-darwin-amd64
mv kubeseal-darwin-amd64 $HOME/bin/kubeseal
chmod +x $HOME/bin/kubeseal

kubectl create secret generic project-secret --dry-run --from-literal=username=admin --dry-run --from-literal=password=4P@s5wOrd2N0 -o yaml | kubeseal > project-sealedsecret.json


kubectl create -f project-sealedsecret.json
kubectl get secret project-secret -o yaml
```

[Re-encryption (advanced)](https://github.com/bitnami-labs/sealed-secrets#re-encryption-advanced)


[Early key renewal](https://github.com/bitnami-labs/sealed-secrets#early-key-renewal) `--key-cutoff-time` = `date -R`


Get public cert from cluster to create
```
kubeseal --fetch-cert > mycert.crt
kubeseal -v 10 --scope-wide --cert mycert.crt <mysecret.yaml > mysealedsecret.json
```


#### Multi-Cluster Sealed Secrets setup.

```
kind create cluster --name cluster-1 --config ~/development/kind-cluster.yaml --kubeconfig kubeconfig-1
kind create cluster --name cluster-2 --config ~/development/kind-cluster.yaml --kubeconfig kubeconfig-2

helm template --name sealed-secrets-controller --namespace kube-system sealed-secrets/ | kubectl --kubeconfig kubeconfig-1 -n kube-system apply -f -

kubectl --kubeconfig kubeconfig-1 create secret generic project-secret --dry-run=client --from-literal=username=admin --dry-run=client --from-literal=password=4P@s5wOrd2N0 -o json  > project-secret.json

kubeseal --kubeconfig kubeconfig-1 < project-secret.json > project-sealed-secret.json
kubectl --kubeconfig kubeconfig-1 apply -f project-sealed-secret.json

helm template --name sealed-secrets-controller --namespace kube-system sealed-secrets/ | kubectl --kubeconfig kubeconfig-2 -n kube-system apply -f -

kubectl --kubeconfig kubeconfig-1 -n kube-system get secret sealed-secrets-XXXXXX -o json > sealing-key-cluster-1.json

kubectl apply -f project-sealed-secret.json --kubeconfig kubeconfig-2

# Note due to not having the signing key the secret isn't decrypted. As a result not created.
kubectl get secrets --kubeconfig kubeconfig-2

# Install they signing key for this sealedsecret
kubectl apply -f sealing-key-cluster-1.json --kubeconfig kubeconfig-2

# Restart the pod to use the second sealing key
kubectl --kubeconfig kubeconfig-2 delete  po sealed-secrets-controller-XXXX-XXXX  -n kube-system


# The project-secret should have been created in the namespace now.
kubectl get secrets  --kubeconfig kubeconfig-2
```

[crypto details](https://github.com/bitnami-labs/sealed-secrets/blob/master/docs/crypto.md)


Kind cluster `kind-cluster.yaml`

```
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
```


#### Important Note(s)

Every 30 days the sealing key is renewed. So understand the [secret rotation](https://github.com/bitnami-labs/sealed-secrets#secret-rotation) process
