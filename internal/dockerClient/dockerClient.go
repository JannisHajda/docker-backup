package dockerClient

import (
	"github.com/JannisHajda/docker-backup/internal/utils"
	docker "github.com/fsouza/go-dockerclient"
)

type DockerClient struct {
	client *docker.Client
}

func NewDockerClient() (*DockerClient, error) {
	err := utils.EnsureAccessToDockerSocket()

	if err != nil {
		return nil, err
	}

	client, err := docker.NewClientFromEnv()

	if err != nil {
		return nil, err
	}

	return &DockerClient{
		client: client,
	}, nil
}

func (dc *DockerClient) GetContainer(id string) (*docker.Container, error) {
	opts := docker.InspectContainerOptions{
		ID: id,
	}

	c, err := dc.client.InspectContainerWithOptions(opts)

	if err != nil {
		return nil, err
	}

	return c, nil
}

func (dc *DockerClient) GetAllContainers() ([]docker.APIContainers, error) {
	opts := docker.ListContainersOptions{}

	containers, err := dc.client.ListContainers(opts)

	if err != nil {
		return nil, err
	}

	return containers, nil
}
