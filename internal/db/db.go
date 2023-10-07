package db

import (
	"database/sql"

	"github.com/JannisHajda/docker-backup/internal/db/drivers"
	"github.com/JannisHajda/docker-backup/internal/db/tables"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	conn     *sql.DB
	driver   drivers.Driver
	pt       *ProjectsTable
	projects []*Project
}

func Connect(driver drivers.Driver) (*Database, error) {
	conn, err := sql.Open(driver.GetName(), driver.GetConnectionString())

	if err != nil {
		return nil, err
	}

	return &Database{conn: conn, driver: driver}, nil
}

func (db *Database) InitTables() error {
	err := db.InitProjectsTable()

	if err != nil {
		return err
	}

	err = tables.InitContainersTable(db.conn, db.driver)

	if err != nil {
		return err
	}

	err = tables.InitProjectContainersTable(db.conn, db.driver)

	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetConnection() *sql.DB {
	return db.conn
}

func (db *Database) AddProject(name string) error {

	if len(db.projects) != 0 {
		for _, project := range db.projects {
			if project.Name == name {
				return ProjectAlreadyExistsError{Name: name, Err: nil}
			}
		}
	}

	p, err := db.pt.Add(name)

	if err != nil {
		_, ok := err.(ProjectAlreadyExistsError)

		if ok {
			p, err := db.pt.GetByName(name)

			if err != nil {
				return err
			}

			db.projects = append(db.projects, p)
			return ProjectAlreadyExistsError{Name: name, Err: err}
		}

		return err
	}

	db.projects = append(db.projects, p)
	return nil
}

func (db *Database) Close() error {
	return db.conn.Close()
}
