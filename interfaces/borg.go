package interfaces

type BorgClient interface {
	GetRepository(path string, passphrase string) (BorgRepository, error)
	CreateRepository(path string, passphrase string) (BorgRepository, error)
	GetOrCreateRepository(path string, passphrase string) (BorgRepository, error)
	GetContainer() DockerContainer
}

type BorgRepository interface {
	Backup(input string) error
	GetPath() string
}
