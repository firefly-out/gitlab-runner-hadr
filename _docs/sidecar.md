# Sidecar

The `sidecar` service checks the status of the `GitLab Runners` and export it
as prometheus metrics for the `decider` to fetch and use.

## Getting the runners status

This is done by sending a `/GET` request to the `GitLabs` API
(`/groups/33/runners`) and checking if the status is online from the following
json response:

```json
[
  {
    "active":true,
    "paused":false,
    "description":"gitlab-runner-76cc4c686c-ds92l",
    "id":1,
    "ip_address":"211.15.10.1",
    "is_shared":false,
    "runner_type":"project_type",
    "name":"",
    "online":true,
    "status":"online"
  },
  {
    "active":true,
    "paused":false,
    "description":"gitlab-runner-76cc4c686c-2sd5f",
    "id":2,
    "ip_address":"211.15.10.3",
    "is_shared":false,
    "runner_type":"project_type",
    "name":"",
    "online":true,
    "status":"online"
  },
  {
    "active":true,
    "paused":false,
    "description":"gitlab-runner-76cc4c686c-2wjrp",
    "id":3,
    "ip_address":"211.231.10.3",
    "is_shared":false,
    "runner_type":"project_type",
    "name":"",
    "online":false,
    "status":"offline"
  }
]
```

Then using the `HOSTNAME` variable that `k8s` injects into every pod, the
`sidecar` compares its `hostname` to the `description` key to find its runner.

## Deployment Example

As the name suggests, there is a need to run the `sidecar` as a `sidecar` to
your `GitLab CI Runners` by adding the following to your runners `Deployment`:

```yml
containers:
  - name: sidecar
    image: zigelboimmisha/gitlab-runner-hadr:main-60081c3
    ports:
      - containerPort: 8080
    args:
      - "sidecar"
      - "-i"
      - "33"
      - "-u"
      - "http://172.28.192.1:80"
      - "-t"
      - "dasdasdasd"
```

## Prometheus Metrics

### Status

The `sidecar` exports the status of the runner by exporting:

```prometheus
# HELP gitlab_runner_hadr_sidecar_online_status Indicates if this runners statusis online using the information received from gitlab.
# TYPE gitlab_runner_hadr_sidecar_online_status gauge
gitlab_runner_hadr_sidecar_online_status{group_id="33",pod_ip="127.0.0.1",pod_name="runner-1",pod_namespace="runners"} 1
```

### Total Runners Available

It also exports how many runners are available to the `GitLab Group` by:

```prometheus
# HELP gitlab_runner_hadr_sidecar_total_runners_available A metric that indicates how many runners were returned from the /get request to the api.
# TYPE gitlab_runner_hadr_sidecar_total_runners_available gauge
gitlab_runner_hadr_sidecar_total_runners_available{group_id="33",pod_ip="127.0.0.1",pod_name="runner-1",pod_namespace="runners"} 3
```

### Monitoring the Sidecar

To monitor the `sidecar`, it also export different statuses:

#### Uptime

```prometheus
# HELP gitlab_runner_hadr_uptime_seconds The uptime of the gitlab runner service in seconds.
# TYPE gitlab_runner_hadr_uptime_seconds gauge
gitlab_runner_hadr_uptime_seconds{pod_ip="127.0.0.1",pod_name="runner-1",pod_namespace="runners"} 285.062274212
```

#### Total Requests

```prometheus
# HELP gitlab_runner_hadr_sidecar_total_status_requests Counting how many /get requests where executed by the sidecar.
# TYPE gitlab_runner_hadr_sidecar_total_status_requests gauge
gitlab_runner_hadr_sidecar_total_status_requests{group_id="33",pod_ip="127.0.0.1",pod_name="runner-1",pod_namespace="runners",status="success"} 58
```
