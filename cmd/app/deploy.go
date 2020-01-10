package app

import (
	"context"
	"github.com/docker/docker/client"
	"github.com/linuxsuren/docker-yaml/cmd/common"
	"github.com/linuxsuren/docker-yaml/pkg"
	"github.com/spf13/cobra"
)

// NewDeployCommand deploy app command
func NewDeployCommand(commonOpts *common.Options) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "deploy",
		Short: "deploy",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var cli *client.Client
			cli, err = client.NewEnvClient()
			if err != nil {
				return
			}

			var apps []pkg.Application
			if apps, err = pkg.GetApplications("docker.yaml"); err == nil {
				for _, app := range apps {
					df := pkg.DockerDeploy{
						App:     app,
						Context: context.Background(),
						Client:  cli,
					}

					if err = df.DeployImage(); err != nil {
						break
					}
				}
			}
			return
		},
	}
	return
}
