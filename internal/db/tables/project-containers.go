package tables

import (
	"database/sql"

	"github.com/JannisHajda/docker-backup/internal/db/drivers"
)

type ProjectContainers struct {
	ProjectId   int64
	ContainerId string
}

func InitProjectContainersTable(conn *sql.DB, driver drivers.Driver) error {
	if driver.GetName() == "sqlite3" {
		_, err := conn.Exec(`
			CREATE TABLE IF NOT EXISTS project_containers (
				project_id INTEGER NOT NULL,
				container_id TEXT NOT NULL,
				PRIMARY KEY (project_id, container_id),
				FOREIGN KEY (project_id) REFERENCES projects(id),
				FOREIGN KEY (container_id) REFERENCES containers(id)
			);
		`)

		return err
	}

	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS project_containers (
			project_id SERIAL NOT NULL,
			container_id TEXT NOT NULL,
			PRIMARY KEY (project_id, container_id),
			FOREIGN KEY (project_id) REFERENCES projects(id),
			FOREIGN KEY (container_id) REFERENCES containers(id)
		);
	`)

	return err
}
