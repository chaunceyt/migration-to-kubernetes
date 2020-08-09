# Kubernetes kubelet notes

The **kubelet** is a daemon that runs on each node within a Kubernetes cluster. 

- [Command line reference](https://kubernetes.io/docs/reference/command-line-tools-reference/kubelet/)
- [Configuring each kubelet in your cluster using kubeadm](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/kubelet-integration/)

##Responsibility

Manage pods that have a `nodeName:` that matches their nodeName.

## Run in standalone mode

Use the kubelet to manage pods without a Kubernetes cluster.


### Centos
https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/

```
cat <<EOF | sudo tee /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-\$basearch
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
exclude=kubelet kubeadm kubectl
EOF
```

```
# Set SELinux in permissive mode (effectively disabling it)
sudo setenforce 0
sudo sed -i 's/^SELINUX=enforcing$/SELINUX=permissive/' /etc/selinux/config
sudo yum install -y kubelet --disableexcludes=kubernetes
sudo systemctl enable --now kubelet
```

Update systemd file: `vi /usr/lib/systemd/system/kubelet.service`

```
# update ExecStart
ExecStart=/usr/bin/kubelet \
  --file-check-frequency 30s \
  --max-pods 10 \
  --minimum-image-ttl-duration 300s \
  --pod-manifest-path=/etc/kubernetes/manifests \
  --sync-frequency 30s \
  --cgroup-driver=systemd
```

Run the following commands

```
systemctl daemon-relaod
systemctl start kubelet
systemctl status kubelet
systemctl enable kubelet
```

Create a static pod and copy to the `--pod-manifest-path`

```
apiVersion: v1
kind: Pod
metadata:
  name: webapp
  labels:
    role: static-pod
spec:
   initContainers:
     - name: msginit
       image: centos:7
       command:
       - "bin/bash"
       - "-c"
       - "echo INIT_DONE > /ic/this"
       volumeMounts:
       - mountPath: /ic
         name: logs
  containers:
    - name: web
      image: nginx
      livenessProbe:
        httpGet:
          path: /probe.txt
          port: 80
        initialDelaySeconds: 2
        periodSeconds: 2
      readinessProbe:
        tcpSocket:
          port: 80
        initialDelaySeconds: 2
        periodSeconds: 2
      lifecycle:
        postStart:
          exec:
            command: ["/bin/sh", "-c", "echo Hello from the postStart handler > /usr/share/message"]
        preStop:
          exec:
            command: ["/bin/sh","-c","nginx -s quit; while killall -0 nginx; do sleep 1; done"]
      ports:
        - name: web
          containerPort: 80
          hostPort: 8081
          protocol: TCP
      volumeMounts:
      - name: docroot
        mountPath: /usr/share/nginx/html
      - name: logs
        mountPath: /ic
  volumes:
    - name: logs
      emptyDir: {}
    - name: docroot
      hostPath:
        path: /mnt/project1/docroot
```

