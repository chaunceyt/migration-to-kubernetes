# Webproject Control

As I continue to teach myself Golang. This hack session I decided to take the [simple API endpoint](https://github.com/chaunceyt/go-examples/tree/master/webproject-api-using-gin) and make it a cli tool.

Initial start was to move to flags. Later I plan to move to "github.com/urfave/cli" later.

## Issue

I work in an environment that uses GitLab + GKE as the CI solution for our developers. We're currently using a helm chart to generate the Kubernetes manifests for each prprojects workload. These workloads contain the following services

- Apache/PHP7.x PHP-FPM
- Mariadb
- Memcached
- Redis
- Solr
- ElasticSearch

## Proof of concept

Create a commandline tool that accepts parameters instructing it to create an environment with similar components stated above. This is ALPHA work.

## Example usage at the moment.

```
releases/webproject-ctl-darwin-amd64 \
  -create \
  -deployment-name=vrt-manager \
  -primary-container-name=vrt-manager \
  -prinary-container-image-tag=chaunceyt/vrt-manager-httpd \
  -primary-container-port=8080 \
  -replicas=1 \
  -domain-name=project-d.kube.domain.tld \
  -namespace=webproject
```

## Output...

```
2020/01/13 20:40:37 Creating pvc...
2020/01/13 20:40:37 Created PVC - Name: "vrt-manager-webfiles-pvc", UID: "c49985c4-6857-4dd8-9bc0-1e37c39167d4"
2020/01/13 20:40:37 Unsupported CacheEngine selected or not defined
Creating webproject deployment...
Created Deployment - Name: "vrt-manager", UID: "6846f34c-717e-4c81-97d5-92b2e241a746"
Creating service for WebProject.
2020/01/13 20:40:37 Created Memcahed Deployment - Name: "vrt-manager-ing", UID: "4c5fccad-3fb1-4eac-8785-54212fefa1db"
Chaunceys-iMac:webprojectctl cthorn$ kubectl get po -n webproject
NAME                           READY   STATUS    RESTARTS   AGE
vrt-manager-6c9fdbdf69-7h2gs   1/1     Running   0          3s
Chaunceys-iMac:webprojectctl cthorn$ kubectl get all,pvc,ing -n webproject
NAME                               READY   STATUS    RESTARTS   AGE
pod/vrt-manager-6c9fdbdf69-7h2gs   1/1     Running   0          21s

NAME                      TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)   AGE
service/vrt-manager-svc   ClusterIP   10.105.153.105   <none>        80/TCP    6m30s

NAME                          READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/vrt-manager   1/1     1            1           21s

NAME                                     DESIRED   CURRENT   READY   AGE
replicaset.apps/vrt-manager-6c9fdbdf69   1         1         1       21s

NAME                                             STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
persistentvolumeclaim/vrt-manager-webfiles-pvc   Bound    pvc-c49985c4-6857-4dd8-9bc0-1e37c39167d4   1Gi        RWO            standard       21s

NAME                                 HOSTS                       ADDRESS   PORTS   AGE
ingress.extensions/vrt-manager-ing   project-d.kube.domain.tld             80      21s
```


Another execution added Redis as CacheEngine and Mysql as the Database engine...

```
  go run webproject/*.go -deployment-name=vrt-manager-03 -primary-container-name=vrt-manager -prinary-container-image-tag=chaunceyt/vrt-manager-httpd -primary-container-port=8080 -replicas=2 -domain-name=project-d.kube.domain.tld -namespace=project-x -cache-engine=redis -database-engine=mariadb -database-engine-image=mysql:5.7
2020/01/13 21:16:36 Creating pvc...
2020/01/13 21:16:37 Created PVC - Name: "vrt-manager-03-webfiles-pvc", UID: "aa5d3721-1514-44c0-915a-42acb0c68ce3"
2020/01/13 21:16:37 Creating pvc...
2020/01/13 21:16:37 Created PVC - Name: "vrt-manager-03-db-pvc", UID: "2f72f42d-d35d-4b47-ac5c-638ea114f0e2"
2020/01/13 21:16:37 Creating database deployment...
2020/01/13 21:16:37 Created Database Deployment - Name: "vrt-manager-03-db", UID: "cc3d9de5-128c-46f5-9c6e-34e4a75ca995"
2020/01/13 21:16:37 Creating redis deployment...
2020/01/13 21:16:37 Created Redis Deployment - Name: "vrt-manager-03-redis", UID: "6f526030-55b7-4ae2-a626-07c99616b341"
2020/01/13 21:16:37 Creating redis service...
2020/01/13 21:16:37 Created Redis Service - Name: "vrt-manager-03-redis-svc", UID: "ba409c15-1ee2-43dd-93aa-64e2fcfa9091"
Creating webproject deployment...
Created Deployment - Name: "vrt-manager-03", UID: "25ec5f35-e21e-4fa3-8983-55ac0c721d63"
Creating service for WebProject.
Created Webproject Service - Name: "vrt-manager-03-svc", UID: "23737a94-6a92-40dd-a9a9-a2d8ddab6a17"
2020/01/13 21:16:38 Created Memcahed Deployment - Name: "vrt-manager-03-ing", UID: "ba67993f-7442-455b-b0a6-29a1841d5f30"
```


