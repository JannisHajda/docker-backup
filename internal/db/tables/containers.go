package tables

import (
	"database/sql"

	"github.com/JannisHajda/docker-backup/internal/db/drivers"
)

type Container struct {
	Id   string
	Name string
}

func InitContainersTable(conn *sql.DB, driver drivers.Driver) error {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS containers (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL
		);
	`)

	return err
}
