package db

import (
	"docker-backup/errors"
	"docker-backup/interfaces"
)

type ContainersTable struct {
	client interfaces.DatabaseClient
}

func NewContainersTable(client interfaces.DatabaseClient) (interfaces.ContainersTable, error) {
	ct := &ContainersTable{client: client}

	err := ct.Init()
	if err != nil {
		return nil, err
	}

	return ct, nil
}

func (ct *ContainersTable) Init() error {
	sql := "CREATE TABLE IF NOT EXISTS containers (id VARCHAR(64) PRIMARY KEY, name VARCHAR(255) NOT NULL);"
	_, err := ct.client.Exec(sql)

	if err != nil {
		return err
	}

	return nil
}

func (ct *ContainersTable) GetByName(name string) (interfaces.DatabaseContainer, error) {
	sql := "SELECT id, name FROM containers WHERE name = $1;"
	result, err := ct.client.Query(sql, name)
	if err != nil {
		return nil, err
	}

	if !result.Next() {
		return nil, errors.NewContainerNotFoundError(nil)
	}

	var id string
	err = result.Scan(&id, &name)
	if err != nil {
		return nil, err
	}

	return NewDatabaseContainer(ct.client, id, name), nil
}

func (ct *ContainersTable) GetByID(id string) (interfaces.DatabaseContainer, error) {
	sql := "SELECT id, name FROM containers WHERE id = $1;"
	result, err := ct.client.Query(sql, id)
	if err != nil {
		return nil, err
	}

	if !result.Next() {
		return nil, nil
	}

	var name string
	err = result.Scan(&id, &name)
	if err != nil {
		return nil, err
	}

	return NewDatabaseContainer(ct.client, id, name), nil
}

func (ct *ContainersTable) GetAll() ([]interfaces.DatabaseContainer, error) {
	sql := "SELECT id, name FROM containers;"
	rows, err := ct.client.Query(sql)
	if err != nil {
		return nil, err
	}

	var containers []interfaces.DatabaseContainer

	for rows.Next() {
		var id string
		var name string

		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}

		containers = append(containers, NewDatabaseContainer(ct.client, id, name))
	}

	return containers, nil
}

func (ct *ContainersTable) Add(c interfaces.DatabaseContainer) (interfaces.DatabaseContainer, error) {
	sql := "INSERT INTO containers (id, name) VALUES ($1, $2);"
	_, err := ct.client.Exec(sql, c.GetID(), c.GetName())
	if err != nil {
		if _, ok := err.(*errors.DuplicateKeyError); ok {
			return nil, errors.NewContainerAlreadyExistsError(err)
		}

		return nil, err
	}

	c, err = ct.GetByID(c.GetID())
	if err != nil {
		return nil, err
	}

	return NewDatabaseContainer(ct.client, c.GetID(), c.GetName()), nil
}

func (ct *ContainersTable) Remove(c interfaces.DatabaseContainer) error {
	sql := "DELETE FROM containers WHERE id = $1;"
	_, err := ct.client.Exec(sql, c.GetID())
	if err != nil {
		return err
	}

	return nil
}
