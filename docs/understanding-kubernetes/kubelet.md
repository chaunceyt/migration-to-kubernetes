# Kubernetes kubelet notes

The **kubelet** is a daemon that runs on each node within a Kubernetes cluster. 

- [Command line reference](https://kubernetes.io/docs/reference/command-line-tools-reference/kubelet/)
- [Configuring each kubelet in your cluster using kubeadm](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/kubelet-integration/)
- [TGI Kubernetes 086: Grokking Kubernetes - The kubelet](https://www.youtube.com/watch?v=CKpSyl5vgK8) ([shownotess](https://github.com/vmware-tanzu/tgik/tree/master/episodes/086))

##Responsibility

Manage pods that have a `nodeName:` that matches their nodeName.

The following `kubectl` commands are "expressed and implemented" by the kubelet's [api](https://github.com/kubernetes/kubernetes/blob/master/pkg/kubelet/server/server.go). The Kubernetes API proxies these commands to kubelet


- exec ([how kubectl exec works](https://erkanerol.github.io/post/how-kubectl-exec-works/))
- attach
- cp
- log


## Default permissions within a cluster

When cluster is provisioned using `kubeadm`. 

```
 kubectl --kubeconfig=/etc/kubernetes/kubelet.conf auth can-i --list
Resources                                                       Non-Resource URLs   Resource Names   Verbs
selfsubjectaccessreviews.authorization.k8s.io                   []                  []               [create]
selfsubjectrulesreviews.authorization.k8s.io                    []                  []               [create]
certificatesigningrequests.certificates.k8s.io/selfnodeclient   []                  []               [create]
                                                                [/api/*]            []               [get]
                                                                [/api]              []               [get]
                                                                [/apis/*]           []               [get]
                                                                [/apis]             []               [get]
                                                                [/healthz]          []               [get]
                                                                [/healthz]          []               [get]
                                                                [/livez]            []               [get]
                                                                [/livez]            []               [get]
                                                                [/openapi/*]        []               [get]
                                                                [/openapi]          []               [get]
                                                                [/readyz]           []               [get]
                                                                [/readyz]           []               [get]
                                                                [/version/]         []               [get]
                                                                [/version/]         []               [get]
                                                                [/version]          []               [get]
                                                                [/version]          []               [get]
```

## Run in standalone mode

Use the kubelet to manage pods without a Kubernetes cluster.


### Centos 7
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
# Ensure kubelet can pull down private images.
# Config to use /root/.docker/config.json
Environment=HOME=/root
User=root
ExecStart=/usr/bin/kubelet \
  --address=127.0.0.1 \
  --hostname-override=127.0.0.1 \
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

```
# https://github.com/kubernetes/kubernetes/blob/master/pkg/kubelet/server/server.go#L328
curl -L -sk https://localhost:10250/pods

# get prom2json installed
curl -L -sk https://localhost:10250/healthz/log
curl -L -sk https://localhost:10250/healthz/ping
curl -L -sk https://localhost:10250/stats/summary
curl -L -sk https://localhost:10250/metrics
curl -L -sk https://localhost:10250/metrics/cadvisor
curl -L -sk https://localhost:10250/metrics/resource

# Logs on host system
curl -L -sk https://localhost:10250/logs
curl -L -sk https://localhost:10250/logs/syslog

# https://github.com/kubernetes/kubernetes/blob/master/pkg/kubelet/server/server.go#L484
curl -L -sk https://localhost:10250/containerLogs/[namspace]/[podName]/[container]

# Run command in a container
# https://github.com/kubernetes/kubernetes/blob/master/pkg/kubelet/server/server.go#L419
curl -sk -XPOST -d 'cmd=ls -ltr /' https://localhost:10250/run/[namspace]/[podName]/[container]

```


## From source

```
# ensure go, make, gcc are installed
cd ~/path/to/development/directory
git clone https://github.com/kubernetes/kubernetes
cd kubernetes
make WHAT=cmd/kubelet
mkdir /etc/kubelet-standalone
cp _output/bin/kubelet /usr/local/bin/kubelet

# ensure executable is in your path.
kubelet -h

# update the systemd file above.

```

`/etc/kubelet-standalone/config.yaml` for `kubelet --config=/etc/kubelet-standalone/config.yaml` 

```
kind: KubeletConfiguration
apiVersion: kubelet.config.k8s.io/v1beta1
address: 0.0.0.0
staticPodPath: /etc/kubernetes/manifests
authentication:
 webhook:
   enabled: false
 anonymous:
   enabled: true
authorization:
  mode: AlwaysAllow
maxPods: 5
containerRuntime: docker
failSwapFalse: false
serializeImagePulls: true
fileCheckFrequency: 20s
```