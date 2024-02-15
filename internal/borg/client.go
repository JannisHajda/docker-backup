package borg

import (
	"docker-backup/errors"
	"docker-backup/interfaces"
	"strings"
)

type BorgClient struct {
	container interfaces.DockerContainer
}

func (b *BorgClient) setPassphrase(passphrase string) {
	b.container.SetEnv("BORG_PASSPHRASE", passphrase)
}

func (b *BorgClient) ensureBorgIsInstalled() error {
	cmd := "which borg"
	_, err := b.container.Exec(cmd)
	if err != nil {
		return errors.NewBorgNotInstalledError(err)
	}

	return nil
}

func NewBorgClient(c interfaces.DockerContainer) (interfaces.BorgClient, error) {
	b := &BorgClient{container: c}

	err := b.ensureBorgIsInstalled()
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (b *BorgClient) GetRepo(path string, passphrase string) (interfaces.BorgRepo, error) {
	b.container.SetEnv("BORG_PASSPHRASE", passphrase)

	_, err := b.container.Exec("borg list " + path)
	if err != nil {
		err = errors.HandleBorgClientError(err)
		return nil, err
	}

	return &BorgRepo{BorgClient: b, path: path, passphrase: passphrase}, nil
}

func (b *BorgClient) GetContainer() interfaces.DockerContainer {
	return b.container
}

func validateEncryptionType(encryptionType string) error {
	supportedTypes := []string{"repokey", "repokey-blake2", "keyfile", "keyfile-blake2", "authenticated", "authenticated-blake2", "none"}
	encryptionType = strings.ToLower(encryptionType)

	for _, t := range supportedTypes {
		if t == encryptionType {
			return nil
		}
	}

	return errors.NewBorgUnknownEncryptionTypeError(encryptionType)
}

func (b *BorgClient) CreateRepo(config interfaces.CreateBorgRepoConfig) (interfaces.BorgRepo, error) {
	b.setPassphrase(config.Passphrase)

	cmd := "borg init"
	err := validateEncryptionType(config.EncryptionType)
	if err != nil {
		return nil, err
	}

	if config.MakeParentDirs {
		cmd += " --make-parent-dirs"
	}

	if config.AppendOnly {
		cmd += " --append-only"
	}

	if config.StorageQuota != "" {
		cmd += " --storage-quota " + config.StorageQuota
	}

	_, err = b.container.Exec(cmd)
	if err != nil {
		err = errors.HandleBorgClientError(err)
		return nil, err
	}

	return &BorgRepo{BorgClient: b, path: config.Path, passphrase: config.Passphrase}, nil
}

func (b *BorgClient) GetOrCreateRepo(config interfaces.CreateBorgRepoConfig) (interfaces.BorgRepo, error) {
	repo, err := b.GetRepo(config.Path, config.Passphrase)

	if err != nil {
		if _, ok := err.(*errors.RepositoryDoesNotExistError); ok {
			repo, err = b.CreateRepo(config)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return repo, nil
}
