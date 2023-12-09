package main

import (
	"encoding/json"
	"fmt"
	"log"
	http "net/http"
	"runner-hadr/pkg"
)

var (
	GetRunnersApiUrl = "/runners/33/runners"
	Runners          = []pkg.RunnerStatus{{
		Active:      true,
		Paused:      false,
		Description: "test-1-20150125",
		Id:          1,
		IpAddress:   "271.15.10.1",
		IsShared:    false,
		RunnerType:  "project_type",
		Name:        "",
		Online:      true,
		Status:      "online"}, {
		Active:      true,
		Paused:      false,
		Description: "test-2-20150125",
		Id:          2,
		IpAddress:   "271.15.10.3",
		IsShared:    false,
		RunnerType:  "project_type",
		Name:        "",
		Online:      true,
		Status:      "online"},
	}
)

func main() {
	http.HandleFunc(GetRunnersApiUrl, handleGet)

	log.Println("Go!")
	http.ListenAndServe(":8080", nil)
}

func handleGet(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		j, err := json.Marshal(Runners)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = w.Write(j)
		if err != nil {
			fmt.Println(err)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "I can't do that.")
	}
}