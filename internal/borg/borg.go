package borg

import (
	"docker-backup/errors"
	"docker-backup/interfaces"
	"docker-backup/internal/helper"
	"fmt"
	"strings"
)

const (
	REPO_DOES_NOT_EXIST            = "Repository .* does not exist"
	REPO_PARENT_DIR_DOES_NOT_EXIST = "parent path of the repo directory (.*) does not exist"
	REPO_ALREADY_EXISTS            = "Repository .* already exists"
	WRONG_PASSPHRASE               = "Wrong passphrase"
	PERMISSION_DENIED              = "Permission denied"
)

type BorgClient struct {
	container interfaces.DockerContainer
}

func NewBorgClient(c interfaces.DockerContainer) (interfaces.BorgClient, error) {
	b := &BorgClient{container: c}

	err := b.ensureBorgIsInstalled()
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (b *BorgClient) setPassphrase(passphrase string) {
	if passphrase == "" {
		b.container.RemoveEnv("BORG_PASSPHRASE")
		return
	}

	b.container.SetEnv("BORG_PASSPHRASE", passphrase)
}

func (b *BorgClient) setKeyfile(keyfile string) {
	if keyfile == "" {
		b.container.RemoveEnv("BORG_KEY_FILE")
		return
	}

	b.container.SetEnv("BORG_KEY_FILE", keyfile)
}

func (b *BorgClient) ensureBorgIsInstalled() error {
	cmd := "which borg"
	_, err := b.container.Exec(cmd)
	if err != nil {
		return errors.NewBorgNotInstalledError(err)
	}

	return nil
}

func (b *BorgClient) handleError(err error) error {
	if helper.RegexMatch(err.Error(), REPO_DOES_NOT_EXIST) {
		return errors.NewRepositoryDoesNotExistError(err)
	}

	if helper.RegexMatch(err.Error(), REPO_PARENT_DIR_DOES_NOT_EXIST) {
		return errors.NewRepositoryParentDirectoryDoesNotExistError(err)
	}

	if helper.RegexMatch(err.Error(), REPO_ALREADY_EXISTS) {
		return errors.NewRepositoryAlreadyExistsError(err)
	}

	if helper.RegexMatch(err.Error(), WRONG_PASSPHRASE) {
		return errors.NewWrongPassphraseError(err)
	}

	if helper.RegexMatch(err.Error(), PERMISSION_DENIED) {
		return errors.NewBorgPermissionDeniedError(err)
	}

	return err
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

func (b *BorgClient) GetRepo(config interfaces.GetBorgRepoConfig) (interfaces.BorgRepo, error) {
	r := NewBorgRepo(b, config.Path, config.Passphrase, config.Keyfile)
	_, err := r.Info()
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (b *BorgClient) CreateRepo(config interfaces.CreateBorgRepoConfig) (interfaces.BorgRepo, error) {
	err := validateEncryptionType(config.EncryptionType)
	if err != nil {
		return nil, err
	}

	cmd := fmt.Sprintf("borg init --encryption=%s ", config.EncryptionType)

	if config.MakeParentDirs {
		cmd += " --make-parent-dirs "
	}

	if config.AppendOnly {
		cmd += " --append-only "
	}

	if config.StorageQuota != "" {
		cmd += " --storage-quota " + config.StorageQuota
	}

	cmd += config.Path

	b.setPassphrase(config.Passphrase)
	b.setKeyfile(config.Keyfile)

	_, err = b.container.Exec(cmd)
	if err != nil {
		return nil, b.handleError(err)
	}

	return NewBorgRepo(b, config.Path, config.Passphrase, config.Keyfile), nil
}
