# Notes

I'm on a team responsible for migrating to a new CI environment. I treating my notes as a JIRA ticket storing them here.

# Issue

We have a number of OpenStack VMs allocated to CI for our projects. Each project has a Jenkins container that manages a develop, qa and review environment for the project. Each of the stated environments have multiple containers running.

List of services currently offered:

- Apache or Nginx
- Mariadb
- Memcached or Redis
- Solr or ElasticSearch
- Varnish

Looking for something to manage our containers that is not docker-compose and Jenkins. We reviewed Docker Swarm, OpenShift, and GKE for container management settling on Google Kubernetes Engine (GKE) and Gitlab. We will use GKE to manage our containers and Gitlab will replace Jenkins.

We need to create a plan to migrate to GKE & Gitlab and implement said plan.

# Analysis

## Actions


- Created a cluster, deployed the needed components, connect the cluster to a Gitlab project and ran a project from the environment. This cluster had two node-pools builder and worker.
- Training and Professional development.
- Created Migration plan
- Invited other projects that wanted to help us better understand the requirements needed to run the cluster.
- Team member created starter project containing a Helm chart for the workloads
- Define labels for namespaces and workloads
- Created new cluster that would become the "production" one used for all projects. Requirements: RBAC with integration with Google Groups, network policies
- Deployed RBAC-Manager to manage RBAC Rolebindings
- Deployed Ingress Nginx to manage incoming traffic into the cluster
- Discovered required annontations for certmanager and proxy-authication
- Deployed Prometheus & Grafana for observability (this was done via the GCP Console)
- Deployed WeaveScope for observability
- Migrated projects off POC cluster to new cluster and deleted it.
- Use `stern` for log tailing since it monitors all of the containers in a pod
- Use `octant` as a graphical debugging tool


## Issues

- What RBAC permissions will allow project team members to debug their workloads without impacting other projects running in the cluster?
- Builder-pool performance (gitlab-runner) (WIP currently creating a builder-pool with localssd drives for projects with 6 or more developers)
- Cluster scaling up and down impacts API server availability resulting in Gitlab losing communication with the runner. GKE like AWS "automatically configures the proper VM size for your master depending on the number of nodes in your cluster" See [Size of master and master components](https://kubernetes.io/docs/setup/best-practices/cluster-large/#size-of-master-and-master-components)
- Ingress size of response (added annotation resolved this)
- Upgrading cert-manager exposed a bug in nginx ingress that resulted in all our projects reverting to self-signed ssl certs. (Upgraded nginx resolved issue)
- Allowing project tech leads to manage project namespace without another rolebinding. They need to be able to delete pods if they're in a regressed state. (WIP)



# Confirmation
1. A project is able to build their projects codebase and create a container image from that build and push to a registry.
2. After a container image is built a process triggers the creation of a web environment using the container image built during the build process.
3. A project should should have a domain per environment. i.e. `develop|qa|review.kube.domain.tld`
4. Each project should have an **Apache**, **Caching** and **Database** container by default.
5. Tech leads should be able to troubleshoot their project environments