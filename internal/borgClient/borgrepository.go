package borgclient

import (
	"docker-backup/interfaces"
	"time"
)

type BorgRepository struct {
	*BorgClient
	name       string
	passphrase string
	key        string
}

func NewBorgRepository(c *BorgClient, name string, passphrase string) (interfaces.BorgRepository, error) {
	return &BorgRepository{
		BorgClient: c, name: name, passphrase: passphrase}, nil
}

func (b *BorgRepository) GetName() string {
	return b.name
}

func (b *BorgRepository) Backup() error {
	b.container.SetEnv("BORG_PASSPHRASE", b.passphrase)
	now := time.Now().Format("2006-01-02T15:04:05")

	input := b.inputDir + "/" + b.name
	output := b.outputDir + "/" + b.name

	_, err := b.container.Exec("borg create " + output + "::" + now + " " + input)
	if err != nil {
		return err
	}

	return nil
}
