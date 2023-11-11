package main

import (
	"fmt"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	Url     string       = "https://herrmit-io-gitlab.com"
	RunnerA RunnerStatus = RunnerStatus{Active: true,
		Paused:      false,
		Description: "test-1-20150125",
		Id:          1,
		IpAddress:   "271.15.10.1",
		IsShared:    false,
		RunnerType:  "project_type",
		Name:        "",
		Online:      true,
		Status:      "online"}
	RunnerB RunnerStatus = RunnerStatus{Active: true,
		Paused:      false,
		Description: "test-2-20150125",
		Id:          2,
		IpAddress:   "271.15.10.3",
		IsShared:    false,
		RunnerType:  "project_type",
		Name:        "",
		Online:      true,
		Status:      "online"}

	Runners = []RunnerStatus{RunnerA, RunnerB}
)

func Test_GetRunnerStatus(t *testing.T) {
	t.Run("fetching the runners status from the GitLab API", func(t *testing.T) {
		// Arrange
		defer gock.Off()
		gock.New(Url).
			Get("/api").
			Reply(200).
			JSON(Runners)

		// Act
		body, err := GetRunnerStatutes(Url + "/api")

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
		body, err := GetRunnerStatutes(Url + "/api")

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
		body, err := GetRunnerStatutes(Url + "/api")

		// Assert
		assert.Nil(t, body)
		assert.ErrorContains(t, err, "failed to decode JSON response")
	})

	t.Run("request to an empty url", func(t *testing.T) {
		// Act
		body, err := GetRunnerStatutes("" + "/api")

		// Assert
		assert.Nil(t, body)
		assert.ErrorContains(t, err, "unsupported protocol scheme")
	})
}
