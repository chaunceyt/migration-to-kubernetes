# Notes...

## Deep Dive

Goal: Understand how to use and support prometheus and alertmanager in a production environment.


- Highly available? https://prometheus.io/docs/introduction/faq/#can-prometheus-be-made-highly-available
- Federation https://prometheus.io/docs/prometheus/latest/federation/
- Hierarchical federation collect time-series data from multiple lower-level servers
- Cross-service federation scrapes selected data from another server
- Security https://prometheus.io/docs/operating/security/
- Client Libraries https://prometheus.io/docs/instrumenting/clientlibs/
- PHP https://github.com/endclothing/prometheus_client_php
- Go https://github.com/prometheus/client_golang


## Usage patterns

- Metric collection
- Visualization
- Alerting


#### Command line options

`prometheus-2.19.2.darwin-amd64/prometheus -h`

```
usage: prometheus [<flags>]

The Prometheus monitoring server

Flags:
  -h, --help                     Show context-sensitive help (also try --help-long and --help-man).
      --version                  Show application version.
      --config.file="prometheus.yml"
                                 Prometheus configuration file path.
      --web.listen-address="0.0.0.0:9090"
                                 Address to listen on for UI, API, and telemetry.
      --web.read-timeout=5m      Maximum duration before timing out read of the request, and closing idle connections.
      --web.max-connections=512  Maximum number of simultaneous connections.
      --web.external-url=<URL>   The URL under which Prometheus is externally reachable (for example, if Prometheus is served via a reverse proxy). Used for generating
                                 relative and absolute links back to Prometheus itself. If the URL has a path portion, it will be used to prefix all HTTP endpoints served
                                 by Prometheus. If omitted, relevant URL components will be derived automatically.
      --web.route-prefix=<path>  Prefix for the internal routes of web endpoints. Defaults to path of --web.external-url.
      --web.user-assets=<path>   Path to static asset directory, available at /user.
      --web.enable-lifecycle     Enable shutdown and reload via HTTP request.
      --web.enable-admin-api     Enable API endpoints for admin control actions.
      --web.console.templates="consoles"
                                 Path to the console template directory, available at /consoles.
      --web.console.libraries="console_libraries"
                                 Path to the console library directory.
      --web.page-title="Prometheus Time Series Collection and Processing Server"
                                 Document title of Prometheus instance.
      --web.cors.origin=".*"     Regex for CORS origin. It is fully anchored. Example: 'https?://(domain1|domain2)\.com'
      --storage.tsdb.path="data/"
                                 Base path for metrics storage.
      --storage.tsdb.retention=STORAGE.TSDB.RETENTION
                                 [DEPRECATED] How long to retain samples in storage. This flag has been deprecated, use "storage.tsdb.retention.time" instead.
      --storage.tsdb.retention.time=STORAGE.TSDB.RETENTION.TIME
                                 How long to retain samples in storage. When this flag is set it overrides "storage.tsdb.retention". If neither this flag nor
                                 "storage.tsdb.retention" nor "storage.tsdb.retention.size" is set, the retention time defaults to 15d. Units Supported: y, w, d, h, m, s,
                                 ms.
      --storage.tsdb.retention.size=STORAGE.TSDB.RETENTION.SIZE
                                 [EXPERIMENTAL] Maximum number of bytes that can be stored for blocks. A unit is required, supported units: B, KB, MB, GB, TB, PB, EB. Ex:
                                 "512MB". This flag is experimental and can be changed in future releases.
      --storage.tsdb.no-lockfile
                                 Do not create lockfile in data directory.
      --storage.tsdb.allow-overlapping-blocks
                                 [EXPERIMENTAL] Allow overlapping blocks, which in turn enables vertical compaction and vertical query merge.
      --storage.tsdb.wal-compression
                                 Compress the tsdb WAL.
      --storage.remote.flush-deadline=<duration>
                                 How long to wait flushing sample on shutdown or config reload.
      --storage.remote.read-sample-limit=5e7
                                 Maximum overall number of samples to return via the remote read interface, in a single query. 0 means no limit. This limit is ignored for
                                 streamed response types.
      --storage.remote.read-concurrent-limit=10
                                 Maximum number of concurrent remote read calls. 0 means no limit.
      --storage.remote.read-max-bytes-in-frame=1048576
                                 Maximum number of bytes in a single frame for streaming remote read response types before marshalling. Note that client might have limit
                                 on frame size as well. 1MB as recommended by protobuf by default.
      --rules.alert.for-outage-tolerance=1h
                                 Max time to tolerate prometheus outage for restoring "for" state of alert.
      --rules.alert.for-grace-period=10m
                                 Minimum duration between alert and restored "for" state. This is maintained only for alerts with configured "for" time greater than grace
                                 period.
      --rules.alert.resend-delay=1m
                                 Minimum amount of time to wait before resending an alert to Alertmanager.
      --alertmanager.notification-queue-capacity=10000
                                 The capacity of the queue for pending Alertmanager notifications.
      --alertmanager.timeout=10s
                                 Timeout for sending alerts to Alertmanager.
      --query.lookback-delta=5m  The maximum lookback duration for retrieving metrics during expression evaluations and federation.
      --query.timeout=2m         Maximum time a query may take before being aborted.
      --query.max-concurrency=20
                                 Maximum number of queries executed concurrently.
      --query.max-samples=50000000
                                 Maximum number of samples a single query can load into memory. Note that queries will fail if they try to load more samples than this
                                 into memory, so this also limits the number of samples a query can return.
      --log.level=info           Only log messages with the given severity or above. One of: [debug, info, warn, error]
      --log.format=logfmt        Output format of log messages. One of: [logfmt, json]
```

#### promtool

`prometheus-2.19.2.darwin-amd64/promtool -h`

```

usage: promtool [<flags>] <command> [<args> ...]

Tooling for the Prometheus monitoring system.

Flags:
  -h, --help     Show context-sensitive help (also try --help-long and --help-man).
      --version  Show application version.

Commands:
  help [<command>...]
    Show help.

  check config <config-files>...
    Check if the config files are valid or not.

  check rules <rule-files>...
    Check if the rule files are valid or not.

  check metrics
    Pass Prometheus metrics over stdin to lint them for consistency and correctness.

    examples:

    $ cat metrics.prom | promtool check metrics

    $ curl -s http://localhost:9090/metrics | promtool check metrics

  query instant <server> <expr>
    Run instant query.

  query range [<flags>] <server> <expr>
    Run range query.

  query series --match=MATCH [<flags>] <server>
    Run series query.

  query labels <server> <name>
    Run labels query.

  debug pprof <server>
    Fetch profiling debug information.

  debug metrics <server>
    Fetch metrics debug information.

  debug all <server>
    Fetch all debug information.

  test rules <test-rule-file>...
    Unit tests for rules.
    
```    

Example usage and output.

```
prometheus-2.19.2.darwin-amd64/promtool query instant http://localhost:9090 up
up{instance="localhost:9091", job="pushgateway"} => 1 @[1595510365.04]
up{instance="localhost:8080", job="cadvisor"} => 1 @[1595510365.04]
up{instance="localhost:9090", job="prometheus"} => 1 @[1595510365.04]
up{instance="localhost:9100", job="workstation"} => 1 @[1595510365.04]
```


#### Run binaries locally to see how things work together.

```

mkdir -p prometheus/temp
cd prometheus

wget https://github.com/prometheus/prometheus/releases/download/v2.19.2/prometheus-2.19.2.darwin-amd64.tar.gz
tar -xvzf prometheus-2.19.2.darwin-amd64.tar.gz

wget https://github.com/prometheus/alertmanager/releases/download/v0.21.0/alertmanager-0.21.0.darwin-amd64.tar.gz
tar -xvzf alertmanager-0.21.0.darwin-amd64.tar.gz

wget https://github.com/prometheus/pushgateway/releases/download/v1.2.0/pushgateway-1.2.0.darwin-amd64.tar.gz
tar -xvzf pushgateway-1.2.0.darwin-amd64.tar.gz

mkdir -p temp/{consoles,console_libraries}

# Bash script start-prometheus.sh
#!/bin/bash
PROMETHEUS_DIR=./prometheus-2.19.2.darwin-amd64
${PROMETHEUS_DIR}/prometheus --config.file ${PROMETHEUS_DIR}/prometheus.yml \
    --storage.tsdb.path ${PROMETHEUS_DIR}/data \
    --web.console.templates=${PROMETHEUS_DIR}/consoles \
    --web.console.libraries=${PROMETHEUS_DIR}/console_libraries

# Bash script start-alertmanager.sh
#!/bin/bash

ALERTMANAGER_DIR=./alertmanager-0.21.0.darwin-amd64
${ALERTMANAGER_DIR}/alertmanager \
  --config.file ${ALERTMANAGER_DIR}/alertmanager.yml \
  --storage.path ${ALERTMANAGER_DIR}/data

# Bash script start-node-exporter.sh
#!/bin/bash

NODE_EXPORTER_DIR=./node_exporter-1.0.1.darwin-amd64
${NODE_EXPORTER_DIR}/node_exporter

# Bash script start-pushgateway.sh
#!/bin/bash

PUSHGATEWAY_DIR=./pushgateway-1.2.0.darwin-amd64
${PUSHGATEWAY_DIR}/pushgateway

# Test 
curl localhost:9090/metrics
```

#### Exporters

Provide metric data that is collected by prometheus

- https://prometheus.io/docs/instrumenting/exporters/
- https://github.com/prometheus/node_exporter/blob/master/README.md

#### Configure an exporter

#### node_exporter

- Scrape config https://prometheus.io/docs/prometheus/latest/configuration/configuration/#scrape_config
- Monitoring a system https://prometheus.io/docs/guides/node-exporter/

```
wget https://github.com/prometheus/node_exporter/releases/download/v1.0.1/node_exporter-1.0.1.darwin-amd64.tar.gz
tat -xvzf node_exporter-1.0.1.darwin-amd64.tar.gz
./node_exporter-1.0.1.darwin-amd64/node_exporter
curl localhost:9100/metrics
```


##### Application monitoring

Apache Exporter: https://github.com/Lusitaniae/apache_exporter

```
sudo vi /etc/systemd/system/apache_exporter.service

[Unit]
Description=Prometheus Apache Exporter
Wants=network-online.target
After=network-online.target

[Service]
User=apache_exporter
Group=apache_exporter
Type=simple
ExecStart=/usr/local/bin/apache_exporter

[Install]
WantedBy=multi-user.target
```

```
sudo systemctl enable apache_exporter
sudo systemctl start apache_exporter
sudo systemctl status apache_exporter
curl localhost:9117/metrics
```


#### Configure Prometheus to Scrape Metrics from an exporter

`vi prometheus.yml` addiing workstation to prometheus config.

```
...

- job_name: 'Workstation'
  static_configs:
  - targets: ['localhost:9100']

- job_name: 'Apache'
    static_configs:
    - targets: ['<APACHE_SERVER_PRIVATE_IP>:9117']
...

```



### Prometheus data model

https://prometheus.io/docs/concepts/data_model/

Prometheus uses a combo of labels and metrics to identify each set of time-series data.

https://prometheus.io/docs/concepts/data_model/#metric-names-and-labels


`node_cpu_seconds_total(cpu="0", jobs="My Server")`


### What Is Time-Series Data?

"Time series data is a collection of quantities that are assembled over even intervals in time and ordered chronologically. The time interval at which data is collection is generally referred to as the time series frequency." 

"Consist of a series of values associated with different points in time."

### Metric Types

https://prometheus.io/docs/concepts/metric_types/
https://prometheus.io/docs/practices/histograms/

- counter - never descrease only increases
- gauge (values can both increase and decrease)

- histograms (more complex)
- historgram sum metric
- historgram count metrics


- quantile metrics
- quantile sum
- quantile count 


### Prometheus Querying
https://prometheus.io/docs/prometheus/latest/querying/basics/

Expression browser
Prometheus http api
visualization tools (grafana)

Examples

`node_cpu_seconds_total`

#### with label(s)

```
node_cpu_seconds_total{cpu="0"} 
node_cpu_seconds_total{cpu!="0"}
node_cpu_seconds_total{mode=~"s.*"}
```

####  range vector selector

```
node_cpu_seconds_total{cpu="0"}[2m] # get vector selector with time series over period of time
node_cpu_seconds_total{cpu="0"} offset 1h # get data from the past with offset modifier
node_cpu_seconds_total{cpu="0"}[5m] offset 1h # get range vector selector with offset modifier

node_cpu_seconds_total{mode=~"user|system"}
node_cpu_seconds_total{mode!~"user|system"}
```


### Operators

https://prometheus.io/docs/prometheus/latest/querying/operators/


Arithmetic binary operators

- + (addition)
- - (subtraction)
- * (multiplication)
- / (division)
- % (modulo)
- ^ (power/exponentiation)

Comparison binary operators


- == (equal)
- != (not-equal)
- > (greater-than)
- < (less-than)
- >= (greater-or-equal)
- <= (less-or-equal)

#### Query Functions

https://prometheus.io/docs/prometheus/latest/querying/functions/


#### HTTP API

https://prometheus.io/docs/prometheus/latest/querying/api/

For cli usage install `jq` to cleanup json response.

```
curl localhost:9090/api/v1/query?query=node_cpu_seconds_total

curl localhost:9090/api/v1/query --data-urlencode "query=node_cpu_seconds_total{cpu=\"0\"}"

sum(rate(node_cpu_seconds_total{job="Linux Server1",mode!='idle'}[5m])) * 100 / 2

sum(rate(node_cpu_seconds_total{job="Linux Server1",mode='idle'}[5m])) by (cpu)

```

Getting data based on range.

```
start=$(date --date '-5 min' +'%Y-%m-%dT%H:%M:%SZ')
end=$(date +'%Y-%m-%dT%H:%M:%SZ')
curl "localhost:9090/api/v1/query_range?query=node_cpu_seconds_total&start=$start&end=$end&step=1m"
```


### Using prometheus visualization methods

- Expression brower https://prometheus.io/docs/visualization/browser/ (native) build and debug queries
- Grafana (external + most common) https://prometheus.io/docs/visualization/grafana/
- Console templates https://prometheus.io/docs/visualization/consoles/ (native)
-  Console template graph library


The "moment" input allows one to control data
Note "graph" view - does not work with range vectors. Have to use instant vector query. i.e. rate()

#### Console templates

Path: `/etc/prometheus/consoles/`

Allows one to create visulization consoles using go template language to create simple HTML files. Add directive when starting prometheus server.

`--web.console.templates=/etc/prometheus/consoles/`

Examples: 

- https://prometheus.io/docs/visualization/consoles/#example-console

#### Console template graph library

- https://prometheus.io/docs/visualization/consoles/#graph-library


### Jobs and Instances

https://prometheus.io/docs/concepts/jobs_instances/

- Instances are individual endpoints prometheus scrapes Usually an instance of a single application
- Jobs a collection of instances, sharing a single purpose. 

### Scrape meta-metrics 

- `scrape_duration_seconds`


### Enable docker to serve metrics to prometheus

On a Linux host `sudo vi /etc/docker/daemon.json`

```
{
  "experimental": true,
  "metrics-addr": "IP_ADDRESS:9323"
}
```

```
sudo systemctl restart docker
curl IP_ADDRESS:9323/metrics
```


### Enable container metrics monitoring

```
docker run -d --restart always --name cadvisor -p 8080:8080 -v "/:/rootfs:ro" -v "/var/run:/var/run:rw" -v "/sys:/sys:ro" -v "/var/lib/docker/:/var/lib/docker:ro" google/cadvisor:latest
```

### Pushgateway

A solution for a push-based model. The default is a "pull" method. An use case for push - a batch job process. In this context prometheus becomes a "man in the middle". Allowing client to push to prometheus.

- https://prometheus.io/docs/instrumenting/pushing/
- Github repo: https://github.com/prometheus/pushgateway
- When to use pushgateway https://prometheus.io/docs/practices/pushing/

#### Installing pushgateway

```
sudo useradd -M -r -s /bin/false pushgateway
wget https://github.com/prometheus/pushgateway/releases/download/v1.2.0/pushgateway-1.2.0.linux-amd64.tar.gz
tar xvfz pushgateway-1.2.0.linux-amd64.tar.gz
sudo cp pushgateway-1.2.0.linux-amd64/pushgateway /usr/local/bin/
sudo chown pushgateway:pushgateway /usr/local/bin/pushgateway
```

`sudo vi /etc/systemd/system/pushgateway.service`

```
[Unit]
Description=Prometheus Pushgateway
Wants=network-online.target
After=network-online.target

[Service]
User=pushgateway
Group=pushgateway
Type=simple
ExecStart=/usr/local/bin/pushgateway

[Install]
WantedBy=multi-user.target
```

```
sudo systemctl enable pushgateway
sudo systemctl start pushgateway
sudo systemctl status pushgateway
curl localhost:9091/metrics
```

Configure Prometheus to scrape a pushgateway

```
sudo vi prometheus.yml

- job_name: 'Pushgateway'
    honor_labels: true
    static_configs:
    - targets: ['localhost:9091']
```

```
sudo systemctl restart prometheus
```


#### Sending data to pushgateway

```
echo "value_of_pi 3.14" | curl --data-binary @- http://localhost:9091/metrics/job/my_job
```

```
cat << EOF | curl --data-binary @- http://localhost:9091/metrics/job/my_job/instance/my_instance
# TYPE temperature gauge
temperature{location="room1"} 31
temperature{location="room2"} 33
# TYPE my_metric gauge
# HELP my_metric An example.
my_metric 5
EOF
```

`curl localhost:9091/metrics`


### Recording rules

Provide a layer of control over data, allowing pre-calculate new metrics

https://prometheus.io/docs/prometheus/latest/configuration/recording_rules/

### Alertmanager

Responsible for handling alerts sent to it by clients. i.e. Prometheus server

- https://prometheus.io/docs/alerting/latest/alertmanager/
- https://github.com/prometheus/alertmanager
- https://prometheus.io/docs/alerting/latest/overview/
- https://prometheus.io/docs/alerting/latest/configuration/


#### Alertrmanager HA

- https://prometheus.io/docs/alerting/latest/alertmanager/#high-availability
- https://github.com/prometheus/alertmanager#high-availability


### Alerting rules

- https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/


### Managing alerts
- Routing
- Grouping 
- Inhibition - suppress an alert if another alert is already firing.
- Silencies - temporarily turn off notifications

## Resources

- "A blog on monitoring, scale and operational Sanity" https://www.robustperception.io/blog
- Documentation https://prometheus.io/docs/prometheus/
- Overview https://prometheus.io/docs/introduction/overview/
- Prometheus configuration: https://prometheus.io/docs/prometheus/latest/configuration/configuration/
- Example configuration file: https://github.com/prometheus/prometheus/blob/release-2.15/config/testdata/conf.good.yml