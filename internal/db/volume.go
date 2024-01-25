package db

import "docker-backup/interfaces"

type DatabaseVolume struct {
	interfaces.DatabaseClient
	name string
	id   int64
}

func NewDatabaseVolume(client interfaces.DatabaseClient, id int64, name string) interfaces.DatabaseVolume {
	return &DatabaseVolume{DatabaseClient: client, id: id, name: name}
}

func (d *DatabaseVolume) GetName() string {
	return d.name
}

func (d *DatabaseVolume) GetID() int64 {
	return d.id
}
