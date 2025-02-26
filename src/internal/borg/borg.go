package borg

import (
	"fmt"
	"github.com/JannisHajda/docker-backup/internal/interfaces"
	"strings"
)

type Client struct {
	shell interfaces.Shell
}

func NewClient(shell interfaces.Shell) (*Client, error) {
	return &Client{shell: shell}, nil
}

func (c *Client) appendAuth(cmd string, pass string) string {
	cmd = fmt.Sprintf("export BORG_PASSPHRASE=%s && %s", pass, cmd)
	return cmd
}

func (c *Client) CreateRepo(path string, pass string) (*Repo, error) {
	cmd := fmt.Sprintf("borg init --encryption=repokey %s", path)
	cmd = c.appendAuth(cmd, pass)
	_, _, exitCode, err := c.shell.Exec(cmd)
	if err != nil {
		return nil, err
	}

	if exitCode != 0 {
		return nil, fmt.Errorf("failed to initialize borg repository: exit code %d", exitCode)
	}

	name := strings.Split(path, "/")[len(strings.Split(path, "/"))-1]
	return &Repo{client: c, Path: path, Archives: nil, Name: name}, nil
}

func (c *Client) GetRepo(path string, pass string) (*Repo, error) {
	cmd := fmt.Sprintf("borg list %s", path)
	cmd = c.appendAuth(cmd, pass)
	_, _, exitCode, err := c.shell.Exec(cmd)
	if err != nil {
		return nil, err
	}

	if exitCode != 0 {
		return nil, fmt.Errorf("failed to list borg repository: exit code %d", exitCode)
	}

	name := strings.Split(path, "/")[len(strings.Split(path, "/"))-1]
	return &Repo{client: c, Path: path, Archives: nil, Name: name}, nil
}
