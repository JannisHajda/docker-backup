package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	dockerClient "github.com/docker/docker/client"
)

type Client struct {
	dockerClient.Client
	ctx context.Context
}

const (
	workerImage = "worker"
)

func NewClient(ctx context.Context) (*Client, error) {
	cli, err := dockerClient.NewClientWithOpts(dockerClient.FromEnv, dockerClient.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &Client{Client: *cli, ctx: ctx}, nil
}

func (cli *Client) GetContainer(id string) (*Container, error) {
	inspect, err := cli.ContainerInspect(cli.ctx, id)
	if err != nil {
		return nil, err
	}

	return NewContainer(inspect, cli), nil
}

func (cli *Client) GetContext() context.Context {
	return cli.ctx
}

func (cli *Client) SpawnWorker(repoName string, repoPassphrase string, mounts []mount.Mount) (*Worker, error) {
	config := container.Config{
		Image: workerImage,
		Cmd:   []string{"tail", "-f", "/dev/null"},
		Env: []string{
			fmt.Sprintf("BORG_REPO=/output/%s", repoName),
			"BORG_PASSPHRASE=" + repoPassphrase,
		},
	}

	hostConfig := container.HostConfig{
		Mounts: mounts,
	}

	resp, err := cli.ContainerCreate(cli.ctx, &config, &hostConfig, nil, nil, "")
	if err != nil {
		return nil, err
	}

	c, err := cli.GetContainer(resp.ID)
	if err != nil {
		return nil, err
	}

	w := Worker{Container: *c, repoName: repoName}

	if err := w.Start(cli.ctx); err != nil {
		return nil, err
	}

	return &w, nil
}
