package pkg

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"time"
)

type DockerDeploy struct {
	Client  *client.Client
	App     Application
	Context context.Context
}

func (d *DockerDeploy) DeployImage() (err error) {
	err = d.StopContainer()
	if err == nil {
		err = d.RemoveContainer()
	}
	if err == nil {
		err = d.RunImage()
	}
	return
}

func (d *DockerDeploy) FindContainer() (containerID string, err error) {
	containers, err := d.Client.ContainerList(d.Context, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, item := range containers {
		if item.Names[0] == d.App.Name {
			containerID = item.ID
			break
		}
	}
	return
}

func (d *DockerDeploy) StopContainer() (err error) {
	timeout := time.Second

	var containerID string
	if containerID, err = d.FindContainer(); err == nil {
		err = d.Client.ContainerStop(d.Context, containerID, &timeout)
	}
	return
}

func (d *DockerDeploy) RemoveContainer() (err error) {
	var containerID string
	if containerID, err = d.FindContainer(); err == nil {
		err = d.Client.ContainerRemove(d.Context, containerID, types.ContainerRemoveOptions{})
	}
	return
}

func (d *DockerDeploy) RunImage() (err error) {
	var containerBody container.ContainerCreateCreatedBody
	containerBody, err = d.Client.ContainerCreate(d.Context, &container.Config{
		Image: d.App.Image,
	}, nil, nil, d.App.Name)
	if err == nil {
		err = d.Client.ContainerStart(d.Context, containerBody.ID, types.ContainerStartOptions{})
	}
	return
}
