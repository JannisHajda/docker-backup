package interfaces

type SSHClient interface {
	AddKeyfile(path string) error
	AddKnownHost(host string) error
}
