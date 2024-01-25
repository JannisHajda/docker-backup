package db

import "docker-backup/interfaces"

type DatabaseContainer struct {
	client interfaces.DatabaseClient
	id     string
	name   string
}

func NewDatabaseContainer(client interfaces.DatabaseClient, id string, name string) interfaces.DatabaseContainer {
	return &DatabaseContainer{client: client, id: id, name: name}
}

func (d *DatabaseContainer) GetID() string {
	return d.id
}

func (d *DatabaseContainer) GetName() string {
	return d.name
}

func (d *DatabaseContainer) GetVolumes() ([]interfaces.DatabaseVolume, error) {
	cvt := d.client.GetContainerVolumesTable()
	return cvt.GetAll(d)
}

func (d *DatabaseContainer) GetVolumeByID(id int64) (interfaces.DatabaseVolume, error) {
	cvt := d.client.GetContainerVolumesTable()
	return cvt.GetByID(d, id)
}

func (d *DatabaseContainer) GetVolumeByName(name string) (interfaces.DatabaseVolume, error) {
	cvt := d.client.GetContainerVolumesTable()
	return cvt.GetByName(d, name)
}

func (d *DatabaseContainer) AddVolume(v interfaces.DatabaseVolume) error {
	cvt := d.client.GetContainerVolumesTable()
	return cvt.Add(d, v)
}

func (d *DatabaseContainer) RemoveVolume(v interfaces.DatabaseVolume) error {
	cvt := d.client.GetContainerVolumesTable()
	return cvt.Remove(d, v)
}
