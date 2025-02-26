package docker

import (
	"bytes"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
	"strings"
)

type Container struct {
	container.InspectResponse
	client *Client
}

func NewContainer(inspect container.InspectResponse, cli *Client) *Container {
	return &Container{inspect, cli}
}

func (c *Container) Exec(cmd string) (string, string, int, error) {
	ctx := c.client.GetContext()

	var execCmd []string
	if strings.ContainsAny(cmd, "><|;&") {
		execCmd = []string{"sh", "-c", cmd}
	} else {
		execCmd = strings.Fields(cmd)
	}

	execConfig := container.ExecOptions{
		Cmd:          execCmd,
		AttachStdout: true,
		AttachStderr: true,
	}

	resp, err := c.client.ContainerExecCreate(ctx, c.ID, execConfig)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to create exec: %w", err)
	}

	attachResp, err := c.client.ContainerExecAttach(ctx, resp.ID, container.ExecAttachOptions{})
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to attach exec: %w", err)
	}
	defer attachResp.Close()

	var stdoutBuf, stderrBuf bytes.Buffer
	_, err = stdcopy.StdCopy(&stdoutBuf, &stderrBuf, attachResp.Reader)
	if err != nil {
		return stdoutBuf.String(), stderrBuf.String(), 0, fmt.Errorf("failed to demultiplex output: %w", err)
	}

	inspect, err := c.client.ContainerExecInspect(ctx, resp.ID)
	if err != nil {
		return stdoutBuf.String(), stderrBuf.String(), 0, fmt.Errorf("failed to inspect exec: %w", err)
	}

	return stdoutBuf.String(), stderrBuf.String(), inspect.ExitCode, nil
}

func (c *Container) Start() error {
	return c.client.ContainerStart(c.client.ctx, c.ID, container.StartOptions{})
}

func (c *Container) Stop() error {
	return c.client.ContainerStop(c.client.ctx, c.ID, container.StopOptions{})
}

func (c *Container) Pause() error {
	return c.client.ContainerPause(c.client.ctx, c.ID)
}

func (c *Container) Unpause() error {
	return c.client.ContainerUnpause(c.client.ctx, c.ID)
}

func (c *Container) Remove() error {
	return c.client.ContainerRemove(c.client.ctx, c.ID, container.RemoveOptions{})
}

func (c *Container) StopAndRemove() error {
	if err := c.Stop(); err != nil {
		return fmt.Errorf("failed to unpause container %s: %v", c.ID, err)
	}

	return c.Remove()
}
