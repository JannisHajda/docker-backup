package errors

import "fmt"

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
