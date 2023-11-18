# Gitlab Runner HADR

This project will help you make sure there is at-least one
[GitLab Runners](https://docs.gitlab.com/runner/) group available for your
[GitLab](https://docs.gitlab.com/ee/user/project/working_with_projects.html)
projects. By checking the status of each `GitLab Runner` via the `GitLabs`
server API, our μserivce will be able to switch the `GitLab Runners` around to
have an Active-Active scenario 24/7.

## Real Life Scenario

An organization has `GitLab Runners` installed in 2 separated clusters
registered with different `Tags`:

| Cluster   | Tags                      |
|-----------|---------------------------|
| Cluster A | `runner-a`, `main-runner` |
| Cluster B | `runner-b`                |

Each [job](https://docs.gitlab.com/ee/ci/jobs/) has to run on the same cluster
as the previous job because of the
[caches](https://docs.gitlab.com/ee/ci/caching/) configured per cluster.

In a scenario where `Cluster A` faints, our developers will be forced to change
their [.gitlab-ci.yml](https://docs.gitlab.com/ee/ci/yaml/?query=.gitlab-ci)
files to use `runner-b` instead of `rnd-runner`.

With our `μserivce`, your teams will not even know that `Cluster A` had any
problems, and they can continue working as usual.

## Sidecar

![sidecar check](/docs/visio/1.%20Sidecar%20check.PNG)

The sidecar will check the status about the `GitLab Runner` and export a
[liveness health check][liveness health check].

[liveness health check]: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/#define-a-liveness-command

![sidecar check](/docs/visio/1.%20Sidecar%20check%20explained.PNG)

## Decider

The decider has 2 jobs:

### Status Check

![sidecar check](/docs/visio/2.%20Decider%20check.PNG)

- Checking the `liveness health check` of each available `GitLab Runner`
installed on the same cluster

- Exporting the statuses (`Online / Total`) to the configured
`GitLab` project

![sidecar check](/docs/visio/2.%20Decider%20check%20explained.PNG)

### Update Runners

![sidecar check](/docs/visio/3.%20Deciders%20action.PNG)

- Read the statuses of both clusters

- Decide if the current cluster is stronger (has more
  `GitLab Runner` available)

#### Stronger

If stringer, it will update the tag list of the `GitLab Runners`
to the configured desired main cluster tag so users will not have to change
their `.gitlab-ci.yml` files.

![sidecar check](/docs/visio/3.%20Deciders%20action%20-%20win.PNG)

#### Weaker

If weaker, it will remove the configured desired main cluster tag from
the tag list so users will stop using those runners without updating their
`.gitlab-ci.yml` file.

![sidecar check](/docs/visio/3.%20Deciders%20action%20-%20lose.PNG)
