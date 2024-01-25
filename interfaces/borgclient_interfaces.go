package interfaces

type BorgClient interface {
	GetRepository(path string, passphrase string) (BorgRepository, error)
	CreateRepository(path string, passphrase string) (BorgRepository, error)
	GetOrCreateRepository(path string, passphrase string) (BorgRepository, error)
}

type BorgRepository interface {
	Archive(inputPath string) error
}
