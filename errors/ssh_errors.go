package errors

type RemoteHostNotFoundError struct {
	*InternalError
}

func NewRemoteHostNotFoundError(e error) *RemoteHostNotFoundError {
	return &RemoteHostNotFoundError{
		InternalError: NewInternalError("Remote host not found error", e),
	}
}

type HostKeyVerificationFailedError struct {
	*InternalError
}

func NewHostKeyVerificationFailedError(e error) *HostKeyVerificationFailedError {
	return &HostKeyVerificationFailedError{
		InternalError: NewInternalError("Host key verification failed error", e),
	}
}
