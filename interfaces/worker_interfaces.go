package interfaces

type Worker interface {
	Backup() error
	Stop() error
}

type LocalBackup interface {
	SetVolumeName(name string)
	GetVolumeName() string
	SetVolume(volume DockerVolume)
	GetVolume() DockerVolume
}

type RemoteBackup interface {
	GetUser() string
	GetHost() string
	GetPath() string
	GetSshKey() string
}
