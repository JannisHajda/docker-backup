package interfaces

type BorgClient interface {
	GetRepository(name string, passphrase string) (BorgRepository, error)
	CreateRepository(name string, passphrase string) (BorgRepository, error)
	GetOrCreateRepository(name string, passphrase string) (BorgRepository, error)
	SetInputDir(inputDir string)
	SetOutputDir(outputDir string)
}

type BorgRepository interface {
	Backup() error
	GetName() string
}
