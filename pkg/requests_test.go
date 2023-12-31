package pkg

import (
	"fmt"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

var (
	Url         string       = "https://beaverr-io-gitlab.com"
	RunnerAName string       = "Runner-1"
	RunnerA     RunnerStatus = RunnerStatus{Active: true,
		Paused:      false,
		Description: RunnerAName,
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
		Description: RunnerBName,
		Id:          2,
		IpAddress:   "271.15.10.3",
		IsShared:    false,
		RunnerType:  "project_type",
		Name:        RunnerBName,
		Online:      false,
		Status:      "offline"}
	RunnerCName string = "Runner-3"
	GroupId     string = "33"
	Token       string = "dsa543ads"
	Runners            = []RunnerStatus{RunnerA, RunnerB}
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
		body, err := GetAllRunnerStatutes(Url+"/api", GroupId, Token)

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
		body, err := GetAllRunnerStatutes(Url+"/api", GroupId, Token)

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
		body, err := GetAllRunnerStatutes(Url+"/api", GroupId, Token)

		// Assert
		assert.Nil(t, body)
		assert.ErrorContains(t, err, "failed to decode JSON response")
	})

	t.Run("request to an empty url", func(t *testing.T) {
		// Act
		body, err := GetAllRunnerStatutes(""+"/api", GroupId, Token)

		// Assert
		assert.Nil(t, body)
		assert.ErrorContains(t, err, "unsupported protocol scheme")
	})
}

func Test_fetchCurrentRunnerStatus(t *testing.T) {
	t.Run("finding the needed runner", func(t *testing.T) {
		// Act
		status, err := FetchCurrentRunnerStatus(RunnerAName, Runners)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, RunnerA, status)
	})

	t.Run("runner does not exist", func(t *testing.T) {
		// Act
		status, err := FetchCurrentRunnerStatus(RunnerCName, Runners)

		// Assert
		assert.NotNil(t, status)
		assert.ErrorContains(t, err, fmt.Sprintf("runner %s was not found", RunnerCName))
	})
}

func Test_checkRunnerStatus(t *testing.T) {
	t.Run("should return true is Online is true", func(t *testing.T) {
		// Act
		status := CheckRunnerStatus(Runners[0], GroupId)

		// Assert
		assert.True(t, status)
	})

	t.Run("should return false is Online is false", func(t *testing.T) {
		// Act
		status := CheckRunnerStatus(Runners[1], GroupId)

		// Assert
		assert.False(t, status)
	})
}
