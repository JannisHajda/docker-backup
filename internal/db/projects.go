package db

type ProjectsTable struct {
	db       *Database
	projects []*Project
}

type Project struct {
	Id   int64
	Name string
}

type ProjectAlreadyExistsError struct {
	Name string
	Err  error
}

func (pae ProjectAlreadyExistsError) Error() string {
	return "Project with name " + pae.Name + " already exists"
}

func (db *Database) InitProjectsTable() error {
	db.projects = []*Project{}
	db.pt = &ProjectsTable{db: db, projects: db.projects}

	if db.driver.GetName() == "sqlite3" {
		_, err := db.conn.Exec(`
			CREATE TABLE IF NOT EXISTS projects (
				id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
				name TEXT NOT NULL UNIQUE
			);
		`)

		return err
	}

	_, err := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS projects (
			id SERIAL PRIMARY KEY NOT NULL,
			name TEXT NOT NULL UNIQUE
		);
	`)

	return err
}

func (pt *ProjectsTable) Add(name string) (*Project, error) {
	p := &Project{Name: name}

	err := pt.db.conn.QueryRow(`
		INSERT INTO projects (name)
		VALUES ($1)
		RETURNING id;
	`, name).Scan(&p.Id)

	if err != nil {
		if pt.db.driver.GetName() == "sqlite3" && err.Error() == "UNIQUE constraint failed: projects.name" {
			return nil, ProjectAlreadyExistsError{Name: name, Err: err}
		}

		if pt.db.driver.GetName() == "postgres" && err.Error() == "pq: duplicate key value violates unique constraint \"projects_name_key\"" {
			return nil, ProjectAlreadyExistsError{Name: name, Err: err}
		}

		return nil, err

	}

	pt.projects = append(pt.projects, p)

	return p, nil
}

func (pt *ProjectsTable) GetById(id int64) (*Project, error) {
	p := &Project{}

	err := pt.db.conn.QueryRow(`
		SELECT id, name
		FROM projects
		WHERE id = $1;
	`, id).Scan(&p.Id, &p.Name)

	if err != nil {
		return nil, err
	}

	return p, nil
}

func (pt *ProjectsTable) GetByName(name string) (*Project, error) {
	p := &Project{}

	err := pt.db.conn.QueryRow(`
		SELECT id, name
		FROM projects
		WHERE name = $1;
	`, name).Scan(&p.Id, &p.Name)

	if err != nil {
		return nil, err
	}

	return p, nil
}
