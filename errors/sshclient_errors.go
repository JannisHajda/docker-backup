package errors

const REMOTE_HOST_NOT_FOOUND = "Could not resolve hostname"

type RemoteHostNotFoundError struct {
	*InternalError
}

func NewRemoteHostNotFoundError(e error) *RemoteHostNotFoundError {
	return &RemoteHostNotFoundError{
		InternalError: NewInternalError("Remote host not found error", e),
	}
}

const HOST_KEY_VERIFICATION_FAILED = "Host key verification failed"

type HostKeyVerificationFailedError struct {
	*InternalError
}

func NewHostKeyVerificationFailedError(e error) *HostKeyVerificationFailedError {
	return &HostKeyVerificationFailedError{
		InternalError: NewInternalError("Host key verification failed error", e),
	}
}
