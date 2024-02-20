package borg

import (
	"docker-backup/errors"
	"docker-backup/interfaces"
	"fmt"
	"strings"
)

type BorgRepo struct {
	*BorgClient
	path       string
	passphrase string
	keyfile    string
}

func NewBorgRepo(client *BorgClient, path string, passphrase string, keyfile string) *BorgRepo {
	r := &BorgRepo{
		BorgClient: client,
		path:       path,
		passphrase: passphrase,
		keyfile:    keyfile,
	}

	return r
}

func (b *BorgRepo) authenticate() {
	b.setPassphrase(b.passphrase)
	b.setKeyfile(b.keyfile)
}

func (b *BorgRepo) validateCompression(compression string) error {
	supportedCompressions := []string{"none", "lz4", "zstd", "zlib", "lzma"}
	compression = strings.ToLower(compression)

	for _, c := range supportedCompressions {
		if c == compression {
			return nil
		}
	}

	return errors.NewBorgUnknownCompressionTypeError(compression)
}

func (b *BorgRepo) ListArchives() (string, error) {
	b.authenticate()

	output, err := b.container.Exec("borg list " + b.path)
	if err != nil {
		return "", b.handleError(err)
	}

	return output, nil
}

func (b *BorgRepo) CreateArchive(config interfaces.CreateBorgArchiveConfig) error {
	err := b.validateCompression(config.Compression)
	if err != nil {
		return err
	}

	b.authenticate()

	sources := strings.Join(config.Sources, " ")
	cmd := fmt.Sprintf("borg create --compression %s %s::%s %s", config.Compression, b.path, config.Name, sources)
	_, err = b.container.Exec(cmd)
	if err != nil {
		return b.handleError(err)
	}

	return nil
}

func (b *BorgRepo) Info() (string, error) {
	b.authenticate()

	output, err := b.container.Exec("borg info " + b.path)
	if err != nil {
		return "", b.handleError(err)
	}

	return output, nil
}
