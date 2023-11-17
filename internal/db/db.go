package db

import (
	"database/sql"
	"strings"

	"github.com/JannisHajda/docker-backup/internal/db/drivers"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	driver drivers.Driver
	pt     *ProjectsTable
	ct     *ContainersTable
	pct    *ProjectContainersTable
}

type SQLCommand struct {
	postgres string
	sqlite3  string
}

func (db *Database) testConnection() error {
	connection, err := db.connect()

	if err != nil {
		return err
	}

	defer connection.Close()

	return connection.Ping()
}

func NewDatabase(driver drivers.Driver) (*Database, error) {
	db := &Database{driver: driver}

	err := db.testConnection()
	if err != nil {
		return nil, err
	}

	db.pt, err = newProjectsTable(db)
	if err != nil {
		return nil, err
	}

	db.ct, err = newContainersTable(db)
	if err != nil {
		return nil, err
	}

	db.pct, err = newProjectContainersTable(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *Database) IsUniqueViolationError(err error) bool {
	return strings.Contains(err.Error(), db.driver.UniqueViolationError())
}

func (db *Database) IsNoRowsError(err error) bool {
	return strings.Contains(err.Error(), db.driver.NoRowsError())
}

func (db *Database) connect() (*sql.DB, error) {
	c, err := sql.Open(db.driver.GetName(), db.driver.GetConnectionString())

	if err != nil {
		return nil, err
	}

	return c, nil
}

func (db *Database) execute(command SQLCommand, args ...interface{}) (sql.Result, error) {
	connection, err := db.connect()

	if err != nil {
		return nil, err
	}

	defer connection.Close()

	if db.driver.GetName() == "sqlite3" {
		return connection.Exec(command.sqlite3, args...)
	}

	return connection.Exec(command.postgres, args...)
}

func (db *Database) queryRow(command SQLCommand, args ...any) (*sql.Row, error) {
	connection, err := db.connect()

	if err != nil {
		return nil, err
	}

	defer connection.Close()

	if db.driver.GetName() == "sqlite3" {
		return connection.QueryRow(command.sqlite3, args...), nil
	}

	return connection.QueryRow(command.postgres, args...), nil
}

func (db *Database) query(command SQLCommand, args ...any) (*sql.Rows, error) {
	connection, err := db.connect()

	if err != nil {
		return nil, err
	}

	defer connection.Close()

	if db.driver.GetName() == "sqlite3" {
		return connection.Query(command.sqlite3, args...)
	}

	return connection.Query(command.postgres, args...)
}

func (db *Database) GetAllProjects() ([]*Project, error) {
	return db.pt.getAll()
}

func (db *Database) GetProjectByName(name string) (*Project, error) {
	return db.pt.getByName(name)
}

func (db *Database) AddProject(name string) (*Project, error) {
	return db.pt.add(name)
}

func (db *Database) GetContainerByID(id string) (*Container, error) {
	return db.ct.getByID(id)
}

func (db *Database) GetAllContainers() ([]*Container, error) {
	return db.ct.getAll()
}

func (db *Database) AddContainer(id string, name string) (*Container, error) {
	return db.ct.add(id, name)
}

func (db *Database) GetOrAddContainer(id string, name string) (*Container, error) {
	container, err := db.GetContainerByID(id)

	if err != nil && db.IsNoRowsError(err) {
		return db.AddContainer(id, name)
	}

	if err != nil {
		return nil, err
	}

	return container, err
}
