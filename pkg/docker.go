package pkg

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"os"
	"time"
)

type DockerDeploy struct {
	Client  *client.Client
	App     Application
	Context context.Context
}

func (d *DockerDeploy) DeployImage() (err error) {
	var needUpdate bool
	if needUpdate, err = d.NeedUpdate(); err != nil {
		return
	}

	if !needUpdate {
		fmt.Println("do not need to update", d.App.Name)
		return
	}

	err = d.StopContainer()
	if err == nil {
		err = d.RemoveContainer()
	}
	if err == nil {
		err = d.RunImage()
	}
	return
}

func (d *DockerDeploy) FindContainerID() (containerID string, err error) {
	var result *types.Container
	if result, err = d.FindContainer(); err == nil && result != nil {
		containerID = result.ID
	}
	return
}

func (d *DockerDeploy) FindContainer() (result *types.Container, err error) {
	containers, err := d.Client.ContainerList(d.Context, types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return
	}

	for _, item := range containers {
		if item.Names[0] == "/"+d.App.Name {
			result = &item
			break
		}
	}
	return
}

func (d *DockerDeploy) StopContainer() (err error) {
	timeout := time.Minute

	var containerID string
	if containerID, err = d.FindContainerID(); err == nil && containerID != "" {
		err = d.Client.ContainerStop(d.Context, containerID, &timeout)
	}
	return
}

func (d *DockerDeploy) RemoveContainer() (err error) {
	var containerID string
	if containerID, err = d.FindContainerID(); err == nil && containerID != "" {
		fmt.Println("prepare to remove container ", containerID)
		err = d.Client.ContainerRemove(d.Context, containerID, types.ContainerRemoveOptions{})
	}
	return
}

func (d *DockerDeploy) NeedUpdate() (needUpdate bool, err error) {
	var result *types.Container
	if result, err = d.FindContainer(); err == nil && result != nil {
		needUpdate = d.App.Image != result.Image
	} else {
		needUpdate = true
	}
	return
}

func (d *DockerDeploy) RunImage() (err error) {
	var reader io.ReadCloser
	reader, err = d.Client.ImagePull(d.Context, d.App.Image, types.ImagePullOptions{})
	if err != nil {
		return
	}
	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		return
	}

	var containerBody container.ContainerCreateCreatedBody
	containerBody, err = d.Client.ContainerCreate(d.Context, d.getConfig(), d.getHostConfig(), nil, d.App.Name)
	if err == nil {
		err = d.Client.ContainerStart(d.Context, containerBody.ID, types.ContainerStartOptions{})
	}
	return
}

func (d *DockerDeploy) getConfig() *container.Config {
	config := &container.Config{
		Image:   d.App.Image,
		Volumes: map[string]struct{}{},
	}
	return config
}

func (d *DockerDeploy) getHostConfig() *container.HostConfig {
	portBindings := make(map[nat.Port][]nat.PortBinding, 0)
	for _, expose := range d.App.Exposes {
		port, _ := nat.NewPort("tcp", fmt.Sprintf("%d", expose.Container))

		portBindings[port] = []nat.PortBinding{{
			HostIP:   "0.0.0.0",
			HostPort: fmt.Sprintf("%d", expose.Host),
		}}
	}

	mounts := make([]mount.Mount, 0)
	for _, volume := range d.App.Volumes {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: volume.Host,
			Target: volume.Container,
		})
	}

	return &container.HostConfig{
		PortBindings: portBindings,
		Mounts:       mounts,
	}
}
