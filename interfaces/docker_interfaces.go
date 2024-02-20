package interfaces

type DockerClient interface {
	GetContainer(id string) (DockerContainer, error)
	CreateContainer(image string, volumes []DockerVolume, binds []DockerBind) (DockerContainer, error)
	CreateVolume(name string) (DockerVolume, error)
	GetVolume(name string) (DockerVolume, error)
}

type DockerContainer interface {
	GetID() string
	GetName() string
	Exec(cmd string) (string, error)
	GetEnv(key string) string
	GetEnvs() map[string]string
	SetEnv(key string, value string)
	RemoveEnv(key string)
	Start() error
	Stop() error
	Remove() error
	StopAndRemove() error
	GetVolumes() []DockerVolume
	GetBinds() []DockerBind
}

type DockerVolume interface {
	GetName() string
	GetMountPoint() string
	SetMountPoint(mountPoint string)
	IsRW() bool
}

type DockerBind interface {
	GetHostPath() string
	GetMountPoint() string
	IsRW() bool
}
