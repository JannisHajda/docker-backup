package db

type ProjectContainersTable struct {
	db *Database
}

func newProjectContainersTable(db *Database) (*ProjectContainersTable, error) {
	pct := &ProjectContainersTable{db: db}
	err := pct.init()
	if err != nil {
		return nil, err
	}

	return pct, nil
}

func (pct *ProjectContainersTable) init() error {
	sql := SQLCommand{
		postgres: `
			CREATE TABLE IF NOT EXISTS project_containers (
				project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
				container_id TEXT REFERENCES containers(id) ON DELETE CASCADE,
				PRIMARY KEY (project_id, container_id)
			);
		`,
		sqlite3: `
			CREATE TABLE IF NOT EXISTS project_containers (
				project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
				container_id TEXT REFERENCES containers(id) ON DELETE CASCADE,
				PRIMARY KEY (project_id, container_id)
			);
		`,
	}

	_, err := pct.db.execute(sql)
	return err
}

func (pct *ProjectContainersTable) add(projectID int, containerID string) error {
	sql := SQLCommand{
		postgres: `INSERT INTO project_containers (project_id, container_id) VALUES ($1, $2)`,
		sqlite3:  `INSERT INTO project_containers (project_id, container_id) VALUES ($1, $2)`,
	}

	_, err := pct.db.execute(sql, projectID, containerID)

	if err != nil {
		if pct.db.IsUniqueViolationError(err) {
			return err
		}

		return err
	}

	return err
}

func (pct *ProjectContainersTable) getAllContainers(projectID int) ([]*Container, error) {
	sql := SQLCommand{
		postgres: `
			SELECT containers.id, containers.name FROM containers
			INNER JOIN project_containers ON containers.id = project_containers.container_id
			WHERE project_containers.project_id = $1
		`,
		sqlite3: `
			SELECT containers.id, containers.name FROM containers
			INNER JOIN project_containers ON containers.id = project_containers.container_id
			WHERE project_containers.project_id = $1
		`,
	}

	rows, err := pct.db.query(sql, projectID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var containers []*Container

	for rows.Next() {
		c := &Container{db: pct.db}
		err = rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, err
		}
		containers = append(containers, c)
	}

	return containers, nil
}
