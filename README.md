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

| Cluster | Tags |
| --- | --- |
| Cluster A | `runner-a, rnd-runner` |
| Cluster B | `runner-b` |

Each [job](https://docs.gitlab.com/ee/ci/jobs/) has to run on the same cluster
as the previous job because of the
[caches](https://docs.gitlab.com/ee/ci/caching/) configured per cluster.

In a scenario where `Cluster A` faints, our developers will be forced to change
their [.gitlab-ci.yml](https://docs.gitlab.com/ee/ci/yaml/?query=.gitlab-ci)
files to use `runner-b` instead of `rnd-runner`.

With our `μserivce`, your teams will not even know that `Cluster A` had any
problems and they can continue working as usual.

## Commands

### Sidecar

The sidecar will check the status about the `GitLab Runner` and export a
[liveness health check][liveness health check].

[liveness health check]: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/#define-a-liveness-command

### Checker

The checker has 2 jobs:

1. Status check:
    - Checking the `liveness health check` of each available `GitLab Runner`
installed on the same cluster
    - Exporting the statuses (`Online / Total`) to the configured
`GitLab` project
1. Executer:
    - Read the statuses of both clusters
    - Decide if the current cluster is stronger (has more
`GitLab Runner` available)
    - Stronger:
        - Update the tag list of the `GitLab Runners` to the desired main tag
    - Weaker:
        - Remove the main tag from the tag list
