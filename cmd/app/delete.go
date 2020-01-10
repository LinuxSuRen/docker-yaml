package app

import (
	"context"
	"github.com/docker/docker/client"
	"github.com/linuxsuren/docker-yaml/cmd/common"
	"github.com/linuxsuren/docker-yaml/pkg"
	"github.com/spf13/cobra"
)

// NewDeleteCommand create app command
func NewDeleteCommand(commonOpts *common.Options) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "delete",
		Short: "delete",
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

					if err = df.RemoveContainer(); err != nil {
						break
					}
				}
			}
			return 
		},
	}
	return
}
