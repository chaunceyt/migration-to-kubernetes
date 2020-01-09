# Google Kubernetes Engine

I've installed Kubernetes from scratch using `kubeadm`, `kops` `ansible` (OpenShift), and find that Google's Kubernetes Engine (GKE) to be my preferred Kubernetes environment to use.

## Provision GKE cluster

**NOTE**: This can be done a number of ways. Using the GCP Console, gcloud command, and terraform to name a few. For this document the gcloud command is used. 

Create script `create-cluster.sh` that will be responsible for creating the initial cluster.

```
#!/bin/bash

PROJECT=example-devcloud
CLUSTER_NAME=devcloud
REGION=us-east4
ZONE_LOCATION=us-east4-b
MACHINE_TYPE=n1-standard-4

gcloud config set project $PROJECT
gcloud config set compute/region $REGION
gcloud config set compute/zone $ZONE_LOCATION

gcloud beta container clusters create $CLUSTER_NAME  \
  --num-nodes=3 \
  --enable-ip-alias \
  --enable-network-policy \
  --machine-type=$MACHINE_TYPE  \
  --security-group="gke-security-groups@example.com" \
  --zone=$ZONE_LOCATION \
  --node-locations=$ZONE_LOCATION

```

Create script `create-builder-nodes.sh` that will create the "group" node pool for gitlab runner builds.

```
#!/bin/bash
# us-west1-b - production
# us-central1-c - testing
PROJECT=web-devcloud
CLUSTER_NAME=devcloud
REGION=us-east4
ZONE_LOCATION=us-east4-b
MACHINE_TYPE=n1-standard-8
NODE_POOL_NAME=builder-pool-n1s8

gcloud config set project $PROJECT
gcloud config set compute/region $REGION
gcloud config set compute/zone $ZONE_LOCATION

gcloud container node-pools create $NODE_POOL_NAME \
  --cluster=$CLUSTER_NAME \
  --machine-type=$MACHINE_TYPE \
  --zone=$ZONE_LOCATION \
  --num-nodes=2 \
  --node-taints=type=runner:NoSchedule \
  --node-labels=type=runner
```

Create script `create-project-node-pool.sh` to create a builder-pool for projects over a certian number of developers.

```
#!/bin/bash
# This script is responsible for creating a node-pool for the PROJECT_NAME

# Utility function: usage
function usage() {
  echo -e "usage: ${0} [project-name]\n"
	echo -e "example usage: ${0} webproject\n"
}

# Ensure we our project-name
if [ -z "${1}" ]
then
	echo "ERROR: Invalid project-name"
	usage
	exit 1
fi

# Set the projectname
PROJECT_NAME=$1

# Set GCP settings.
PROJECT=web-devcloud
CLUSTER_NAME=devcloud
REGION=us-east4
ZONE_LOCATION=us-east4-b
MACHINE_TYPE=n1-standard-2
NODE_POOL_NAME=builder-pool-${PROJECT_NAME}

gcloud config set project $PROJECT
gcloud config set compute/region $REGION
gcloud config set compute/zone $ZONE_LOCATION

# Create our node-pool.
gcloud container node-pools create $NODE_POOL_NAME \
  --cluster=$CLUSTER_NAME \
  --machine-type=$MACHINE_TYPE \
  --zone=$ZONE_LOCATION \
  --num-nodes=1 \
  --local-ssd-count=1 \
  --node-taints=type=${PROJECT_NAME}-runner:NoSchedule \
  --node-labels=type=${PROJECT_NAME}-runner
```
