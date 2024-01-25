package dockerclient

import (
	"context"
	"docker-backup/interfaces"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

type DockerContainer struct {
	*DockerClient
	id      string
	name    string
	envVars map[string]string
	volumes []interfaces.DockerVolume
}

func NewDockerContainer(c *DockerClient, id string, name string, volumes []interfaces.DockerVolume) interfaces.DockerContainer {
	return &DockerContainer{DockerClient: c, id: id, name: name, envVars: map[string]string{},
		volumes: volumes}
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

func (d *DockerContainer) SetEnv(key string, value string) {
	d.envVars[key] = value
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
		Cmd:          strings.Fields(cmd),
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false, // Set to true if you want a TTY-like behavior
		Env:          d.envMapToSlice(),
	})
	if err != nil {
		return "", err
	}

	execID := execCreateResp.ID

	execAttachResp, err := d.client.ContainerExecAttach(context.Background(), execID, types.ExecStartCheck{})
	if err != nil {
		return "", err
	}
	defer execAttachResp.Close()

	// Read the output from the attached response
	var output strings.Builder
	_, err = io.Copy(&output, execAttachResp.Reader)
	if err != nil {
		return "", err
	}

	// Start the execution
	err = d.client.ContainerExecStart(context.Background(), execID, types.ExecStartCheck{})
	if err != nil {
		return "", err
	}

	// Wait for the command to finish
	execInspectResp, err := d.client.ContainerExecInspect(context.Background(), execID)
	if err != nil {
		return "", err
	}

	if execInspectResp.ExitCode != 0 {
		return "", fmt.Errorf("Command failed with exit code %d: %s", execInspectResp.ExitCode, output.String())
	}

	return output.String(), nil
}
