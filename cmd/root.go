package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"runner-hadr/metrics"
)

var (
	GitLabUrl           string // The url of our GitLab instance to send our API requests to
	GitLabGroupId       string // The group id to check its runners
	GitLabToken         string // Token to access the api
	RunnersPerPage      string // How many runners to show, its needed as the default runner number is 20
	RunnersDeployment   string // The name of our runners deployment
	RunnersNamespace    string // Where are the runners deployed at
	StatusCheckInterval int    // How often do we want to check on our runners

	rootCmd = &cobra.Command{
		Use:   "gitlab-runner-hadr",
		Short: "gitlab-runner-hadr is a cli that allows your GitLab Runners to have HADR",
		Long: `gitlab-runner-hadr is a cli that allows your GitLab Runners to have HADR
                by implementing 2 micro-services that will run along side your
                GitLab Runners to make sure your Runners are available without
                the need to change your .gitlab-ci.yml files.
                Complete documentation is available at https://github.com/firefly-out/gitlab-runner-hadr`,
	}

	// Command for the Sidecar
	sidecarCmd = &cobra.Command{
		Use:   "sidecar",
		Short: "Executes the sidecar mirco-service that will export the status of your GitLab Runner",
		Example: `	cli sidecar -i 33 -u http://localhost:8080
	Increasing the status check interval to 10 seconds:
	cli sidecar -i 33 -u http://localhost:8080 -s 10`,
		Run: func(cmd *cobra.Command, args []string) {
			Sidecar(GitLabUrl, GitLabGroupId, GitLabToken, RunnersPerPage, StatusCheckInterval)
		},
	}

	// Command for the Decider
	deciderCmd = &cobra.Command{
		Use:     "decider",
		Short:   "Executes the decider mirco-service that will check the status of your GitLab Runner",
		Example: `	cli decider`,
		Run: func(cmd *cobra.Command, args []string) {
			Decider(RunnersNamespace, RunnersDeployment)
		},
	}
)

func Execute() {
	sidecarCmd.Flags().StringVarP(&GitLabUrl, "gitlab-url", "u", "", "The GitLabs url to check if the runner is connected to it")
	sidecarCmd.Flags().StringVarP(&GitLabGroupId, "gitlab-group-id", "i", "", "The group ID the runner is installed on")
	sidecarCmd.Flags().StringVarP(&GitLabToken, "gitlab-token", "t", "", "Access token to read the API")
	sidecarCmd.Flags().StringVarP(&RunnersPerPage, "runners-per-list", "r", "1000", "How many runners to request from the API")
	sidecarCmd.Flags().IntVarP(&StatusCheckInterval, "status-check-interval", "s", 5, "Interval for checking the status of the runner in seconds")
	deciderCmd.Flags().StringVarP(&RunnersNamespace, "namespace", "n", "gitlab-runners", "The namespace your gitlab runners are deployed on")
	deciderCmd.Flags().StringVarP(&RunnersDeployment, "deployment", "d", "gitlab-runner", "The name of your gitlab runners deployment")

	// Add sidecar to the root command
	rootCmd.AddCommand(sidecarCmd)
	sidecarCmd.MarkFlagRequired("gitlab-url")
	sidecarCmd.MarkFlagRequired("gitlab-group-id")
	sidecarCmd.MarkFlagRequired("gitlab-token")

	// Adding decider command to the root command
	rootCmd.AddCommand(deciderCmd)
	deciderCmd.MarkFlagRequired("namespace")

	// Start serving the prometheus metrics
	go metrics.RunUptimeMetrics()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
