package docker

import (
	"bytes"
	"context"
	"docker-backup/interfaces"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
)

type DockerContainer struct {
	*DockerClient
	id      string
	name    string
	envVars map[string]string
	volumes []interfaces.DockerVolume
	binds   []interfaces.DockerBind
}

func NewDockerContainer(c *DockerClient, id string, name string, volumes []interfaces.DockerVolume, binds []interfaces.DockerBind) interfaces.DockerContainer {
	return &DockerContainer{DockerClient: c, id: id, name: name, envVars: map[string]string{},
		volumes: volumes, binds: binds}
}

func (d *DockerContainer) GetID() string {
	return d.id
}

func (d *DockerContainer) GetName() string {
	return d.name
}

func (d *DockerContainer) GetVolumes() []interfaces.DockerVolume {
	return d.volumes
}

func (d *DockerContainer) GetBinds() []interfaces.DockerBind {
	return d.binds
}

func (d *DockerContainer) SetEnv(key string, value string) {
	d.envVars[key] = value
}

func (d *DockerContainer) RemoveEnv(key string) {
	delete(d.envVars, key)
}

func (d *DockerContainer) GetEnv(key string) string {
	return d.envVars[key]
}

func (d *DockerContainer) GetEnvs() map[string]string {
	return d.envVars
}

func (d *DockerContainer) Start() error {
	return d.client.ContainerStart(context.Background(), d.id, types.ContainerStartOptions{})
}

func (d *DockerContainer) Stop() error {
	return d.client.ContainerStop(context.Background(), d.id, container.StopOptions{})
}

func (d *DockerContainer) Remove() error {
	return d.client.ContainerRemove(context.Background(), d.id, types.ContainerRemoveOptions{})
}

func (d *DockerContainer) StopAndRemove() error {
	err := d.Stop()
	if err != nil {
		return err
	}

	return d.Remove()
}

func (d *DockerContainer) envMapToSlice() []string {
	var envs []string
	for key, value := range d.envVars {
		envs = append(envs, fmt.Sprintf("%s=%s", key, value))
	}

	return envs
}

func (d *DockerContainer) Exec(cmd string) (string, error) {
	execCreateResp, err := d.client.ContainerExecCreate(context.Background(), d.id, types.ExecConfig{
		Cmd:          []string{"/bin/bash", "-c", cmd},
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
		Env:          d.envMapToSlice(), // You may customize the environment if needed
	})

	if err != nil {
		return "", err
	}

	res, err := d.client.ContainerExecAttach(context.Background(), execCreateResp.ID, types.ExecStartCheck{})
	if err != nil {
		return "", err
	}

	defer res.Close()

	var outBuf, errBuf bytes.Buffer
	outputDone := make(chan error)

	go func() {
		_, err := stdcopy.StdCopy(&outBuf, &errBuf, res.Reader)
		outputDone <- err
	}()

	select {
	case err := <-outputDone:
		if err != nil {
			return "", err
		}

	case <-context.Background().Done():
		return "", context.Background().Err()
	}

	stdout, err := io.ReadAll(&outBuf)
	if err != nil {
		return "", err
	}

	stderr, err := io.ReadAll(&errBuf)
	if err != nil {
		return "", err
	}

	inspect, err := d.client.ContainerExecInspect(context.Background(), execCreateResp.ID)
	if err != nil {
		return "", err
	}

	if inspect.ExitCode != 0 {
		return "", fmt.Errorf("Command failed with exit code %d: %s", inspect.ExitCode, stderr)
	}

	return string(stdout), nil
}
