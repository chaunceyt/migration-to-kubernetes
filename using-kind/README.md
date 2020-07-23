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

### SPIKE: Gatekeeper - Policy Controller for Kubernetes

- [Gatekeeper](https://github.com/open-policy-agent/gatekeeper)
- [Webinar: K8s with OPA Gatekeeper](https://www.youtube.com/watch?v=v4wJE3I8BYM)

```
kind create cluster --name opa-gatekeeper --config kind-config.yaml --image kindest/node:v1.18.0

# Download the gatekeeper manifest and review what is being created.
wget https://raw.githubusercontent.com/open-policy-agent/gatekeeper/master/deploy/gatekeeper.yaml

# install
# https://github.com/open-policy-agent/gatekeeper#how-to-use-gatekeeper
kubectl apply -f gatekeeper.yaml

# work through the libraries to better understand how the ConstraintTemplate work.
# https://github.com/open-policy-agent/gatekeeper#constraint-templates
git clone https://github.com/open-policy-agent/gatekeeper.git
cd gatekeeper/library/pod-security-policy
bash test.sh

```

### SPIKE: Vertical Pod Autoscaling

Vertical Pod autoscaling involves adjusting a Pod's CPU and memory requests

```
mkdir ~/development/kind-metrics-server
cd ~/development/kind-metrics-server
wget https://github.com/kubernetes-sigs/metrics-server/releases/download/v0.3.6/components.yaml
# edit the deployment adding the following directives
# - --kubelet-insecure-tls
# - --kubelet-preferred-address-types=InternalIP
# or run kubectl patch deployment metrics-server -n kube-system -p '{"spec":{"template":{"spec":{"containers":[{"name":"metrics-server","args":["--cert-dir=/tmp", "--secure-port=4443", "--kubelet-insecure-tls","--kubelet-preferred-address-types=InternalIP"]}]}}}}'

kubectl apply -f components.yaml
# Wait for awhile to get some metrics
kubectl top no

# Install VPA
git clone https://github.com/kubernetes/autoscaler
cd autoscaler/vertical-pod-autoscaler/
cd hack
./vpa-up.sh
# Walk through the examples here
https://cloud.google.com/kubernetes-engine/docs/how-to/vertical-pod-autoscaling

Modes:
- Recommendation 
- Auto

```

### SPIKE: Horizontal Pod Autoscaler

The Horizontal Pod Autoscaler automatically scales the number of pods in a replication controller, deployment, replica set or stateful set based on observed CPU utilization 

