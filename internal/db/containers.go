package db

type ContainersTable struct {
	db         *Database
	containers []*Container
}

type Container struct {
	Id   string
	Name string
}

type ContainerAlreadyExistsError struct {
	Name string
	Err  error
}

func (cae ContainerAlreadyExistsError) Error() string {
	return "Container with name " + cae.Name + " already exists"
}

type ContainerNotFoundError struct {
	Id   string
	Name string
	Err  error
}

func (cne ContainerNotFoundError) Error() string {
	return "Container with id " + cne.Id + " and name " + cne.Name + " not found"
}

func (db *Database) InitContainersTable() error {
	db.ct = &ContainersTable{db: db, containers: []*Container{}}

	_, err := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS containers (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		);
	`)

	return err
}

func (ct *ContainersTable) Add(id, name string) (*Container, error) {
	c := &Container{Id: id, Name: name}

	_, err := ct.db.conn.Exec("INSERT INTO containers (id, name) VALUES ($1, $2)", id, name)

	if err != nil {
		if ct.db.driver.GetName() == "sqlite3" && (err.Error() == "UNIQUE constraint failed: containers.id") || (err.Error() == "UNIQUE constraint failed: containers.name") {
			return nil, ContainerAlreadyExistsError{Name: name, Err: nil}
		}

		if ct.db.driver.GetName() == "postgres" && err.Error() == "pq: duplicate key value violates unique constraint \"containers_pkey\"" {
			return nil, ContainerAlreadyExistsError{Name: name, Err: nil}
		}

		return nil, err
	}

	ct.containers = append(ct.containers, c)

	return c, nil
}

func (ct *ContainersTable) GetById(id string) (*Container, error) {
	c := &Container{Id: id}
	err := ct.db.conn.QueryRow("SELECT name FROM containers WHERE id = $1", id).Scan(&c.Name)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, ContainerNotFoundError{Id: id, Name: "", Err: nil}
		}

		return nil, err
	}

	return c, nil
}

func (ct *ContainersTable) GetByName(name string) (*Container, error) {
	c := &Container{Name: name}
	err := ct.db.conn.QueryRow("SELECT id FROM containers WHERE name = $1", name).Scan(&c.Id)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, ContainerNotFoundError{Id: "", Name: name, Err: nil}
		}

		return nil, err
	}

	return c, nil
}

func (ct *ContainersTable) GetOrCreate(id, name string) (*Container, error) {
	c, err := ct.GetById(id)

	if err != nil {
		if _, ok := err.(ContainerNotFoundError); ok {
			return ct.Add(id, name)
		}

		return nil, err
	}

	return c, nil
}
