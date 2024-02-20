package tests

import (
	"docker-backup/errors"
	"docker-backup/interfaces"
	"docker-backup/internal/borg"
	"docker-backup/mocks"
	goerrors "errors"
	"testing"
)

func GetMockContainer() *mocks.DockerContainer {
	container := mocks.DockerContainer{}
	container.On("Exec", "which borg").Return("/usr/bin/borg", nil)

	container.On("SetEnv", "BORG_REPO", "/repo").Return()
	container.On("SetEnv", "BORG_PASSPHRASE", "passphrase").Return()
	container.On("SetEnv", "BORG_KEY_FILE", "keyfile").Return()
	container.On("RemoveEnv", "BORG_REPO").Return()
	container.On("RemoveEnv", "BORG_PASSPHRASE").Return()
	container.On("RemoveEnv", "BORG_KEY_FILE").Return()

	return &container
}

func TestCreateBorgRepoInvalidEncryption(t *testing.T) {
	client := &borg.BorgClient{}
	_, err := client.CreateRepo(interfaces.CreateBorgRepoConfig{
		EncryptionType: "invalid",
	})

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if _, ok := err.(*errors.BorgUnknownEncryptionTypeError); !ok {
		t.Errorf("Expected error of type BorgUnknownEncryptionTypeError, got %T", err)
	}
}

func TestCreateBorgArchiveInvalidCompression(t *testing.T) {
	repo := &borg.BorgRepo{}
	err := repo.CreateArchive(interfaces.CreateBorgArchiveConfig{
		Compression: "invalid",
	})

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if _, ok := err.(*errors.BorgUnknownCompressionTypeError); !ok {
		t.Errorf("Expected error of type BorgUnknownCompressionTypeError, got %T", err)
	}
}

func TestBorgClientNotInstalled(t *testing.T) {
	container := mocks.DockerContainer{}
	container.On("Exec", "which borg").Return("", goerrors.New("exit status 1"))

	_, err := borg.NewBorgClient(&container)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	var borgNotInstalledError *errors.BorgNotInstalledError
	if !goerrors.As(err, &borgNotInstalledError) {
		t.Errorf("Expected error of type BorgNotInstalledError, got %T", err)
	}
}

func TestBorgClientGetNonExistingRepo(t *testing.T) {
	container := GetMockContainer()
	container.On("Exec", "borg info /repo").Return("", goerrors.New("Repository does not exist"))

	client, _ := borg.NewBorgClient(container)
	repo, err := client.GetRepo(interfaces.GetBorgRepoConfig{
		Path: "/repo",
	})

	if err == nil {
		t.Error("Expected error, got nil")
	}

	var repoDoesNotExist *errors.RepositoryDoesNotExistError
	if !goerrors.As(err, &repoDoesNotExist) {
		t.Errorf("Expected error of type BorgRepoDoesNotExistError, got %T", err)
	}

	if repo != nil {
		t.Errorf("Expected nil, got %T", repo)
	}
}

func TestBorgClientGetRepoWrongPassphrase(t *testing.T) {
	container := GetMockContainer()
	container.On("Exec", "borg info /repo").Return("", goerrors.New("Wrong passphrase"))

	client, _ := borg.NewBorgClient(container)
	repo, err := client.GetRepo(interfaces.GetBorgRepoConfig{
		Path:       "/repo",
		Passphrase: "passphrase",
	})

	if err == nil {
		t.Error("Expected error, got nil")
	}

	var wrongPassphrase *errors.WrongPassphraseError
	if !goerrors.As(err, &wrongPassphrase) {
		t.Errorf("Expected error of type WrongPassphraseError, got %T", err)
	}

	if repo != nil {
		t.Errorf("Expected nil, got %T", repo)
	}
}
