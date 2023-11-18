package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
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
