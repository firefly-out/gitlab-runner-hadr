package main

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
