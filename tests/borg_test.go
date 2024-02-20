package tests

import (
	"docker-backup/errors"
	"docker-backup/interfaces"
	"docker-backup/internal/borg"
	"docker-backup/mocks"
	goerrors "errors"
	"testing"
)

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
