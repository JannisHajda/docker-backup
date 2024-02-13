package docker

import (
	"context"
	"docker-backup/interfaces"
	"fmt"
	"github.com/docker/docker/api/types/volume"

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
	var binds []interfaces.DockerBind
	for _, m := range c.Mounts {
		if m.Type == "volume" {
			volumes = append(volumes, NewDockerVolume(m.Name, m.Source, m.RW))
		} else if m.Type == "bind" {
			binds = append(binds, NewDockerBind(m.Source, m.Destination, m.RW))
		}
	}

	return NewDockerContainer(d, c.ID, c.Name, volumes, binds), nil
}

func (d *DockerClient) CreateContainer(image string, volumes []interfaces.DockerVolume, binds []interfaces.DockerBind) (interfaces.DockerContainer, error) {
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

	for _, b := range binds {
		hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
			Type:   "bind",
			Source: b.GetSource(),
			Target: b.GetTarget(),
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

	return NewDockerContainer(d, container.GetID(), container.GetName(), volumes, binds), nil
}

func (d *DockerClient) GetVolume(name string) (interfaces.DockerVolume, error) {
	v, err := d.client.VolumeInspect(context.Background(), name)
	if err != nil {
		return nil, err
	}

	return NewDockerVolume(v.Name, v.Mountpoint, true), nil
}

func (d *DockerClient) CreateVolume(name string) (interfaces.DockerVolume, error) {
	v, err := d.client.VolumeCreate(context.Background(), volume.CreateOptions{
		Name: name,
	})

	if err != nil {
		return nil, err
	}

	fmt.Printf("Created volume %s\n", v.Name)

	return NewDockerVolume(name, v.Mountpoint, true), nil
}
