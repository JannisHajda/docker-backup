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
	projects []*tables.Project
}

func Connect(driver drivers.Driver) (*Database, error) {
	conn, err := sql.Open(driver.GetName(), driver.GetConnectionString())

	if err != nil {
		return nil, err
	}

	return &Database{conn: conn, driver: driver}, nil
}

func (db *Database) InitTables() error {
	err := tables.InitProjectsTable(db.conn)

	if err != nil {
		return err
	}

	db.projects = []*tables.Project{}

	return nil
}

func (db *Database) Close() error {
	return db.conn.Close()
}
