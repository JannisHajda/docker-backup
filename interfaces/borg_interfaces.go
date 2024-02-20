package interfaces

type CreateBorgRepoConfig struct {
	Path           string
	EncryptionType string
	Passphrase     string
	Keyfile        string
	AppendOnly     bool
	StorageQuota   string
	MakeParentDirs bool
}

type GetBorgRepoConfig struct {
	Path       string
	Passphrase string
	Keyfile    string
}

type BorgClient interface {
	GetRepo(config GetBorgRepoConfig) (BorgRepo, error)
	CreateRepo(config CreateBorgRepoConfig) (BorgRepo, error)
}

type CreateBorgArchiveConfig struct {
	Sources     []string
	Name        string
	Compression string
}

type BorgRepo interface {
	Info() (string, error)
	ListArchives() (string, error)
	CreateArchive(config CreateBorgArchiveConfig) error
}
