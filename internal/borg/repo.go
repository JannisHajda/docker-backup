package borg

import (
	"docker-backup/errors"
	"time"
)

type BorgRepo struct {
	*BorgClient
	path       string
	passphrase string
	key        string
}

func (b *BorgRepo) GetPath() string {
	return b.path
}

func (b *BorgRepo) GetArchives() (string, error) {
	b.container.SetEnv("BORG_PASSPHRASE", b.passphrase)

	output, err := b.container.Exec("borg list " + b.path)
	if err != nil {
		err = errors.HandleBorgClientError(err)
		return "", err
	}

	return output, nil
}

func (r *BorgRepo) Backup(input string) error {
	r.container.SetEnv("BORG_PASSPHRASE", r.passphrase)
	now := time.Now().Format("2006-01-02T15:04:05")

	_, err := r.container.Exec("borg create " + r.path + "::" + now + " " + input)
	if err != nil {
		return err
	}

	return nil
}
