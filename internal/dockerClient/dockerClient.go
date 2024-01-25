package dockerclient

import (
	"context"
	"docker-backup/interfaces"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"

	"os/exec"
)

type DockerClient struct {
	client client.Client
}

func ensureDockerIsInstalled() error {
	cmd := exec.Command("docker", "version")
	err := cmd.Run()

	return err
}

func NewDockerClient() (interfaces.DockerClient, error) {
	err := ensureDockerIsInstalled()
	if err != nil {
		return nil, err
	}

	client, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	return &DockerClient{client: *client}, nil
}

func (d *DockerClient) GetContainer(id string) (interfaces.DockerContainer, error) {
	c, err := d.client.ContainerInspect(context.Background(), id)
	if err != nil {
		return nil, err
	}

	var volumes []interfaces.DockerVolume
	for _, v := range c.Mounts {
		if v.Type == "volume" {
			volumes = append(volumes, NewDockerVolume(v.Name, v.RW))
		}
	}

	return NewDockerContainer(d, c.ID, c.Name, volumes), nil
}

func (d *DockerClient) CreateContainer(image string, volumes []interfaces.DockerVolume) (interfaces.DockerContainer, error) {
	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{},
	}

	for _, v := range volumes {
		hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
			Type:   "volume",
			Source: v.GetName(),
			Target: v.GetMountPoint(),
		})
	}

	c, err := d.client.ContainerCreate(context.Background(), &container.Config{
		Image: image,
	}, hostConfig, nil, nil, "")

	if err != nil {
		return nil, err
	}

	container, err := d.GetContainer(c.ID)
	if err != nil {
		return nil, err
	}

	return NewDockerContainer(d, container.GetID(), container.GetName(), volumes), nil
}
