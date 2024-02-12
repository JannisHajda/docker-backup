package sshclient

import (
	"docker-backup/errors"
	"docker-backup/interfaces"
	"fmt"
	"strings"
)

type SSHClient struct {
	container interfaces.DockerContainer
}

func NewSSHClient(container interfaces.DockerContainer, keyfiles []interfaces.DockerBind, hosts []string) (interfaces.SSHClient, []error) {
	client := &SSHClient{container: container}
	var errs []error
	err := client.initSSHFolder()
	if err != nil {
		errs = append(errs, err)
		return nil, errs
	}

	for _, keyfile := range keyfiles {
		err = client.AddKey(keyfile.GetTarget())
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
	if errors.IsErrOfKind(e, errors.HOST_KEY_VERIFICATION_FAILED) {
		return errors.NewHostKeyVerificationFailedError(e)
	}

	if errors.IsErrOfKind(e, errors.REMOTE_HOST_NOT_FOOUND) {
		return errors.NewRemoteHostNotFoundError(e)
	}

	return fmt.Errorf("unknown error")
}

func (s *SSHClient) initSSHFolder() error {
	cmd := "mkdir -p ~/.ssh"
	_, err := s.container.Exec(cmd)
	if err != nil {
		return err
	}

	return nil
}

func (s *SSHClient) AddKey(path string) error {
	keyName := strings.Split(path, "/")[len(strings.Split(path, "/"))-1]

	// copy key to ssh folder
	cmd := fmt.Sprintf("cp %s /root/.ssh/%s", path, keyName)
	_, err := s.container.Exec(cmd)
	if err != nil {
		return err
	}

	// add key to ssh-agent
	cmd = fmt.Sprintf("eval `ssh-agent` && ssh-add /root/.ssh/%s", keyName)
	_, err = s.container.Exec(cmd)
	if err != nil {
		return err
	}

	return nil
}

func (s *SSHClient) AddKnownHost(host string) error {
	cmd := fmt.Sprintf("ssh-keyscan -H %s >> /root/.ssh/known_hosts", host)
	_, err := s.container.Exec(cmd)
	if err != nil {
		return err
	}

	return nil
}
