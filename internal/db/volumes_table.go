package db

import (
	"docker-backup/errors"
	"docker-backup/interfaces"
)

type VolumesTable struct {
	client interfaces.DatabaseClient
}

func NewVolumesTable(client interfaces.DatabaseClient) (interfaces.VolumesTable, error) {
	vt := &VolumesTable{client: client}

	err := vt.Init()
	if err != nil {
		return nil, err
	}

	return vt, nil
}

func (vt *VolumesTable) Init() error {
	sql := "CREATE TABLE IF NOT EXISTS volumes (id SERIAL PRIMARY KEY, name VARCHAR(255) NOT NULL);"
	_, err := vt.client.Exec(sql)

	if err != nil {
		return err
	}

	return nil
}

func (vt *VolumesTable) GetByName(name string) (interfaces.DatabaseVolume, error) {
	sql := "SELECT id, name FROM volumes WHERE name = $1;"
	result, err := vt.client.Query(sql, name)
	if err != nil {
		return nil, err
	}

	if !result.Next() {
		return nil, errors.NewVolumeNotFoundError(nil)
	}

	var id int64
	err = result.Scan(&id, &name)
	if err != nil {
		return nil, err
	}

	return NewDatabaseVolume(vt.client, id, name), nil
}

func (vt *VolumesTable) GetByID(id int64) (interfaces.DatabaseVolume, error) {
	sql := "SELECT id, name FROM volumes WHERE id = $1;"
	result, err := vt.client.Query(sql, id)
	if err != nil {
		return nil, err
	}

	if !result.Next() {
		return nil, errors.NewVolumeNotFoundError(nil)
	}

	var name string
	err = result.Scan(&id, &name)
	if err != nil {
		return nil, err
	}

	return NewDatabaseVolume(vt.client, id, name), nil
}

func (vt *VolumesTable) GetAll() ([]interfaces.DatabaseVolume, error) {
	sql := "SELECT id, name FROM volumes;"
	rows, err := vt.client.Query(sql)
	if err != nil {
		return nil, err
	}

	var volumes []interfaces.DatabaseVolume
	for rows.Next() {
		var id int64
		var name string

		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}

		volume := NewDatabaseVolume(vt.client, id, name)
		if err != nil {
			return nil, err
		}

		volumes = append(volumes, volume)
	}

	return volumes, nil
}

func (vt *VolumesTable) Add(v interfaces.DatabaseVolume) (interfaces.DatabaseVolume, error) {
	sql := "INSERT INTO volumes (name) VALUES ($1);"
	_, err := vt.client.Exec(sql, v.GetName())
	if err != nil {
		if _, ok := err.(*errors.DuplicateKeyError); ok {
			return nil, errors.NewVolumeAlreadyExistsError(err)
		}

		return nil, err
	}

	v, err = vt.GetByName(v.GetName())
	if err != nil {
		return nil, err
	}

	return NewDatabaseVolume(vt.client, v.GetID(), v.GetName()), nil
}

func (vt *VolumesTable) Remove(v interfaces.DatabaseVolume) error {
	sql := "DELETE FROM volumes WHERE id = $1;"
	_, err := vt.client.Exec(sql, v.GetID())
	if err != nil {
		return err
	}

	return nil
}
