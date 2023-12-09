package pkg

import (
	"fmt"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	Url         string       = "https://beaverr-io-gitlab.com"
	RunnerAName string       = "Runner-1"
	RunnerA     RunnerStatus = RunnerStatus{Active: true,
		Paused:      false,
		Description: "test-1-20150125",
		Id:          1,
		IpAddress:   "271.15.10.1",
		IsShared:    false,
		RunnerType:  "project_type",
		Name:        RunnerAName,
		Online:      true,
		Status:      "online"}
	RunnerBName string       = "Runner-2"
	RunnerB     RunnerStatus = RunnerStatus{Active: true,
		Paused:      false,
		Description: "test-2-20150125",
		Id:          2,
		IpAddress:   "271.15.10.3",
		IsShared:    false,
		RunnerType:  "project_type",
		Name:        RunnerBName,
		Online:      false,
		Status:      "offline"}
	RunnerCName string = "Runner-3"
	GroupId     string = "33"

	Runners = []RunnerStatus{RunnerA, RunnerB}
)

func Test_getAllRunnerStatutes(t *testing.T) {
	t.Run("fetching the runners status from the GitLab API", func(t *testing.T) {
		// Arrange
		defer gock.Off()
		gock.New(Url).
			Get("/api").
			Reply(200).
			JSON(Runners)

		// Act
		body, err := getAllRunnerStatutes(Url+"/api", GroupId)

		// Assert
		assert.Equal(t, 2, len(body))
		assert.Equal(t, RunnerA.Id, body[0].Id)
		assert.Equal(t, RunnerB.Id, body[1].Id)
		assert.Nil(t, err)
	})

	t.Run("unauthorized error was received from the GET request", func(t *testing.T) {
		// Arrange
		defer gock.Off()
		gock.New(Url).
			Get("/api").
			Reply(401).
			JSON(map[string]string{"error": "Unauthorized"})

		// Act
		body, err := getAllRunnerStatutes(Url+"/api", GroupId)

		// Assert
		assert.Nil(t, body)
		assert.Equal(t, fmt.Errorf("API is not available at the moment"), err)
	})

	t.Run("incorrect json returned from the API", func(t *testing.T) {
		// Arrange
		defer gock.Off()
		gock.New(Url).
			Get("/api").
			Reply(200).
			JSON(map[string]string{"foo": "bar"})

		// Act
		body, err := getAllRunnerStatutes(Url+"/api", GroupId)

		// Assert
		assert.Nil(t, body)
		assert.ErrorContains(t, err, "failed to decode JSON response")
	})

	t.Run("request to an empty url", func(t *testing.T) {
		// Act
		body, err := getAllRunnerStatutes(""+"/api", GroupId)

		// Assert
		assert.Nil(t, body)
		assert.ErrorContains(t, err, "unsupported protocol scheme")
	})
}

func Test_fetchCurrentRunnerStatus(t *testing.T) {
	t.Run("finding the needed runner", func(t *testing.T) {
		// Act
		status, err := fetchCurrentRunnerStatus(RunnerAName, Runners)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, RunnerA, status)
	})

	t.Run("runner does not exist", func(t *testing.T) {
		// Act
		status, err := fetchCurrentRunnerStatus(RunnerCName, Runners)

		// Assert
		assert.NotNil(t, status)
		assert.ErrorContains(t, err, fmt.Sprintf("Runner %s was not found", RunnerCName))
	})
}

func Test_checkRunnerStatus(t *testing.T) {
	t.Run("should return true is Online is true", func(t *testing.T) {
		// Act
		status := checkRunnerStatus(Runners[0], GroupId)

		// Assert
		assert.True(t, status)
	})

	t.Run("should return false is Online is false", func(t *testing.T) {
		// Act
		status := checkRunnerStatus(Runners[1], GroupId)

		// Assert
		assert.False(t, status)
	})
}
