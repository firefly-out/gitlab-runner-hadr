package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"time"
)

var (
	PodName      = os.Getenv("POD_NAME")
	PodIp        = os.Getenv("POD_IP")
	PodNamespace = os.Getenv("POD_NAMESPACE")

	startTime = time.Now()
	uptime    = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gitlab_runner_hadr_uptime_seconds",
			Help: "The uptime of the gitlab runner service in seconds.",
		},
		[]string{"pod_name", "pod_ip", "pod_namespace"},
	)
	TotalRunnersAvailableCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gitlab_runner_hadr_sidecar_total_runners_available",
			Help: "A metric that indicates how many runners were returned from the /get request to the api.",
		},
		[]string{"pod_name", "pod_ip", "pod_namespace", "group_id"},
	)
	RunnerOnlineStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gitlab_runner_hadr_sidecar_online_status",
			Help: "Indicates if this runners status is online using the information received from gitlab.",
		},
		[]string{"pod_name", "pod_ip", "pod_namespace", "group_id"},
	)
	TotalStatusRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gitlab_runner_hadr_sidecar_total_status_requests",
			Help: "Counting how many /get requests where executed by the sidecar.",
		},
		[]string{"status", "pod_name", "pod_ip", "pod_namespace", "group_id"},
	)
)

func init() {
	prometheus.MustRegister(uptime)
}

// RunUptimeMetrics exports the prometheus metrics to the world.
func RunUptimeMetrics() {
	go func() {
		for {
			// Update the uptime metric every second
			uptime.With(prometheus.Labels{"pod_name": PodName, "pod_ip": PodIp, "pod_namespace": PodNamespace}).Set(time.Since(startTime).Seconds())
			time.Sleep(time.Second)
		}
	}()

	// Expose the registered metrics via an HTTP endpoint
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}

// IncreaseTotalStatusRequests to increase the gitlab_runner_hadr_sidecar_total_status_requests metric.
func IncreaseTotalStatusRequests(status, gitlabGroupId string) {
	TotalStatusRequests.With(
		prometheus.Labels{
			"status":        status,
			"pod_name":      PodName,
			"pod_ip":        PodIp,
			"pod_namespace": PodNamespace,
			"group_id":      gitlabGroupId}).Inc()
}

// ChangeRunnerOnlineStatus to change the gitlab_runner_hadr_sidecar_online_status metric.
func ChangeRunnerOnlineStatus(runnerStatus float64, gitlabGroupId string) {
	RunnerOnlineStatus.With(
		prometheus.Labels{
			"pod_name":      PodName,
			"pod_ip":        PodIp,
			"pod_namespace": PodNamespace,
			"group_id":      gitlabGroupId}).Set(runnerStatus)
}

// ChangeTotalRunnersAvailableCount to change the gitlab_runner_hadr_sidecar_total_runners_available metric.
func ChangeTotalRunnersAvailableCount(runnersReturned float64, gitlabGroupId string) {
	TotalRunnersAvailableCount.With(prometheus.Labels{
		"pod_name":      PodName,
		"pod_ip":        PodIp,
		"pod_namespace": PodNamespace,
		"group_id":      gitlabGroupId}).Set(runnersReturned)
}
