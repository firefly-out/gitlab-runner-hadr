package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// GetAllRunnerStatutes tries to fetch details about all available runners from the given API url.
// This method returns an array of the RunnerStatus struct.
func GetAllRunnerStatutes(url string) (runners []RunnerStatus, err error) {
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

// SidecarExecutor runs the sidecar for achieving the following goals:
//   - Getting the status of the pods Runner
//   - Exporting status using prometheus metrics
func SidecarExecutor(gitlabBaseUrl, gitlabGroupId string, statusInterval int) {
	var fullApiUrl = fmt.Sprintf("%s/groups/%s/runners", gitlabBaseUrl, gitlabGroupId)
	fmt.Printf("Starting the GitLab Runner HADR Sidecar to check for Runner: %s\n"+
		"Using \"%s\" to check for the runners status\n", "test", fullApiUrl)

	for {
		body, err := GetAllRunnerStatutes(fullApiUrl)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(body)

		duration := time.Duration(statusInterval) * time.Second // Calculate duration using the VARIABLE value
		time.Sleep(duration)
	}
}
