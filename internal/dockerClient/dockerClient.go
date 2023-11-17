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

type DockerContainer struct {
	ID      string
	Name    string
	Volumes []DockerVolume
}

type DockerVolume struct {
	Name     string
	ReadOnly bool
}

func (dc *DockerClient) GetContainerByID(id string) (*DockerContainer, error) {
	opts := docker.InspectContainerOptions{
		ID: id,
	}

	c, err := dc.client.InspectContainerWithOptions(opts)

	if err != nil {
		return nil, err
	}

	return &DockerContainer{
		ID:      c.ID,
		Name:    c.Name,
		Volumes: getContainerVolumes(c),
	}, nil
}

func getContainerVolumes(c *docker.Container) []DockerVolume {
	mounts := c.HostConfig.Mounts
	volumes := []DockerVolume{}

	for _, m := range mounts {
		if m.Type == "volume" {
			volumes = append(volumes, DockerVolume{
				Name:     m.Source,
				ReadOnly: m.ReadOnly,
			})
		}
	}

	return volumes
}

func (dc *DockerClient) GetAllContainers() ([]docker.APIContainers, error) {
	opts := docker.ListContainersOptions{}

	containers, err := dc.client.ListContainers(opts)

	if err != nil {
		return nil, err
	}

	return containers, nil
}
