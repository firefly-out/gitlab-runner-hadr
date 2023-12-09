package pkg

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

var (
	startTime = time.Now()
	uptime    = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "gitlab_runner_hadr_uptime_seconds",
			Help: "The uptime of the gitlab runner service in seconds.",
		},
	)
	TotalRunnersAvailableCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "gitlab_runner_hadr_sidecar_total_runners_available",
			Help: "A metric that indicates how many runners were returned from the /get request to the api.",
		},
	)
	RunnerOnlineStatus = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "gitlab_runner_hadr_sidecar_online_status",
			Help: "Indicates if this runners status is online using the information received from gitlab.",
		},
	)
	TotalStatusRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gitlab_runner_hadr_sidecar_total_status_requests",
			Help: "Counting how many /get requests where executed by the sidecar.",
		},
		[]string{"status"},
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
			uptime.Set(time.Since(startTime).Seconds())
			time.Sleep(time.Second)
		}
	}()

	// Expose the registered metrics via an HTTP endpoint
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
