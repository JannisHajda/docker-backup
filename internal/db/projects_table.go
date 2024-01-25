package db

import (
	"docker-backup/errors"
	"docker-backup/interfaces"
)

type ProjectsTable struct {
	client interfaces.DatabaseClient
}

func NewProjectTable(client interfaces.DatabaseClient) (interfaces.ProjectsTable, error) {
	pt := &ProjectsTable{client: client}

	err := pt.Init()
	if err != nil {
		return nil, err
	}

	return pt, nil
}

func (pt *ProjectsTable) GetClient() interfaces.DatabaseClient {
	return pt.client
}

func (pt *ProjectsTable) SetClient(client interfaces.DatabaseClient) {
	pt.client = client
}

func (pt *ProjectsTable) Init() error {
	sql := "CREATE TABLE IF NOT EXISTS projects (id SERIAL PRIMARY KEY, name VARCHAR(255) NOT NULL);"
	_, err := pt.client.Exec(sql)

	if err != nil {
		return err
	}

	return nil
}

func (pt *ProjectsTable) GetByName(name string) (interfaces.DatabaseProject, error) {
	sql := "SELECT id, name FROM projects WHERE name = $1;"
	result, err := pt.client.Query(sql, name)
	if err != nil {
		return nil, err
	}

	if !result.Next() {
		return nil, nil
	}

	var id int64
	err = result.Scan(&id, &name)
	if err != nil {
		if itemNotFoundError, ok := err.(*errors.ItemNotFoundError); ok {
			return nil, errors.NewProjectAlreadyExistsError(itemNotFoundError)
		}

		return nil, err
	}

	return NewDatabaseProject(pt.client, id, name), nil
}

func (pt *ProjectsTable) GetByID(id int64) (interfaces.DatabaseProject, error) {
	sql := "SELECT id, name FROM projects WHERE id = $1;"
	result, err := pt.client.Query(sql, id)
	if err != nil {
		if itemNotFoundError, ok := err.(*errors.ItemNotFoundError); ok {
			return nil, errors.NewProjectNotFoundError(itemNotFoundError)
		}

		return nil, err
	}

	if !result.Next() {
		return nil, nil
	}

	var name string
	err = result.Scan(&id, &name)
	if err != nil {
		if itemNotFoundError, ok := err.(*errors.ItemNotFoundError); ok {
			return nil, errors.NewProjectAlreadyExistsError(itemNotFoundError)
		}

		return nil, err
	}

	return NewDatabaseProject(pt.client, id, name), nil
}

func (pt *ProjectsTable) GetAll() ([]interfaces.DatabaseProject, error) {
	sql := "SELECT id, name FROM projects;"
	rows, err := pt.client.Query(sql)
	if err != nil {
		return nil, err
	}

	var projects []interfaces.DatabaseProject
	for rows.Next() {
		var id int64
		var name string

		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}

		project := NewDatabaseProject(pt.client, id, name)
		if err != nil {
			return nil, err
		}

		projects = append(projects, project)
	}

	return projects, nil
}

func (pt *ProjectsTable) Add(p interfaces.DatabaseProject) (interfaces.DatabaseProject, error) {
	sql := "INSERT INTO projects (name) VALUES ($1);"
	_, err := pt.client.Exec(sql, p.GetName())
	if err != nil {
		if itemAlreadyExistsError, ok := err.(*errors.DuplicateKeyError); ok {
			return nil, errors.NewProjectAlreadyExistsError(itemAlreadyExistsError)
		}

		return nil, err
	}

	p, err = pt.GetByName(p.GetName())
	if err != nil {
		return nil, err
	}

	return NewDatabaseProject(pt.client, p.GetID(), p.GetName()), nil
}

func (pt *ProjectsTable) Remove(p interfaces.DatabaseProject) error {
	sql := "DELETE FROM projects WHERE id = $1;"
	_, err := pt.client.Exec(sql, p.GetID())
	if err != nil {
		if itemNotFoundError, ok := err.(*errors.ItemNotFoundError); ok {
			return errors.NewProjectNotFoundError(itemNotFoundError)
		}

		return err
	}

	return nil
}
