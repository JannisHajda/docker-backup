package db

import "fmt"

type ProjectContainersTable struct {
	db                *Database
	projectContainers []*ProjectContainers
}

type ProjectContainers struct {
	ProjectId   int64
	ContainerId string
}

type ProjectContainerAlreadyExistsError struct {
	ProjectId   int64
	ContainerId string
	Err         error
}

func (pcae ProjectContainerAlreadyExistsError) Error() string {
	return "ProjectContainer with project id " + fmt.Sprint(pcae.ProjectId) + " and container id " + pcae.ContainerId + " already exists"
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

func (pct *ProjectContainersTable) handleForeignKeyConstraint(projectId int64, containerId string) error {
	if _, err := pct.db.pt.GetById(projectId); err != nil {
		return ProjectNotFoundError{Id: projectId, Name: "", Err: err}
	}

	if _, err := pct.db.ct.GetById(containerId); err != nil {
		return ContainerNotFoundError{Id: "", Name: containerId, Err: err}
	}

	return nil
}

func (pct *ProjectContainersTable) Add(projectId int64, containerId string) (*ProjectContainers, error) {
	pc := &ProjectContainers{ProjectId: projectId, ContainerId: containerId}

	_, err := pct.db.conn.Exec("INSERT INTO project_containers (project_id, container_id) VALUES ($1, $2)", projectId, containerId)

	if err != nil {
		if pct.db.driver.GetName() == "sqlite3" {
			if err.Error() == "UNIQUE constraint failed: project_containers.project_id, project_containers.container_id" {

				return nil, ProjectContainerAlreadyExistsError{ProjectId: projectId, ContainerId: containerId, Err: nil}
			}

			if err.Error() == "FOREIGN KEY constraint failed" {
				return nil, pct.handleForeignKeyConstraint(projectId, containerId)
			}
		}

		if pct.db.driver.GetName() == "postgres" {
			if err.Error() == "pq: duplicate key value violates unique constraint \"project_containers_pkey\"" {
				return nil, ProjectContainerAlreadyExistsError{ProjectId: projectId, ContainerId: containerId, Err: nil}

			}

			if err.Error() == "pq: insert or update on table \"project_containers\" violates foreign key constraint \"project_containers_project_id_fkey\"" {
				return nil, pct.handleForeignKeyConstraint(projectId, containerId)
			}
		}

		return nil, err
	}

	pct.projectContainers = append(pct.projectContainers, pc)

	return pc, nil
}
