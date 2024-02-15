package errors

import "fmt"

func HandleBorgClientError(e error) error {
	if IsErrOfKind(e, REPOSITORY_DOES_NOT_EXIST) {
		return NewRepositoryDoesNotExistError(e)
	}

	if IsErrOfKind(e, REPOSITORY_PARENT_DIRECTORY_DOES_NOT_EXIST) {
		return NewRepositoryParentDirectoryDoesNotExistError(e)
	}

	if IsErrOfKind(e, REPOSITORY_ALREADY_EXISTS) {
		return NewRepositoryAlreadyExistsError(e)
	}

	if IsErrOfKind(e, WRONG_PASSPHRASE) {
		return NewWrongPassphraseError(e)
	}

	if IsErrOfKind(e, BORG_PERMISSION_DENIED) {
		return NewBorgPermissionDeniedError(e)
	}

	return fmt.Errorf("unknown error")
}

type BorgUnknownEncryptionTypeError struct {
	*InternalError
}

func NewBorgUnknownEncryptionTypeError(encryptionType string) *BorgUnknownEncryptionTypeError {
	return &BorgUnknownEncryptionTypeError{
		InternalError: NewInternalError(fmt.Sprintf("Unknown encryption type: %s", encryptionType), nil),
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

const REPOSITORY_DOES_NOT_EXIST = `Repository (.*) does not exist.`

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

const REPOSITORY_PARENT_DIRECTORY_DOES_NOT_EXIST = `parent path of the repo directory (.*) does not exist`

func NewRepositoryParentDirectoryDoesNotExistError(e error) *RepositoryParentDirectoryDoesNotExistError {
	return &RepositoryParentDirectoryDoesNotExistError{
		InternalError: NewInternalError(fmt.Sprintf("Repository parent directory does not exist"), e),
	}
}

const REPOSITORY_ALREADY_EXISTS = `repository already exists at`

type RepositoryAlreadyExistsError struct {
	*InternalError
}

func NewRepositoryAlreadyExistsError(e error) *RepositoryAlreadyExistsError {
	return &RepositoryAlreadyExistsError{
		InternalError: NewInternalError(fmt.Sprintf("Repository already exists"), e),
	}
}

const WRONG_PASSPHRASE = `passphrase supplied in BORG_PASSPHRASE, by BORG_PASSCOMMAND or via BORG_PASSPHRASE_FD is incorrect`

type WrongPassphraseError struct {
	*InternalError
}

func NewWrongPassphraseError(e error) *WrongPassphraseError {
	return &WrongPassphraseError{
		InternalError: NewInternalError(fmt.Sprintf("Wrong passphrase"), e),
	}
}

const BORG_PERMISSION_DENIED = `Permission denied`

type BorgPermissionDeniedError struct {
	*InternalError
}

func NewBorgPermissionDeniedError(e error) *BorgPermissionDeniedError {
	return &BorgPermissionDeniedError{
		InternalError: NewInternalError(fmt.Sprintf("Borg permission denied"), e),
	}
}
