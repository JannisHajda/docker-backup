package errors

import (
	"fmt"
)

type BorgUnknownEncryptionTypeError struct {
	*InternalError
}

func NewBorgUnknownEncryptionTypeError(encryptionType string) *BorgUnknownEncryptionTypeError {
	return &BorgUnknownEncryptionTypeError{
		InternalError: NewInternalError(fmt.Sprintf("Unknown encryption type: %s", encryptionType), nil),
	}
}

type BorgUnknownCompressionTypeError struct {
	*InternalError
}

func NewBorgUnknownCompressionTypeError(compression string) *BorgUnknownCompressionTypeError {
	return &BorgUnknownCompressionTypeError{
		InternalError: NewInternalError(fmt.Sprintf("Unknown compression type: %s", compression), nil),
	}
}

type BorgNotInstalledError struct {
	*InternalError
}

func NewBorgNotInstalledError(e error) *BorgNotInstalledError {
	return &BorgNotInstalledError{
		InternalError: NewInternalError(fmt.Sprintf("Borg not installed"), e),
	}
}

type RepositoryDoesNotExistError struct {
	*InternalError
}

func NewRepositoryDoesNotExistError(e error) *RepositoryDoesNotExistError {
	return &RepositoryDoesNotExistError{
		InternalError: NewInternalError(fmt.Sprintf("Repository does not exist"), e),
	}
}

type RepositoryParentDirectoryDoesNotExistError struct {
	*InternalError
}

func NewRepositoryParentDirectoryDoesNotExistError(e error) *RepositoryParentDirectoryDoesNotExistError {
	return &RepositoryParentDirectoryDoesNotExistError{
		InternalError: NewInternalError(fmt.Sprintf("Repository parent directory does not exist"), e),
	}
}

type RepositoryAlreadyExistsError struct {
	*InternalError
}

func NewRepositoryAlreadyExistsError(e error) *RepositoryAlreadyExistsError {
	return &RepositoryAlreadyExistsError{
		InternalError: NewInternalError(fmt.Sprintf("Repository already exists"), e),
	}
}

type WrongPassphraseError struct {
	*InternalError
}

func NewWrongPassphraseError(e error) *WrongPassphraseError {
	return &WrongPassphraseError{
		InternalError: NewInternalError(fmt.Sprintf("Wrong passphrase"), e),
	}
}

type BorgPermissionDeniedError struct {
	*InternalError
}

func NewBorgPermissionDeniedError(e error) *BorgPermissionDeniedError {
	return &BorgPermissionDeniedError{
		InternalError: NewInternalError(fmt.Sprintf("Borg permission denied"), e),
	}
}
