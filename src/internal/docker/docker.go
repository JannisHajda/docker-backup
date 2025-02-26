package docker

import (
	"context"
	dockerClient "github.com/docker/docker/client"
)

type Client struct {
	dockerClient.Client
	ctx         context.Context
	workerImage string
}

func NewClient(ctx context.Context, workerImage string) (*Client, error) {
	cli, err := dockerClient.NewClientWithOpts(dockerClient.FromEnv, dockerClient.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &Client{Client: *cli, ctx: ctx, workerImage: workerImage}, nil
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
