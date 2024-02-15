package interfaces

type CreateBorgRepoConfig struct {
	Path           string
	EncryptionType string
	Passphrase     string
	Key            string
	AppendOnly     bool
	StorageQuota   string
	MakeParentDirs bool
}

type BorgClient interface {
	GetRepo(path string, passphrase string) (BorgRepo, error)
	CreateRepo(config CreateBorgRepoConfig) (BorgRepo, error)
	GetOrCreateRepo(config CreateBorgRepoConfig) (BorgRepo, error)
	GetContainer() DockerContainer
}

type BorgRepo interface {
	Backup(input string) error
	GetPath() string
}
