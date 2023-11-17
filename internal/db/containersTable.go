package db

type ContainersTable struct {
	db         *Database
	containers []*Container
}

func newContainersTable(db *Database) (*ContainersTable, error) {
	ct := &ContainersTable{db: db}
	err := ct.init()
	if err != nil {
		return nil, err
	}

	return ct, nil
}

func (ct *ContainersTable) init() error {
	sql := SQLCommand{
		postgres: `
			CREATE TABLE IF NOT EXISTS containers (
				id TEXT PRIMARY KEY,
				name TEXT NOT NULL UNIQUE
			);
		`,
		sqlite3: `
			CREATE TABLE IF NOT EXISTS containers (
				id TEXT PRIMARY KEY,
				name TEXT NOT NULL UNIQUE
			);
		`,
	}

	_, err := ct.db.execute(sql)
	return err
}

func (ct *ContainersTable) add(id string, name string) (*Container, error) {
	sql := SQLCommand{
		postgres: `INSERT INTO containers (id, name) VALUES ($1, $2) RETURNING id`,
		sqlite3:  `INSERT INTO containers (id, name) VALUES ($1, $2) RETURNING id`,
	}

	c := &Container{
		db:   ct.db,
		ID:   id,
		Name: name,
	}

	row, err := ct.db.queryRow(sql, id, name)

	if err != nil {
		return nil, err
	}

	err = row.Scan(&c.ID)
	if err != nil {
		if ct.db.IsUniqueViolationError(err) {
			return nil, err
		}
	}

	return c, nil
}

func (ct *ContainersTable) getByID(id string) (*Container, error) {
	sql := SQLCommand{
		postgres: `SELECT id, name FROM containers WHERE id = $1`,
		sqlite3:  `SELECT id, name FROM containers WHERE id = $1`,
	}

	c := &Container{
		db: ct.db,
		ID: id,
	}

	row, err := ct.db.queryRow(sql, id)
	if err != nil {
		return nil, err
	}

	err = row.Scan(&c.ID, &c.Name)
	if err != nil {
		if ct.db.IsNoRowsError(err) {
			return nil, err
		}

		return nil, err
	}

	return c, nil
}

func (ct *ContainersTable) getAll() ([]*Container, error) {
	sql := SQLCommand{
		postgres: `SELECT id, name FROM containers`,
		sqlite3:  `SELECT id, name FROM containers`,
	}

	rows, err := ct.db.query(sql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	containers := []*Container{}

	for rows.Next() {
		c := &Container{
			db: ct.db,
		}

		err := rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, err
		}

		containers = append(containers, c)
	}

	return containers, nil
}
