package worker

import "docker-backup/interfaces"

type LocalBackup struct {
	volumeName string
	v          interfaces.DockerVolume
}

type RemoteBackup struct {
	user    string
	host    string
	path    string
	ssh_key string
}

func NewLocalBackup(volumeName string) interfaces.LocalBackup {
	return &LocalBackup{volumeName: volumeName}
}

func (l *LocalBackup) SetVolume(volume interfaces.DockerVolume) {
	l.v = volume
}

func (l *LocalBackup) SetVolumeName(volumeName string) {
	l.volumeName = volumeName
}

func (l *LocalBackup) GetVolumeName() string {
	return l.volumeName
}

func (l *LocalBackup) GetVolume() interfaces.DockerVolume {
	return l.v
}

func NewRemoteBackup(user string, host string, path string, ssh_key string) interfaces.RemoteBackup {
	return &RemoteBackup{user: user, host: host, path: path, ssh_key: ssh_key}
}

func (r *RemoteBackup) GetUser() string {
	return r.user
}

func (r *RemoteBackup) GetHost() string {
	return r.host
}

func (r *RemoteBackup) GetPath() string {
	return r.path
}

func (r *RemoteBackup) GetSshKey() string {
	return r.ssh_key
}
