package docker

import "docker-backup/interfaces"

type DockerBind struct {
	HostPath   string
	MountPoint string
	RW         bool
}

func NewDockerBind(hostPath string, mountPoint string, rw bool) interfaces.DockerBind {
	return &DockerBind{HostPath: hostPath, MountPoint: mountPoint, RW: rw}
}

func (d *DockerBind) GetHostPath() string {
	return d.HostPath
}

func (d *DockerBind) GetMountPoint() string {
	return d.MountPoint
}

func (d *DockerBind) IsRW() bool {
	return d.RW
}
