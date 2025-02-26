package rclone

import (
	"fmt"
	"github.com/JannisHajda/docker-backup/internal/interfaces"
)

type Client struct {
	shell   interfaces.Shell
	Remotes []*Remote
}

func NewClient(shell interfaces.Shell) (*Client, error) {
	return &Client{shell: shell}, nil
}

type Remote struct {
	Provider string
	Name     string
	User     string
	Pass     string
}

type SyncConfig struct {
	InputPath  string
	OutputPath string
	Remote     Remote
}

func (c *Client) Sync(config SyncConfig) error {
	cmd := fmt.Sprintf("rclone sync %s %s:%s", config.InputPath, config.Remote.Name, config.OutputPath)
	_, stderr, exitCode, err := c.shell.Exec(cmd)
	if err != nil {
		return err
	}

	if exitCode != 0 {
		return fmt.Errorf("failed to sync rclone: exit code %d, stderr: %s", exitCode, stderr)
	}

	return nil
}
