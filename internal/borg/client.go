package borg

import (
	"docker-backup/errors"
	"docker-backup/interfaces"
)

type BorgClient struct {
	container interfaces.DockerContainer
}

func (b *BorgClient) ensureBorgIsInstalled() error {
	_, err := b.container.Exec("borg --version")
	if err != nil {
		return errors.NewBorgNotInstalledError(err)
	}

	return nil
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
		err = errors.HandleBorgClientError(err)
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
		err = errors.HandleBorgClientError(err)
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
