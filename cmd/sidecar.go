package cmd

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"os"
	"runner-hadr/metrics"
	"runner-hadr/pkg"
	"time"
)

// Sidecar runs the sidecar for achieving the following goals:
//   - Getting the status of the pods Runner
//   - Exporting status using prometheus metrics
func Sidecar(gitlabBaseUrl, gitlabGroupId, token, perPage string, statusInterval int) {
	var fullApiUrl = fmt.Sprintf("%s/groups/%s/runners?per_page=%s", gitlabBaseUrl, gitlabGroupId, perPage)
	var runnerName = os.Getenv("HOSTNAME")
	var duration = time.Duration(statusInterval) * time.Second // Calculate duration using the VARIABLE value

	prometheus.Register(metrics.TotalStatusRequests)
	prometheus.Register(metrics.TotalRunnersAvailableCount)
	prometheus.Register(metrics.RunnerOnlineStatus)

	fmt.Printf("Starting the GitLab Runner HADR Sidecar to check for Runner: %s\n"+
		"Using \"%s\" to check for the runners status\n", runnerName, fullApiUrl)
	for {
		allRunnersAvailable, err := pkg.GetAllRunnerStatutes(fullApiUrl, gitlabGroupId, token)
		if err != nil {
			fmt.Println(err)
		}
		metrics.ChangeTotalRunnersAvailableCount(float64(len(allRunnersAvailable)), gitlabGroupId)

		status, err := pkg.FetchCurrentRunnerStatus(runnerName, allRunnersAvailable)
		if err != nil {
			fmt.Println(err)
		}
		pkg.CheckRunnerStatus(status, gitlabGroupId)

		time.Sleep(duration)
	}
}
