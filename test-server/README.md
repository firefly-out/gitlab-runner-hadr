# GitLab Instance Mock

This server will mock the GitLabs API.

## GET /groups/33/runners

Executing a `/GET` request to `/groups/33/runners` will result in the following
json:

```json
[
  {
    "active":true,
    "paused":false,
    "description":"test-1-20150125",
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
    "description":"test-2-20150125",
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
    "description":"test-3-20150125",
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

## How to test locally?

1st there is a need to run this server, can be done by `go run main.go`.

Then, when the server is running, there is a need to execute the sidecar command
by:

1. Building the [Dockerfile](/Dockerfile)
2. Running it locally using the following terminal command to mimic a `k8s` pod

```shell
docker run -e HOSTNAME='runner-1' -e POD_NAME='runner-1'
-e POD_NAMESPACE='runners' -e POD_IP='127.0.0.1' -p 8081:8080
<image_name>:<tag> sidecar -i 33 -u http://<machines_ip>:8080
```

## Prometheus Metrics

The sidecar service exports the following prometheus metrics:

```prometheus
# HELP gitlab_runner_hadr_sidecar_online_status Indicates if this runners statusis online using the information received from gitlab.
# TYPE gitlab_runner_hadr_sidecar_online_status gauge
gitlab_runner_hadr_sidecar_online_status{group_id="33",pod_ip="127.0.0.1",pod_name="runner-1",pod_namespace="runners"} 1

# HELP gitlab_runner_hadr_sidecar_total_runners_available A metric that indicates how many runners were returned from the /get request to the api.
# TYPE gitlab_runner_hadr_sidecar_total_runners_available gauge
gitlab_runner_hadr_sidecar_total_runners_available{group_id="33",pod_ip="127.0.0.1",pod_name="runner-1",pod_namespace="runners"} 3

# HELP gitlab_runner_hadr_sidecar_total_status_requests Counting how many /get requests where executed by the sidecar.
# TYPE gitlab_runner_hadr_sidecar_total_status_requests gauge
gitlab_runner_hadr_sidecar_total_status_requests{group_id="33",pod_ip="127.0.0.1",pod_name="runner-1",pod_namespace="runners",status="success"} 58

# HELP gitlab_runner_hadr_uptime_seconds The uptime of the gitlab runner service in seconds.
# TYPE gitlab_runner_hadr_uptime_seconds gauge
gitlab_runner_hadr_uptime_seconds{pod_ip="127.0.0.1",pod_name="runner-1",pod_namespace="runners"} 285.062274212
```
