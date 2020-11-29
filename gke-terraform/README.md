# Create GKE cluster

* 10-cluster-configuration.tf - create GKE cluster
* 11-builder-configuration.tf - create GKE node-pool with taint for gitlab-runners only
* 11-worker-configuration.tf - create GKE node-pool to run project workloads

```
# Validate
terraform validate
terraform fmt -check=true

# Review
alias convert_report="jq -r '([.resource_changes[].change.actions?]|flatten)|{\"create\":(map(select(.==\"create\"))|length),\"update\":(map(select(.==\"update\"))|length),\"delete\":(map(select(.==\"delete\"))|length)}'"

terraform plan -out=plan.tfplan
terraform show --json $PLAN | convert_report > tfplan.json

# Plan
terraform plan

# Apply
terraform apply -auto-approve

# Destroy
terraform destroy -auto-approve

```

## Trigger cluster node autoscale event

Using n1-standard-1 instances

```
kubectl create deploy nginx --image=nginx
kubectl scale deploy nginx --replicas=240

```
