package borgclient

import (
	"docker-backup/errors"
	"docker-backup/interfaces"
	"time"
)

type BorgRepository struct {
	*BorgClient
	path       string
	passphrase string
	key        string
}

func NewBorgRepository(c *BorgClient, path string, passphrase string) (interfaces.BorgRepository, error) {
	return &BorgRepository{BorgClient: c, path: path, passphrase: passphrase}, nil
}

func (b *BorgRepository) GetPath() string {
	return b.path
}

func (b *BorgRepository) GetArchives() (string, error) {
	b.container.SetEnv("BORG_PASSPHRASE", b.passphrase)

	output, err := b.container.Exec("borg list " + b.path)
	if err != nil {
		err = errors.HandleBorgClientError(err)
		return "", err
	}

	return output, nil
}

func (r *BorgRepository) Backup(input string) error {
	r.container.SetEnv("BORG_PASSPHRASE", r.passphrase)
	now := time.Now().Format("2006-01-02T15:04:05")

	_, err := r.container.Exec("borg create " + r.path + "::" + now + " " + input)
	if err != nil {
		return err
	}

	return nil
}
