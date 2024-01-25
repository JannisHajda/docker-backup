package errors

import "fmt"

type InternalError struct {
	Message string
	Cause   error
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("Internal error: %s", e.Message)
}

func NewInternalError(message string, cause error) *InternalError {
	return &InternalError{
		Message: message,
		Cause:   cause,
	}
}
