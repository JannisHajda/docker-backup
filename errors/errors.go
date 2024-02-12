package errors

import (
	"fmt"
	"regexp"
)

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

// detect regex errors
func IsErrOfKind(err error, kind string) bool {
	message := err.Error()
	re := regexp.MustCompile(kind)
	matches := re.FindStringSubmatch(message)

	if len(matches) > 0 {
		return true
	}

	return false
}
