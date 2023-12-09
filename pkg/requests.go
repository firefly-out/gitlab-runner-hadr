package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"os"
	"time"
)

// SidecarExecutor runs the sidecar for achieving the following goals:
//   - Getting the status of the pods Runner
//   - Exporting status using prometheus metrics
func SidecarExecutor(gitlabBaseUrl, gitlabGroupId string, statusInterval int) {
	var fullApiUrl = fmt.Sprintf("%s/groups/%s/runners", gitlabBaseUrl, gitlabGroupId)
	var runnerName = os.Getenv("HOSTNAME")
	prometheus.Register(TotalStatusRequests)
	prometheus.Register(TotalRunnersAvailableCount)
	prometheus.Register(RunnerOnlineStatus)

	fmt.Printf("Starting the GitLab Runner HADR Sidecar to check for Runner: %s\n"+
		"Using \"%s\" to check for the runners status\n", runnerName, fullApiUrl)
	for {
		allRunnersAvailable, err := getAllRunnerStatutes(fullApiUrl)
		if err != nil {
			fmt.Println(err)
			TotalStatusRequests.With(prometheus.Labels{"status": "failed"}).Inc()
		} else {
			TotalStatusRequests.With(prometheus.Labels{"status": "success"}).Inc()
		}
		TotalRunnersAvailableCount.Set(float64(len(allRunnersAvailable)))

		status, err := fetchCurrentRunnerStatus(runnerName, allRunnersAvailable)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(status)

		if status.Online {
			RunnerOnlineStatus.Set(1)
		} else {
			RunnerOnlineStatus.Set(0)
		}

		duration := time.Duration(statusInterval) * time.Second // Calculate duration using the VARIABLE value
		time.Sleep(duration)
	}
}

// getAllRunnerStatutes tries to fetch details about all available runners from the given API url.
// This method returns an array of the RunnerStatus struct.
func getAllRunnerStatutes(url string) (runners []RunnerStatus, err error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("API is not available at the moment")
	}

	// Decode the JSON response into an array of RunnerStatus structs
	err = json.NewDecoder(res.Body).Decode(&runners)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %v", err)
	}

	return runners, nil
}

// fetchCurrentRunnerStatus to iterate all the received RunnerStatus array
// to find the needed runner and return its status.
func fetchCurrentRunnerStatus(runnerName string, runners []RunnerStatus) (status RunnerStatus, err error) {
	for _, currentRunner := range runners {
		if currentRunner.Name == runnerName {
			return currentRunner, nil
		}
	}
	return RunnerStatus{}, fmt.Errorf("Runner %s was not found\n", runnerName)
}
