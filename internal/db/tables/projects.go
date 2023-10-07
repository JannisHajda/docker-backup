package tables

import "database/sql"

type Project struct {
	Id   int64
	Name string
}

func InitProjectsTable(conn *sql.DB) error {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS projects (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL
		);
	`)

	return err
}
