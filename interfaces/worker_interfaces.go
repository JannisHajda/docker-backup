package interfaces

type Worker interface {
	Backup() error
	Stop() error
}

type Backup struct {
	Path       string
	Passphrase string
	Keyfile    string
}

type LocalBackup struct {
	Backup
	VolumeName string
}

type RemoteBackup struct {
	Backup
	User   string
	Host   string
	SSHKey string
}
