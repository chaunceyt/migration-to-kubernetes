# Using Kind clusters

## What is Kind?

Kind is a tool that allows one to run Kubernetes inside docker containers. It affords on the ability to create local clusters fast and with ease.

- [Kind official documentation](https://kind.sigs.k8s.io/)
- [Kind GitHub repository](https://github.com/kubernetes-sigs/kind)


# Create Kind cluster config

We have a number of options here. Use the default CNI or use Calico's

## Deploy using the kindnet CNI (default)

Create a manifest file named `kind-config.yaml` containing

```
kind: Cluster
apiVersion: kind.sigs.k8s.io/v1alpha3
nodes:
- role: control-plane
- role: worker
- role: worker
- role: worker
```

Create the cluster.

```
kind create cluster --config kind-config.yaml
```

## Deploy using the calico CNI

Create a manifest file named `kind-config-calico.yaml` containing

```
kind: Cluster
apiVersion: kind.sigs.k8s.io/v1alpha3
nodes:
- role: control-plane
- role: worker
- role: worker
- role: worker
networking:
  disableDefaultCNI: true # disable kindnet
  podSubnet: 192.168.0.0/16 # default subnet for Calico
```

Create the cluster.

```
kind create cluster --config kind-calico.yaml

# Install Calico 
kubectl apply -f https://docs.projectcalico.org/v3.8/manifests/calico.yaml
```

## Deploy with kube proxy mode: ipvs

[ipvs](https://kubernetes.io/blog/2018/07/09/ipvs-based-in-cluster-load-balancing-deep-dive/) "(IP Virtual Server) is built on top of the Netfilter and implements transport-layer load balancing as part of the Linux kernel."

```
kind: Cluster
apiVersion: kind.sigs.k8s.io/v1alpha3
nodes:
- role: control-plane
- role: worker
- role: worker
- role: worker
kubeadmConfigPatches:
- |
  apiVersion: kubeproxy.config.k8s.io/v1alpha1
  kind: KubeProxyConfiguration
  metadata:
    name: config
  mode: "ipvs"
``` 

## Deploy a specific version of Kubernets

```
kind create cluster --name project1 --config kind-calico.yaml --image kindest/node:v1.18.0
```

## Using kind to test kubernetes

### Install direnv  

```
# Follow instructions here. 
https://direnv.net/
```

### Install gimme

```
mkdir ~/bin
curl -sL -o ~/bin/gimme https://raw.githubusercontent.com/travis-ci/gimme/master/gimme
chmod +x ~/bin/gimme

gimmie -k # list of go version it can install
gimmie -l # current version
```

### Setup development environment

```
mkdir k8s-development
cd k8s-development
mkdir -p go/src

gimme stable

# gimme stable >> .envrc if you are using direnv

vi .envrc
export GOPATH="${PWD}/go"

go get k8s.io/kubernetes
cd $GOPATH/src/k8s.io/kubernetes
git tag # list current version available.

kind build node-image --image=master
kind create cluster --name test-master  --image master

# cleanup
docker system prune --volumes 

```


