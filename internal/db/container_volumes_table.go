package db

import (
	"docker-backup/errors"
	"docker-backup/interfaces"
)

type ContainerVolumesTable struct {
	client interfaces.DatabaseClient
}

func NewContainerVolumesTable(client interfaces.DatabaseClient) (interfaces.ContainerVolumesTable, error) {
	cvt := &ContainerVolumesTable{client: client}

	err := cvt.Init()
	if err != nil {
		return nil, err
	}

	return cvt, nil
}

func (cvt *ContainerVolumesTable) Init() error {
	sql := "CREATE TABLE IF NOT EXISTS container_volumes (container_id VARCHAR(64) REFERENCES containers(id) ON DELETE CASCADE, volume_id INTEGER REFERENCES volumes(id) ON DELETE CASCADE);"
	_, err := cvt.client.Exec(sql)

	if err != nil {
		return err
	}

	return nil
}

func (cvt *ContainerVolumesTable) GetByName(c interfaces.DatabaseContainer, volume_name string) (interfaces.DatabaseVolume, error) {
	sql := `SELECT id, name FROM volumes WHERE id IN (SELECT volume_id FROM container_volumes WHERE container_id = $1 AND volume_id = (SELECT id FROM volumes WHERE name = $2));`
	result, err := cvt.client.Query(sql, c.GetID(), volume_name)
	if err != nil {
		return nil, err
	}

	if !result.Next() {
		return nil, nil
	}

	var volumeId int64
	var volumeName string
	err = result.Scan(&volumeId, &volumeName)
	if err != nil {
		return nil, err
	}

	return NewDatabaseVolume(cvt.client, volumeId, volumeName), nil
}

func (cvt *ContainerVolumesTable) GetByID(c interfaces.DatabaseContainer, volume_id int64) (interfaces.DatabaseVolume, error) {
	sql := `SELECT id, name FROM volumes WHERE id IN (SELECT volume_id FROM container_volumes WHERE container_id = $1 AND volume_id = $2);`
	result, err := cvt.client.Query(sql, c.GetID(), volume_id)
	if err != nil {
		return nil, err
	}

	if !result.Next() {
		return nil, nil
	}

	var volumeId int64
	var volumeName string
	err = result.Scan(&volumeId, &volumeName)
	if err != nil {
		return nil, err
	}

	return NewDatabaseVolume(cvt.client, volumeId, volumeName), nil
}

func (cvt *ContainerVolumesTable) GetAll(c interfaces.DatabaseContainer) ([]interfaces.DatabaseVolume, error) {
	sql := `SELECT id, name FROM volumes WHERE id IN (SELECT volume_id FROM container_volumes WHERE container_id = $1);`
	result, err := cvt.client.Query(sql, c.GetID())
	if err != nil {
		return nil, err
	}

	var volumes []interfaces.DatabaseVolume
	for result.Next() {
		var volumeId int64
		var volumeName string
		err = result.Scan(&volumeId, &volumeName)
		if err != nil {
			return nil, err
		}

		volumes = append(volumes, NewDatabaseVolume(cvt.client, volumeId, volumeName))
	}

	return volumes, nil
}

func (cvt *ContainerVolumesTable) Add(c interfaces.DatabaseContainer, v interfaces.DatabaseVolume) error {
	sql := "INSERT INTO container_volumes (container_id, volume_id) VALUES ($1, $2);"
	_, err := cvt.client.Exec(sql, c.GetID(), v.GetID())
	if err != nil {
		if _, ok := err.(*errors.DuplicateKeyError); ok {
			return errors.NewVolumeAlreadyInContainerError(err)
		}

		return err
	}

	return nil
}

func (cvt *ContainerVolumesTable) Remove(c interfaces.DatabaseContainer, v interfaces.DatabaseVolume) error {
	sql := "DELETE FROM container_volumes WHERE container_id = $1 AND volume_id = $2;"
	_, err := cvt.client.Exec(sql, c.GetID(), v.GetID())
	if err != nil {
		return err
	}

	return nil
}
