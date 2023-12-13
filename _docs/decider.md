# Decider

The `decider` checks the statuses of all runners deployed on the cluster via
reading the metrics the `sidecar` exports.
Then writes it to a file in the `GitLab` instance for sharing the clusters
status with other `deciders` from other clusters.

The `decider` also checks which cluster is `"stronger"` (having more runners
available) to change the runners `tags` so our clients won't have to change
their `.gitlab-ci.yml` files at all and run their pipelines on a different
cluster without feeling a thing.

## Pre-Conditions

This service uses the
[k8s.io/client-go/kubernetes](https://github.com/kubernetes/client-go) which
allows it to execute API calls to the cluster its deployed on.
Therefor, there is a need to deploy a `Role` alongside a `RoleBinding`
for listing the pods:

```yml
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  namespace: gitlab-runners
  name: pod-reader
rules:
  - apiGroups: [""]
    resources: ["pods"]
    verbs: ["get", "list", "watch"]
```

```yml
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: read-pods
  namespace: gitlab-runners
subjects:
  - kind: ServiceAccount
    name: default
    namespace: gitlab-runners
roleRef:
  kind: Role
  name: pod-reader
  apiGroup: rbac.authorization.k8s.io
```

## Deployment

The `decider` should be deployed on the runners cluster:

```yml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: decider
spec:
  replicas: 1
  selector:
    matchLabels:
      app: decider
  template:
    metadata:
      labels:
        app: decider
    spec:
      containers:
        - name: sidecar
          image: zigelboimmisha/gitlab-runner-hadr:main-60081c3
          imagePullPolicy: Always
          args:
            - "decider"
```
