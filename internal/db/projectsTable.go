package db

import "fmt"

type ProjectsTable struct {
	db *Database
}

func newProjectsTable(db *Database) (*ProjectsTable, error) {
	pt := &ProjectsTable{db: db}
	err := pt.init()
	if err != nil {
		return nil, err
	}

	return pt, nil
}

func (pt *ProjectsTable) init() error {
	sql := SQLCommand{
		postgres: `
			CREATE TABLE IF NOT EXISTS projects (
				id SERIAL PRIMARY KEY NOT NULL,
				name TEXT NOT NULL UNIQUE
			);
		`,
		sqlite3: `
			CREATE TABLE IF NOT EXISTS projects (
				id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
				name TEXT NOT NULL UNIQUE
			);
		`,
	}

	_, err := pt.db.execute(sql)
	return err
}

func (pt *ProjectsTable) add(name string) (*Project, error) {
	sql := SQLCommand{
		postgres: `INSERT INTO projects (name) VALUES ($1) RETURNING id`,
		sqlite3:  `INSERT INTO projects (name) VALUES ($1) RETURNING id`,
	}

	p := &Project{
		db:   pt.db,
		Name: name,
	}

	row, err := pt.db.queryRow(sql, name)
	if err != nil {
		return nil, err
	}

	err = row.Scan(&p.ID)
	if err != nil {
		if pt.db.IsUniqueViolationError(err) {
			return nil, err
		}

		return nil, err
	}

	fmt.Printf("%+v\n", p)

	return &Project{
		db:   pt.db,
		Name: name,
	}, nil
}

func (pt *ProjectsTable) getByName(name string) (*Project, error) {
	sql := SQLCommand{
		postgres: `SELECT id, name FROM projects WHERE name = $1`,
		sqlite3:  `SELECT id, name FROM projects WHERE name = $1`,
	}

	row, err := pt.db.queryRow(sql, name)
	if err != nil {
		return nil, err
	}

	p := &Project{
		db: pt.db,
	}

	err = row.Scan(&p.ID, &p.Name)
	if err != nil {
		if pt.db.IsNoRowsError(err) {
			return nil, err
		}

		return nil, err
	}

	return p, nil
}

func (pt *ProjectsTable) getAll() ([]*Project, error) {
	sql := SQLCommand{
		postgres: `SELECT id, name FROM projects`,
		sqlite3:  `SELECT id, name FROM projects`,
	}

	rows, err := pt.db.query(sql)
	if err != nil {
		return nil, err
	}

	projects := []*Project{}

	for rows.Next() {
		p := &Project{
			db: pt.db,
		}

		err := rows.Scan(&p.ID, &p.Name)
		if err != nil {
			return nil, err
		}

		projects = append(projects, p)
	}

	return projects, nil
}
