package tables

import (
	"database/sql"

	"github.com/JannisHajda/docker-backup/internal/db/drivers"
)

type ProjectsTable struct {
	conn   *sql.DB
	driver drivers.Driver
}

type Project struct {
	Id   int64
	Name string
}

type ProjectAlreadyExistsError struct {
	Name string
	Err  error
}

func (pae ProjectAlreadyExistsError) Error() string {
	return "Project with name " + pae.Name + " already exists"
}

func InitProjectsTable(conn *sql.DB, driver drivers.Driver) error {
	if driver.GetName() == "sqlite3" {
		_, err := conn.Exec(`
			CREATE TABLE IF NOT EXISTS projects (
				id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
				name TEXT NOT NULL UNIQUE
			);
		`)

		return err
	}

	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS projects (
			id SERIAL PRIMARY KEY NOT NULL,
			name TEXT NOT NULL UNIQUE
		);
	`)

	return err
}

func GetProjectsTable(conn *sql.DB, driver drivers.Driver) *ProjectsTable {
	return &ProjectsTable{conn: conn, driver: driver}
}

func (pt *ProjectsTable) Add(name string) (*Project, error) {
	project := &Project{Name: name}

	err := pt.conn.QueryRow(`
		INSERT INTO projects (name)
		VALUES ($1)
		RETURNING id;
	`, name).Scan(&project.Id)

	if err != nil {
		if pt.driver.GetName() == "sqlite3" && err.Error() == "UNIQUE constraint failed: projects.name" {
			return nil, ProjectAlreadyExistsError{Name: name, Err: err}
		}

		if pt.driver.GetName() == "postgres" && err.Error() == "pq: duplicate key value violates unique constraint \"projects_name_key\"" {
			return nil, ProjectAlreadyExistsError{Name: name, Err: err}
		}

		return nil, err
	}

	return project, nil
}

func (pt *ProjectsTable) GetById(id int64) (*Project, error) {
	project := &Project{}

	err := pt.conn.QueryRow(`
		SELECT id, name
		FROM projects
		WHERE id = $1;
	`, id).Scan(&project.Id, &project.Name)

	if err != nil {
		return nil, err
	}

	return project, nil
}

func (pt *ProjectsTable) GetByName(name string) (*Project, error) {
	project := &Project{}

	err := pt.conn.QueryRow(`
		SELECT id, name
		FROM projects
		WHERE name = $1;
	`, name).Scan(&project.Id, &project.Name)

	if err != nil {
		return nil, err
	}

	return project, nil
}
