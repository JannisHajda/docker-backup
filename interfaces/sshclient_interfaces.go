package interfaces

type SSHClient interface {
	AddKey(path string) error
	AddKnownHost(host string) error
}
