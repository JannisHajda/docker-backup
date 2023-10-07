package db

type ProjectContainersTable struct {
	db                *Database
	projectContainers []*ProjectContainers
}

type ProjectContainers struct {
	ProjectId   int64
	ContainerId string
}

func (db *Database) InitProjectContainersTable() error {
	db.pct = &ProjectContainersTable{db: db, projectContainers: []*ProjectContainers{}}

	if db.driver.GetName() == "sqlite3" {
		_, err := db.conn.Exec(`
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

	_, err := db.conn.Exec(`
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
