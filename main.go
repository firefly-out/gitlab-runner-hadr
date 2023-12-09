package main

import (
	"fmt"
	"runner-hadr/pkg"
)

var (
	BaseUrl          = "http://localhost:8080"
	GetRunnersApiUrl = "/runners/33/runners"
)

func main() {
	body, err := pkg.GetAllRunnerStatutes(BaseUrl + GetRunnersApiUrl)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(body)
}
