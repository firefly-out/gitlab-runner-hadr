# GitLab Instance Mock

This server will mock the GitLabs API.

## GET /runners/33/runners

Executing a `/GET` request to `/runners/33/runners` will result in the following json:

```json
[
  {
    "active":true,
    "paused":false,
    "description":"test-1-20150125",
    "id":1,
    "ip_address":"271.15.10.1",
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
    "ip_address":"271.15.10.3",
    "is_shared":false,
    "runner_type":"project_type",
    "name":"",
    "online":true,
    "status":"online"
  }
]
```
