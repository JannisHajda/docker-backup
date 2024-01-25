package borgclient

import (
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
	return &BorgRepository{
		BorgClient: c, path: path, passphrase: passphrase}, nil
}

func (b *BorgRepository) GetPath() string {
	return b.path
}

func (b *BorgRepository) Archive(inputPath string) error {
	b.worker.SetEnv("BORG_PASSPHRASE", b.passphrase)
	now := time.Now().Format("2006-01-02T15:04:05")
	_, err := b.worker.Exec("borg create " + b.path + "::" + now + " " + inputPath)
	if err != nil {
		return err
	}

	return nil
}
