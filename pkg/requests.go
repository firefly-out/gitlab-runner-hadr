package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runner-hadr/metrics"
)

var (
	httpClient = &http.Client{}
)

// GetAllRunnerStatutes tries to fetch details about all available runners from the given API url.
// This method returns an array of the RunnerStatus struct.
func GetAllRunnerStatutes(url, gitlabGroupId, token string) (runners []RunnerStatus, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", token)
	res, err := httpClient.Do(req)
	if err != nil {
		metrics.IncreaseTotalStatusRequests("failed", gitlabGroupId)
		return nil, err
	}
	defer res.Body.Close()
	metrics.IncreaseTotalStatusRequests("success", gitlabGroupId)

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

// FetchCurrentRunnerStatus to iterate all the received RunnerStatus array
// to find the needed runner and return its status.
func FetchCurrentRunnerStatus(runnerName string, runners []RunnerStatus) (status RunnerStatus, err error) {
	for _, currentRunner := range runners {
		if currentRunner.Description == runnerName {
			return currentRunner, nil
		}
	}
	return RunnerStatus{}, fmt.Errorf("runner %s was not found", runnerName)
}

// CheckRunnerStatus to return true if the runners status is online, false otherwise.
func CheckRunnerStatus(runner RunnerStatus, gitlabGroupId string) (runnersStatus bool) {
	if runner.Online {
		metrics.ChangeRunnerOnlineStatus(1, gitlabGroupId)
		runnersStatus = true
	} else {
		metrics.ChangeRunnerOnlineStatus(0, gitlabGroupId)
		runnersStatus = false
	}

	return runnersStatus
}
