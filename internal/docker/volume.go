package docker

import "docker-backup/interfaces"

type DockerVolume struct {
	Name       string
	MountPoint string
	RW         bool
}

func NewDockerVolume(name string, mountPoint string, rw bool) interfaces.DockerVolume {
	return &DockerVolume{Name: name, MountPoint: mountPoint, RW: rw}
}

func (d *DockerVolume) GetName() string {
	return d.Name
}

func (d *DockerVolume) GetMountPoint() string {
	return d.MountPoint
}

func (d *DockerVolume) SetMountPoint(mountPoint string) {
	d.MountPoint = mountPoint
}

func (d *DockerVolume) IsRW() bool {
	return d.RW
}
