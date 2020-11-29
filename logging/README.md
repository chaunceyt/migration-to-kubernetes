# Logging using Elasticsearch, Kibana and fluent-bit

```
# Install elasticsearch and kibana
kubectl apply --recursive -f manifests

# Create a workload with fluent-bit as sidecar
kubectl apply --recursive -f example
```

## To quickly see Kibana output

`kubectl -n kube-logging port-forward service/kibana 5601:5601`

goto [http://localhost:5601/](http://localhost:5601/)
