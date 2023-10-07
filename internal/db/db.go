package db

import (
	"database/sql"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	conn   *sql.DB
	driver DbDriver
}

type DbDriver interface {
	GetName() string
	GetConnectionString() string
}

func Connect(driver DbDriver) (*Database, error) {
	conn, err := sql.Open(driver.GetName(), driver.GetConnectionString())

	if err != nil {
		return nil, err
	}

	return &Database{conn: conn, driver: driver}, nil
}

func (db *Database) Close() error {
	return db.conn.Close()
}
