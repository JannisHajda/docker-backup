package borgclient

import (
	"docker-backup/errors"
	"docker-backup/interfaces"
	"fmt"
	"regexp"
)

type BorgClient struct {
	container interfaces.DockerContainer
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

// Remote-backup permission denied error
func (b *BorgClient) isPermissionDeniedError(output string) bool {
	re := regexp.MustCompile(`Permission denied`)
	matches := re.FindStringSubmatch(output)

	if len(matches) > 0 {
		return true
	}

	return false
}

func (b *BorgClient) wrongPassphraseError(output string) bool {
	re := regexp.MustCompile("passphrase supplied in BORG_PASSPHRASE, by BORG_PASSCOMAND or via BORG_PASSPHRASE_FD is incorrect")
	matches := re.FindStringSubmatch(output)

	if len(matches) > 0 {
		return true
	}

	return false
}

func (b *BorgClient) ensureBorgIsInstalled() error {
	_, err := b.container.Exec("borg --version")
	if err != nil {
		return err
	}

	return nil
}

func (b *BorgClient) isRemoteHostNotFoundError(output string) bool {
	re := regexp.MustCompile("Could not resolve hostname")
	matches := re.FindStringSubmatch(output)

	if len(matches) > 0 {
		return true
	}

	return false
}

func (b *BorgClient) isHostKeyVerificationFailedError(output string) bool {
	re := regexp.MustCompile("Host key verification failed")
	matches := re.FindStringSubmatch(output)

	if len(matches) > 0 {
		return true
	}

	return false
}

func NewBorgClient(c interfaces.DockerContainer) (interfaces.BorgClient, error) {
	bc := &BorgClient{container: c}
	err := bc.ensureBorgIsInstalled()

	if err != nil {
		return nil, err
	}

	return bc, nil
}

func (b *BorgClient) GetRepository(path string, passphrase string) (interfaces.BorgRepository, error) {
	b.container.SetEnv("BORG_PASSPHRASE", passphrase)

	_, err := b.container.Exec("borg list " + path)
	if err != nil {
		err = b.handleError(err)
		return nil, err
	}

	return NewBorgRepository(b, path, passphrase)
}

func (b *BorgClient) GetContainer() interfaces.DockerContainer {
	return b.container
}

func (b *BorgClient) CreateRepository(path string, passphrase string) (interfaces.BorgRepository, error) {
	b.container.SetEnv("BORG_PASSPHRASE", passphrase)

	_, err := b.container.Exec("borg init " + path + " -e repokey-blake2 --make-parent-dirs")
	if err != nil {
		err = b.handleError(err)
		return nil, err
	}

	return NewBorgRepository(b, path, passphrase)
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
