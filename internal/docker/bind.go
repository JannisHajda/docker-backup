package docker

import "docker-backup/interfaces"

type DockerBind struct {
	Source string
	Target string
	RW     bool
}

func NewDockerBind(source string, target string, rw bool) interfaces.DockerBind {
	return &DockerBind{Source: source, Target: target, RW: rw}
}

func (d *DockerBind) GetSource() string {
	return d.Source
}

func (d *DockerBind) GetTarget() string {
	return d.Target
}

func (d *DockerBind) IsRW() bool {
	return d.RW
}
