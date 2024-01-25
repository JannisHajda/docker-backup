package errors

import "fmt"

type DuplicateKeyError struct {
	*InternalError
}

func NewDuplicateKeyError(e error) *DuplicateKeyError {
	return &DuplicateKeyError{
		InternalError: NewInternalError(fmt.Sprintf("Duplicate key error"), e),
	}
}

type ItemNotFoundError struct {
	*InternalError
}

func NewItemNotFoundError(e error) *ItemNotFoundError {
	return &ItemNotFoundError{
		InternalError: NewInternalError(fmt.Sprintf("Item not found error"), e),
	}
}

type ProjectNotFoundError struct {
	*InternalError
}

func NewProjectNotFoundError(e error) *ProjectNotFoundError {
	return &ProjectNotFoundError{
		InternalError: NewInternalError(fmt.Sprintf("Project not found error"), e),
	}
}

type ProjectAlreadyExistsError struct {
	*InternalError
}

func NewProjectAlreadyExistsError(e error) *ProjectAlreadyExistsError {
	return &ProjectAlreadyExistsError{
		InternalError: NewInternalError(fmt.Sprintf("Project already exists error"), e),
	}
}

type ContainerNotFoundError struct {
	*InternalError
}

func NewContainerNotFoundError(e error) *ContainerNotFoundError {
	return &ContainerNotFoundError{
		InternalError: NewInternalError(fmt.Sprintf("Container not found error"), e),
	}
}

type ContainerAlreadyExistsError struct {
	*InternalError
}

func NewContainerAlreadyExistsError(e error) *ContainerAlreadyExistsError {
	return &ContainerAlreadyExistsError{
		InternalError: NewInternalError(fmt.Sprintf("Container already exists error"), e),
	}
}

type ContainerAlreadyInProjectError struct {
	*InternalError
}

func NewContainerAlreadyInProjectError(e error) *ContainerAlreadyInProjectError {
	return &ContainerAlreadyInProjectError{
		InternalError: NewInternalError(fmt.Sprintf("Container already in project error"), e),
	}
}

type VolumeNotFoundError struct {
	*InternalError
}

func NewVolumeNotFoundError(e error) *VolumeNotFoundError {
	return &VolumeNotFoundError{
		InternalError: NewInternalError(fmt.Sprintf("Volume not found error"), e),
	}
}

type VolumeAlreadyExistsError struct {
	*InternalError
}

func NewVolumeAlreadyExistsError(e error) *VolumeAlreadyExistsError {
	return &VolumeAlreadyExistsError{
		InternalError: NewInternalError(fmt.Sprintf("Volume already exists error"), e),
	}
}

type VolumeAlreadyInContainerError struct {
	*InternalError
}

func NewVolumeAlreadyInContainerError(e error) *VolumeAlreadyInContainerError {
	return &VolumeAlreadyInContainerError{
		InternalError: NewInternalError(fmt.Sprintf("Volume already in container error"), e),
	}
}

type VolumeNotInContainerError struct {
	*InternalError
}

func NewVolumeNotInContainerError(e error) *VolumeNotInContainerError {
	return &VolumeNotInContainerError{
		InternalError: NewInternalError(fmt.Sprintf("Volume not in container error"), e),
	}
}
