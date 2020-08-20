# Notes on Kubernetes secrets

- [Documentation](https://kubernetes.io/docs/concepts/configuration/secret/)
- ["Sealed Secrets" for Kubernetes](https://github.com/bitnami-labs/sealed-secrets)
- [Managing secrets deployment in Kubernetes using Sealed Secrets](https://aws.amazon.com/blogs/opensource/managing-secrets-deployment-in-kubernetes-using-sealed-secrets/)
- [TGI Kubernetes 120: CSI and Secrets!](https://www.youtube.com/watch?v=IznsHhKL428)
- [TGI Kubernetes 113: Kubernetes Secrets Take 3](https://www.youtube.com/watch?v=an9D2FyFwR0)

## Install and test.

```
helm fetch stable/sealed-secrets --untar
helm template --name sealed-secrets-controller --namespace kube-system sealed-secrets/ | kubectl -n kube-system apply -f -
wget https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.12.5/kubeseal-darwin-amd64
mv kubeseal-darwin-amd64 $HOME/bin/kubeseal
chmod +x $HOME/bin/kubeseal

kubectl create secret generic project-secret --dry-run --from-literal=username=admin --dry-run --from-literal=password=4P@s5wOrd2N0 -o yaml | kubeseal > project-secret.yaml

```

