# Kubernetes Container Storage Interface

[Documentation](https://kubernetes-csi.github.io/docs/introduction.html)

- [external-resizer](https://kubernetes-csi.github.io/docs/external-resizer.html) `allowVolumeExpansion: true` should be in the output of `kubectl get sc -o yaml`
- [external-snapshotter](https://kubernetes-csi.github.io/docs/external-snapshotter.html)
- [volume cloning](https://kubernetes.io/docs/concepts/storage/volume-pvc-datasource/)

## Implementing Snapshotter

Problem: Project having databases 10GB or larger needing to seed an issue branch workload. Takes a long time to import database. 

Solution: Use `VolumeSnapshots` until cloning is available.

### Starting with a new cluster

Create GKE cluster with the addon `GcePersistentDiskCsiDriver` enabled

```
gcloud beta container \
	--project "[PROJECT_NAME]" clusters create "cthorn-test-01" \
	--zone "us-central1-c" \
	--node-locations "us-central1-c" \
	--no-enable-basic-auth \
	--cluster-version "1.18.10-gke.601" \
	--release-channel "rapid" \
	--machine-type "e2-medium" \
	--image-type "COS" \
	--disk-type "pd-standard" \
	--disk-size "100" \
	--metadata disable-legacy-endpoints=true \
	--scopes "https://www.googleapis.com/auth/devstorage.read_only","https://www.googleapis.com/auth/logging.write","https://www.googleapis.com/auth/monitoring","https://www.googleapis.com/auth/servicecontrol","https://www.googleapis.com/auth/service.management.readonly","https://www.googleapis.com/auth/trace.append" \
	--num-nodes "2" \
	--enable-ip-alias \
	--enable-network-policy \
	--no-enable-master-authorized-networks \
	--addons HorizontalPodAutoscaling,HttpLoadBalancing,GcePersistentDiskCsiDriver \
	--enable-autoupgrade \
	--enable-autorepair \
	--enable-vertical-pod-autoscaling
```

When `GcePersistentDiskCsiDriver` is enabled a number of `StorageClasses` are created.

```
kubectl get sc -l k8s-app=gcp-compute-persistent-disk-csi-driver
NAME           PROVISIONER             RECLAIMPOLICY   VOLUMEBINDINGMODE      ALLOWVOLUMEEXPANSION   AGE
premium-rwo    pd.csi.storage.gke.io   Delete          WaitForFirstConsumer   true                   2m37s
standard-rwo   pd.csi.storage.gke.io   Delete          WaitForFirstConsumer   true                   2m37s
```

Create snapshotclass that uses the provisioner `pd.csi.storage.gke.io`.

```
apiVersion: snapshot.storage.k8s.io/v1beta1
kind: VolumeSnapshotClass
metadata:
  name: csi-pd-snapclass
driver: pd.csi.storage.gke.io
deletionPolicy: Delete
```


Create a persistent volume claim (PVC) that uses the storageclass `standard-rwo`

```
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: webapp-docroot
  namespace: project-01
spec:
  accessModes:
  - ReadWriteOnce
  storageClassName: standard-rwo
  resources:
    requests:
      storage: 6Gi
```

Create a pod that requests the PVC created

```
apiVersion: v1
kind: Pod
metadata:
  name: web-server
  namespace: project-01
spec:
  containers:
   - name: web-server
     image: nginx
     volumeMounts:
       - mountPath: /var/lib/www/html
         name: docroot
  volumes:
   - name: docroot
     persistentVolumeClaim:
       claimName: webapp-docroot
       readOnly: false
```


Create a snapshot of the web-server PVC 

```
apiVersion: snapshot.storage.k8s.io/v1beta1
kind: VolumeSnapshot
metadata:
  name: webapp-docroot-snapshot
  namespace: project-01
spec:
  volumeSnapshotClassName: csi-pd-snapclass
  source:
    persistentVolumeClaimName: webapp-docroot
```

To use a snapshot of an existing PVC our PVC will look like this. NOTE the `dataSource` directive.

```
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: issue-000-webapp-docroot
spec:
  storageClassName: standard-rwo
  dataSource:
    name: webapp-docroot-snapshot
    kind: VolumeSnapshot
    apiGroup: snapshot.storage.k8s.io
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 6Gi
```
 
Create a pod to use the new PVC

```
apiVersion: v1
kind: Pod
metadata:
  name: issue-01-web-server
  namespace: project-01
spec:
  containers:
   - name: web-server
     image: nginx
     volumeMounts:
       - mountPath: /var/lib/www/html
         name: docroot
  volumes:
   - name: docroot
     persistentVolumeClaim:
       claimName: issue-000-webapp-docroot
       readOnly: false
```

## Workflow?
 
- issue/bug/no-ticket branches will use submit a volubemsnapshot request of the main branches database pv. Then submit a PVC with `datasource:` pointing to that snapshot.
- Deploy workload with `.spec.volumes.persistenVolmeClain.claimName` pointing to PVC created from volumespanshot