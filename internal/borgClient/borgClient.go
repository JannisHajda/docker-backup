package borgclient

import (
	"docker-backup/errors"
	"docker-backup/interfaces"
	"fmt"
	"regexp"
)

type BorgClient struct {
	worker interfaces.Worker
}

func (b *BorgClient) handleError(e error) error {
	output := e.Error()
	if b.isRepositoryDoesNotExistError(output) {
		return errors.NewRepositoryDoesNotExistError(e)
	}

	if b.isRepositoryParentDirectoryDoesNotExistError(output) {
		return errors.NewRepositoryParentDirectoryDoesNotExistError(e)
	}

	if b.isRepositoryAlreadyExistsError(output) {
		return errors.NewRepositoryAlreadyExistsError(e)
	}

	return fmt.Errorf("unknown error")
}

func (b *BorgClient) isRepositoryDoesNotExistError(output string) bool {
	re := regexp.MustCompile(`Repository (.*) does not exist.`)
	matches := re.FindStringSubmatch(output)

	if len(matches) > 0 {
		return true
	}

	return false
}

func (b *BorgClient) isRepositoryParentDirectoryDoesNotExistError(output string) bool {
	re := regexp.MustCompile(`parent path of the repo directory (.*) does not exist`)
	matches := re.FindStringSubmatch(output)

	if len(matches) > 0 {
		return true
	}

	return false
}

func (b *BorgClient) isRepositoryAlreadyExistsError(output string) bool {
	re := regexp.MustCompile(`repository already exists at`)
	matches := re.FindStringSubmatch(output)

	if len(matches) > 0 {
		return true
	}

	return false
}

func (b *BorgClient) ensureBorgIsInstalled() error {
	_, err := b.worker.Exec("borg --version")
	if err != nil {
		return err
	}

	return nil
}

func NewBorgClient(w interfaces.Worker) (interfaces.BorgClient, error) {
	bc := &BorgClient{worker: w}
	err := bc.ensureBorgIsInstalled()

	if err != nil {
		return nil, err
	}

	return bc, nil
}

func (b *BorgClient) GetRepository(path string, passphrase string) (interfaces.BorgRepository, error) {
	b.worker.SetEnv("BORG_PASSPHRASE", passphrase)

	_, err := b.worker.Exec("borg list " + path)
	if err != nil {
		err = b.handleError(err)
		return nil, err
	}

	return NewBorgRepository(b, path, passphrase)
}

func (b *BorgClient) CreateRepository(path string, passphrase string) (interfaces.BorgRepository, error) {
	b.worker.SetEnv("BORG_PASSPHRASE", passphrase)

	_, err := b.worker.Exec("borg init " + path + " -e repokey-blake2 --make-parent-dirs")
	if err != nil {
		err = b.handleError(err)
		return nil, err
	}

	return NewBorgRepository(b, path, passphrase)
}

// create/get all repositories
// extract key from config file for new repositories
func (b *BorgClient) PreBackup() error {

	return nil
}

func (b *BorgClient) GetOrCreateRepository(path string, passphrase string) (interfaces.BorgRepository, error) {
	repo, err := b.GetRepository(path, passphrase)
	if err != nil {
		if _, ok := err.(*errors.RepositoryDoesNotExistError); ok {
			repo, err = b.CreateRepository(path, passphrase)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return repo, nil
}
