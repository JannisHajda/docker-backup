package borgClient

import (
	"os"
	"os/exec"

	"github.com/JannisHajda/docker-backup/internal/utils"
)

type BorgClient struct {
}

func NewBorgClient() (*BorgClient, error) {
	err := utils.EnsureBorgInstalled()

	if err != nil {
		return nil, err
	}

	return &BorgClient{}, nil
}

type BorgRepoAlreadyExistsError struct {
	Err error
}

func (brae BorgRepoAlreadyExistsError) Error() string {
	return "Borg repo already exists"
}

func (bc *BorgClient) InitializeRepo(name string, path string) error {
	// Check if the repo directory already exists
	_, err := os.Stat(path)
	if err == nil {
		return BorgRepoAlreadyExistsError{Err: err}
	} else if !os.IsNotExist(err) {
		return err
	}

	// Create the repo folder
	err = os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}

	// Initialize the repo
	cmd := exec.Command("borg", "init", "--encryption=repokey", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
