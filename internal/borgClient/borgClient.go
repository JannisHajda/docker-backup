package borgclient

import (
	"docker-backup/errors"
	"docker-backup/interfaces"
	"fmt"
	"regexp"
)

type BorgClient struct {
	inputDir  string
	outputDir string
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

func (b *BorgClient) ensureBorgIsInstalled() error {
	_, err := b.container.Exec("borg --version")
	if err != nil {
		return err
	}

	return nil
}

func NewBorgClient(c interfaces.DockerContainer, inputDir string, outputDir string) (interfaces.BorgClient, error) {
	bc := &BorgClient{container: c, inputDir: inputDir, outputDir: outputDir}
	err := bc.ensureBorgIsInstalled()

	if err != nil {
		return nil, err
	}

	return bc, nil
}

func (b *BorgClient) GetRepository(name string, passphrase string) (interfaces.BorgRepository, error) {
	path := b.outputDir + "/" + name
	b.container.SetEnv("BORG_PASSPHRASE", passphrase)

	_, err := b.container.Exec("borg list " + path)
	if err != nil {
		err = b.handleError(err)
		return nil, err
	}

	return NewBorgRepository(b, path, passphrase)
}

func (b *BorgClient) CreateRepository(name string, passphrase string) (interfaces.BorgRepository, error) {
	path := b.outputDir + "/" + name
	b.container.SetEnv("BORG_PASSPHRASE", passphrase)

	_, err := b.container.Exec("borg init " + path + " -e repokey-blake2 --make-parent-dirs")
	if err != nil {
		err = b.handleError(err)
		return nil, err
	}

	return NewBorgRepository(b, path, passphrase)
}

func (b *BorgClient) GetOrCreateRepository(name string, passphrase string) (interfaces.BorgRepository, error) {
	repo, err := b.GetRepository(name, passphrase)
	if err != nil {
		if _, ok := err.(*errors.RepositoryDoesNotExistError); ok {
			repo, err = b.CreateRepository(name, passphrase)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return repo, nil
}

func (b *BorgClient) SetInputDir(inputDir string) {
	b.inputDir = inputDir
}

func (b *BorgClient) SetOutputDir(outputDir string) {
	b.outputDir = outputDir
}
