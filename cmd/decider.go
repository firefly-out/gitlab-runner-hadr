package cmd

import (
	"runner-hadr/pkg"
	"time"
)

// Decider runs the decider
func Decider(namespace, deployment string) (err error) {
	client, err := pkg.CreateK8sClient()
	if err != nil {
		return err
	}

	for {
		pkg.ListPods(client, namespace, deployment)
		time.Sleep(time.Second * 10)
	}
}
