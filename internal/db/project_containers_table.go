package db

import (
	"docker-backup/errors"
	"docker-backup/interfaces"
)

type ProjectContainersTable struct {
	client interfaces.DatabaseClient
}

func NewProjectContainersTable(client interfaces.DatabaseClient) (interfaces.ProjectContainersTable, error) {
	pct := &ProjectContainersTable{client: client}

	err := pct.Init()
	if err != nil {
		return nil, err
	}

	return pct, nil
}

func (pct *ProjectContainersTable) Init() error {
	sql := `CREATE TABLE IF NOT EXISTS project_containers (project_id INTEGER NOT NULL, container_id VARCHAR(64) NOT NULL, PRIMARY KEY (project_id, container_id), FOREIGN KEY (project_id) REFERENCES projects (id), FOREIGN KEY (container_id) REFERENCES containers (id));`
	_, err := pct.client.Exec(sql)

	if err != nil {
		return err
	}

	return nil
}

func (pct *ProjectContainersTable) GetByName(p interfaces.DatabaseProject, container_name string) (interfaces.DatabaseContainer, error) {
	sql := `SELECT id, name FROM containers WHERE id IN (SELECT container_id FROM project_containers WHERE project_id = $1 AND container_id = (SELECT id FROM containers WHERE name = $2));`
	result, err := pct.client.Query(sql, p.GetID(), container_name)
	if err != nil {
		return nil, err
	}

	if !result.Next() {
		return nil, nil
	}

	var containerId string
	var containerName string
	err = result.Scan(&containerId, &containerName)

	if err != nil {
		return nil, err
	}

	return NewDatabaseContainer(pct.client, containerId, containerName), nil
}

func (pct *ProjectContainersTable) GetByID(p interfaces.DatabaseProject, container_id string) (interfaces.DatabaseContainer, error) {
	sql := `SELECT id, name FROM containers WHERE id IN (SELECT container_id FROM project_containers WHERE project_id = $1 AND container_id = $2);`
	result, err := pct.client.Query(sql, p.GetID(), container_id)
	if err != nil {
		return nil, err
	}

	if !result.Next() {
		return nil, nil
	}

	var containerId string
	var containerName string
	err = result.Scan(&containerId, &containerName)

	if err != nil {
		return nil, err
	}

	return NewDatabaseContainer(pct.client, containerId, containerName), nil
}

func (pct *ProjectContainersTable) GetAll(p interfaces.DatabaseProject) ([]interfaces.DatabaseContainer, error) {
	sql := `SELECT id, name FROM containers WHERE id IN (SELECT container_id FROM project_containers WHERE project_id = $1);`
	result, err := pct.client.Query(sql, p.GetID())
	if err != nil {
		return nil, err
	}

	var containers []interfaces.DatabaseContainer
	for result.Next() {
		var containerId string
		var containerName string

		err = result.Scan(&containerId, &containerName)
		if err != nil {
			return nil, err
		}

		c := NewDatabaseContainer(pct.client, containerId, containerName)
		containers = append(containers, c)
	}

	return containers, nil
}

func (pct *ProjectContainersTable) Add(p interfaces.DatabaseProject, c interfaces.DatabaseContainer) error {
	sql := "INSERT INTO project_containers (project_id, container_id) VALUES ($1, $2);"
	_, err := pct.client.Exec(sql, p.GetID(), c.GetID())
	if err != nil {
		if _, ok := err.(*errors.DuplicateKeyError); ok {
			return errors.NewContainerAlreadyInProjectError(err)
		}

		return err
	}

	return nil
}

func (pct *ProjectContainersTable) Remove(p interfaces.DatabaseProject, c interfaces.DatabaseContainer) error {
	sql := "DELETE FROM project_containers WHERE project_id = $1 AND container_id = $2;"
	_, err := pct.client.Exec(sql, p.GetID(), c.GetID())
	if err != nil {
		return err
	}

	return nil
}
