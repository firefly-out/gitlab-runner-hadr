package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"runner-hadr/pkg"
)

var (
	GitLabUrl           string
	GitLabGroupId       string
	StatusCheckInterval int

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
			pkg.SidecarExecutor(GitLabUrl, GitLabGroupId, StatusCheckInterval)
		},
	}
)

func Execute() {
	sidecarCmd.Flags().StringVarP(&GitLabUrl, "gitlab-url", "u", "", "The GitLabs url to check if the runner is connected to it")
	sidecarCmd.Flags().StringVarP(&GitLabGroupId, "gitlab-group-id", "i", "", "The group ID the runner is installed on")
	sidecarCmd.Flags().IntVarP(&StatusCheckInterval, "status-check-interval", "s", 5, "Interval for checking the status of the runner in seconds")

	// Add sidecar and decider commands to the root command
	rootCmd.AddCommand(sidecarCmd)
	sidecarCmd.MarkFlagRequired("gitlab-url")      // Mark the "gitlab-url" flag as required
	sidecarCmd.MarkFlagRequired("gitlab-group-id") // Mark the "gitlab-group-id" flag as required

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
