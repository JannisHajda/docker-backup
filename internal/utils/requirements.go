package utils

import (
	"os"
)

type NoAccessToDockerSocketError struct {
	Err error
}

func (nats NoAccessToDockerSocketError) Error() string {
	return "No access to docker socket"
}

func EnsureAccessToDockerSocket() error {
	socketPath := "/var/run/docker.sock"
	_, err := os.Stat(socketPath)

	if os.IsNotExist(err) {
		return NoAccessToDockerSocketError{Err: err}
	}

	if os.IsPermission(err) {
		return NoAccessToDockerSocketError{Err: err}
	}

	if err != nil {
		return err
	}

	return nil
}