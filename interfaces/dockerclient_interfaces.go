package interfaces

type DockerClient interface {
	GetContainer(id string) (DockerContainer, error)
	CreateContainer(image string, volumes []DockerVolume) (DockerContainer, error)
}

type DockerContainer interface {
	GetID() string
	GetName() string
	Exec(cmd string) (string, error)
	GetEnv(key string) string
	GetEnvs() map[string]string
	SetEnv(key string, value string)
	Start() error
	Stop() error
	Remove() error
	StopAndRemove() error
	GetVolumes() []DockerVolume
}

type DockerVolume interface {
	GetName() string
	GetMountPoint() string
	SetMountPoint(mountPoint string)
	IsRW() bool
}
