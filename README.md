# Notes

I'm on a team responsible for migrating to a new CI environment. My role on this team is "Lead cluster administrator" being responsible for the installation and setup of the final production environment and day to day management of the environment. My current title within the company is Software Architect.

I treating my notes as a JIRA ticket created and assigned to me.

# Issue
**Subject:** Perform migration to Kubernetes

We have a number of OpenStack VMs allocated to CI for our projects. Each project has a Jenkins container that manages a develop, qa and review environment for the project. Each of the stated environments have multiple containers running.

Looking for something to manage our containers that is not docker-compose and Jenkins. We reviewed Docker Swarm, OpenShift, and GKE for container management settling on Google Kubernetes Engine (GKE) and Gitlab. We will use GKE to manage our containers and Gitlab will replace Jenkins.

We need to create a plan to migrate to GKE & Gitlab and implement said plan.

# Analysis

## Project Requirements
- Gitlab project
- TLS
- Authentication (Via google account or AWS Cognito)

## List of services

- Starter project with predefined `.gitlab-ci.yml` and Helm chart
- Apache or Nginx
- Mariadb
- Memcached or Redis
- Solr (custom subchart) or ElasticSearch
- Varnish
- Mailhog
- Extra domain with "basic auth" just for browerstack, cypress and other testing tools
- Not limited to current list of services.


## Actions


- Team created a cluster, deployed the needed components, connect the cluster to a Gitlab project and ran a project from the environment. This cluster had two node-pools builder and worker. RBAC Manager (per user)
- [Training](docs/TRAINING-PROFDEV.md) and Professional development. ([CKA/CKAD](certs/))
- Reviewed Knative serving considering scale-to-zero for resource management
- Reviewed Network Policies. GKE uses Tigera's [Calico](https://www.projectcalico.org/) [Demos](https://docs.projectcalico.org/v3.11/security/tutorials/kubernetes-policy-demo/kubernetes-demo)
- Reviewed [Velero](https://velero.io/) for DR and as migration tool. Allowing us to migrate to other cloud providers Kubernetes offering with little effort if we needed to.
- Reviewed SysDig Monitoring and watched their youtube playlist [here](https://www.youtube.com/playlist?list=PLrUjPk-W0lae7KuCFvmdbWj9Powm7Ryu0). (currently using in production)
- Reviewed Horizonal Pod Autoscaler ([HPA](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/))
- Team created migration plan
- Team invited other projects that wanted to help us better understand the requirements needed to run the cluster.
- Team member created starter project containing a Helm chart for the workloads
- Defined labels for namespaces (RBAC-Manager requirement)
- Created [new cluster](docs/GKE-SETUP.md) that would become the "production" one used for all projects. Requirements: RBAC with integration with Google Groups, network policies
- Deployed [RBAC-Manager](docs/RBAC-MANAGER.md) to manage RBAC Rolebindings (per google group)
- Team member created an external IP for the ingress controller
- Deployed [Ingress Nginx](docs/INGRESS-NGINX.md) to manage incoming traffic into the cluster
- Team discovered required annontations for certmanager and proxy-authentication
- Deployed [Prometheus & Grafana](docs/PROMETHEUS-GRAFANA.md) for observability (this was done via the GCP Console)
- Deployed [WeaveScope](docs/WEAVESCOPE.md) for observability
- Created helm chart to provision project's namespace
- Migrated projects off POC cluster to new cluster and deleted it.
- Use `stern` for log tailing since it monitors all of the containers in a pod
- Use `octant` as a graphical debugging tool



## Issues

- What RBAC permissions will allow project team members to debug their workloads without impacting other projects running in the cluster?
- Builder-pool performance (gitlab-runner) (WIP currently creating a builder-pool with localssd drives for projects with 6 or more developers)
- Cluster scaling up and down impacts API server availability resulting in Gitlab losing communication with the runner. Cause: GKE like AWS "automatically configures the proper VM size for your master depending on the number of nodes in your cluster" See [Size of master and master components](https://kubernetes.io/docs/setup/best-practices/cluster-large/#size-of-master-and-master-components) Solution: Disabled autoscaling for now.
- Ingress size of response (added annotation to resolved this)
- Upgrading cert-manager exposed a bug in nginx ingress that resulted in all our projects reverting to self-signed ssl certs. (Upgraded ingress nginx resolved issue)
- Allowing project tech leads to manage project namespace without another rolebinding. They need to be able to delete pods if they're in a regressed state. (WIP)



# Confirmation
1. A project is able to build their projects codebase and create a container image from that build and push to a registry.
2. After a container image is built a process triggers the creation of a web environment using the container image built during the build process.
3. A project should should have a domain per environment. i.e. `develop|qa|review.kube.domain.tld`
4. Each project should have an **Apache**, **Caching** and **Database** container by default.
5. Tech leads should be able to troubleshoot their project environments

## Proof of Concept

Re: Confirmation #2 "a process triggers the creation of a web environment..."

I started working a commandline tool that creates many of the services listed above. Initial start of that can be found [here](webproject-ctl).

Things missing

- ConfigMap for environmental variables
- Solr Search option
- ElasticSearch option
- Sidecar support
