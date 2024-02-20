package ssh

import (
	"docker-backup/errors"
	"docker-backup/interfaces"
	"docker-backup/internal/helper"
	"fmt"
	"strings"
)

type SSHClient struct {
	container interfaces.DockerContainer
}

const (
	SSH_FOLDER = "~/.ssh"

	HOST_KEY_VERIFICATION_FAILED = "Host key verification failed"
	REMOTE_HOST_NOT_FOOUND       = "Could not resolve hostname"
)

func NewSSHClient(container interfaces.DockerContainer, keyfiles []interfaces.DockerBind, hosts []string) (interfaces.SSHClient, []error) {
	client := &SSHClient{container: container}
	var errs []error
	err := client.initSSHFolder()
	if err != nil {
		errs = append(errs, err)
		return nil, errs
	}

	for _, keyfile := range keyfiles {
		keyfilePath := keyfile.GetMountPoint()
		err = client.AddKeyfile(keyfilePath)
		if err != nil {
			errs = append(errs, err)
		}
	}

	for _, host := range hosts {
		err = client.AddKnownHost(host)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return client, errs
	}

	return client, nil
}

func (s *SSHClient) handleError(e error) error {
	if helper.RegexMatch(e.Error(), HOST_KEY_VERIFICATION_FAILED) {
		return errors.NewHostKeyVerificationFailedError(e)
	}

	if helper.RegexMatch(e.Error(), REMOTE_HOST_NOT_FOOUND) {
		return errors.NewRemoteHostNotFoundError(e)
	}

	return e
}

func (s *SSHClient) initSSHFolder() error {
	cmd := "mkdir -p ~/.ssh"
	_, err := s.container.Exec(cmd)
	if err != nil {
		return s.handleError(err)
	}

	return nil
}

func (s *SSHClient) AddKeyfile(keyfile string) error {
	keyName := strings.Split(keyfile, "/")[len(strings.Split(keyfile, "/"))-1]

	// copy key to ssh folder
	cmd := fmt.Sprintf("cp %s %s/%s", keyfile, SSH_FOLDER, keyName)
	_, err := s.container.Exec(cmd)
	if err != nil {
		return s.handleError(err)
	}

	// add key to ssh-agent
	cmd = fmt.Sprintf("eval `ssh-agent` && ssh-add %s/%s", SSH_FOLDER, keyName)
	_, err = s.container.Exec(cmd)
	if err != nil {
		return s.handleError(err)
	}

	return nil
}

func (s *SSHClient) AddKnownHost(host string) error {
	cmd := fmt.Sprintf("ssh-keyscan -H %s >> %s/known_hosts", host, SSH_FOLDER)
	_, err := s.container.Exec(cmd)
	if err != nil {
		return s.handleError(err)
	}

	return nil
}
