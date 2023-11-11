package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// RunnerStatus represents the json returned from gitlab-server/api/v4/runners.
type RunnerStatus struct {
	Active      bool   `json:"active"`
	Paused      bool   `json:"paused"`
	Description string `json:"description"`
	Id          int64  `json:"id"`
	IpAddress   string `json:"ip_address"`
	IsShared    bool   `json:"is_shared"`
	RunnerType  string `json:"runner_type"`
	Name        string `json:"name"`
	Online      bool   `json:"online"`
	Status      string `json:"status"`
}

// GetRunnerStatutes tries to fetch details about all available runners from the given API url.
// This method returns an array of the RunnerStatus struct.
func GetRunnerStatutes(url string) (runners []RunnerStatus, err error) {
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
