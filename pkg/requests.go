package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpClient = &http.Client{}
)

// SidecarExecutor runs the sidecar for achieving the following goals:
//   - Getting the status of the pods Runner
//   - Exporting status using prometheus metrics
func SidecarExecutor(gitlabBaseUrl, gitlabGroupId, token, perPage string, statusInterval int) {
	var fullApiUrl = fmt.Sprintf("%s/groups/%s/runners?per_page=%s", gitlabBaseUrl, gitlabGroupId, perPage)
	var runnerName = os.Getenv("HOSTNAME")
	var duration = time.Duration(statusInterval) * time.Second // Calculate duration using the VARIABLE value

	prometheus.Register(TotalStatusRequests)
	prometheus.Register(TotalRunnersAvailableCount)
	prometheus.Register(RunnerOnlineStatus)

	fmt.Printf("Starting the GitLab Runner HADR Sidecar to check for Runner: %s\n"+
		"Using \"%s\" to check for the runners status\n", runnerName, fullApiUrl)
	for {
		allRunnersAvailable, err := getAllRunnerStatutes(fullApiUrl, gitlabGroupId, token)
		if err != nil {
			fmt.Println(err)
		}
		changeTotalRunnersAvailableCount(float64(len(allRunnersAvailable)), gitlabGroupId)

		status, err := fetchCurrentRunnerStatus(runnerName, allRunnersAvailable)
		if err != nil {
			fmt.Println(err)
		}
		checkRunnerStatus(status, gitlabGroupId)

		time.Sleep(duration)
	}
}

// getAllRunnerStatutes tries to fetch details about all available runners from the given API url.
// This method returns an array of the RunnerStatus struct.
func getAllRunnerStatutes(url, gitlabGroupId, token string) (runners []RunnerStatus, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", token)
	res, err := httpClient.Do(req)
	if err != nil {
		increaseTotalStatusRequests("failed", gitlabGroupId)
		return nil, err
	}
	defer res.Body.Close()
	increaseTotalStatusRequests("success", gitlabGroupId)

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

// checkRunnerStatus to return true if the runners status is online, false otherwise.
func checkRunnerStatus(runner RunnerStatus, gitlabGroupId string) (runnersStatus bool) {
	if runner.Online {
		changeRunnerOnlineStatus(1, gitlabGroupId)
		runnersStatus = true
	} else {
		changeRunnerOnlineStatus(0, gitlabGroupId)
		runnersStatus = false
	}

	return runnersStatus
}
