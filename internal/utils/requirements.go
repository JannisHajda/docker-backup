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

type BorgNotInstalledError struct {
	Err error
}

func (bnie BorgNotInstalledError) Error() string {
	return "Borg is not installed"
}

type BorgNoPermissionError struct {
	Err error
}

func (bnpe BorgNoPermissionError) Error() string {
	return "No permission to execute borg"
}

func EnsureBorgInstalled() error {
	_, err := os.Stat("/usr/bin/borg")

	if os.IsNotExist(err) {
		return BorgNotInstalledError{Err: err}
	}

	if os.IsPermission(err) {
		return BorgNoPermissionError{Err: err}
	}

	if err != nil {
		return err
	}

	return nil
}
