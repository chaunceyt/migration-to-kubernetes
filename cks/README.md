# CKS notes

```
kubectl apply --recursive -f manifests
kubectl auth can-i --list --as system:serviceaccount:cks:cks-sa -n cks
```

```
# https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.19.md#v1193
curl -LO https://dl.k8s.io/v1.19.3/kubernetes-server-linux-amd64.tar.gz
sha512sum kubernetes-server-linux-amd64.tar.gz | awk '{print $1}'
# Confirm sha512 hash.

```

## RBAC 

- Principle of least privilege (POLP)
- `kubectl auth can-i --list`

## clustersigning request

```
openssl genrsa -out cthorn.key 4096
openssl req -new -key cthorn.key -out cthorn.csr

CSR_CERT=$(cat cthorn.csr | base64 -w 0)

cat > cthorn-csr.yaml

apiVersion: certificates.k8s.io/v1
kind: CertificateSigningRequest
metadata:
  name: cthorn
spec:
  groups:
  - system:authenticated
  request: $CSR_CERT
  signerName: kubernetes.io/kube-apiserver-client
  usages:
  - client auth

```

### Approve certificate signing request


```
kubectl certificate approve cthorn
```

## Get ca, cert and key


```
kubectl config view --raw | grep "certificate-authority-data:" | awk '{print $2}' | base64 -d > ca
kubectl config view --raw | grep "client-certificate-data: " | awk '{print $2}' | base64 -d > crt
kubectl config view --raw | grep "client-key-data:" | awk '{print $2}' | base64 -d > key

API_SERVER=(kubectl config view --raw | grep "server:" | awk '{print $2}')

curl https://${API_SERVER} --cacert ca --cert crt --key key
```



```
cd /etc/kubernetes/pki
openssl x509 -in apiserver.crt -text

```
### CIS Benchmarks

[GKE CIS results](https://cloud.google.com/kubernetes-engine/docs/concepts/cis-benchmarks)
[starboard](https://github.com/aquasecurity/starboard)


```
# master
sudo docker run --pid=host -v /etc:/etc:ro -v /var:/var:ro -t aquasec/kube-bench:latest master --version 1.19

# worker 
sudo docker run --pid=host -v /etc:/etc:ro -v /var:/var:ro -t aquasec/kube-bench:latest master --version 1.19
```

## Enable Audit

## Prepare for node upgrade

```
kubectl drain
kubectl uncordon
```

	