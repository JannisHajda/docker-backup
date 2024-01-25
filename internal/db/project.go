package db

import "docker-backup/interfaces"

type DatabaseProject struct {
	client interfaces.DatabaseClient
	id     int64
	name   string
}

func NewDatabaseProject(client interfaces.DatabaseClient, id int64, name string) interfaces.DatabaseProject {
	return &DatabaseProject{client: client, id: id, name: name}
}

func (d *DatabaseProject) GetID() int64 {
	return d.id
}

func (d *DatabaseProject) GetName() string {
	return d.name
}

func (d *DatabaseProject) GetContainers() ([]interfaces.DatabaseContainer, error) {
	pct := d.client.GetProjectContainersTable()
	return pct.GetAll(d)
}

func (d *DatabaseProject) GetContainerByID(id string) (interfaces.DatabaseContainer, error) {
	pct := d.client.GetProjectContainersTable()
	return pct.GetByID(d, id)
}

func (d *DatabaseProject) GetContainerByName(name string) (interfaces.DatabaseContainer, error) {
	pct := d.client.GetProjectContainersTable()
	return pct.GetByName(d, name)
}

func (d *DatabaseProject) AddContainer(c interfaces.DatabaseContainer) error {
	pct := d.client.GetProjectContainersTable()
	return pct.Add(d, c)
}

func (d *DatabaseProject) RemoveContainer(c interfaces.DatabaseContainer) error {
	pct := d.client.GetProjectContainersTable()
	return pct.Remove(d, c)
}
