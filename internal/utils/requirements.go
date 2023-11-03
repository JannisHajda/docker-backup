package utils

import (
	"os"
	"os/exec"
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
	// run borg --version
	// if it fails, return BorgNotInstalledError
	// if it succeeds, return nil
	cmd := exec.Command("borg", "--version")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

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
