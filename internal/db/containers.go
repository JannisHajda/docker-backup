package db

type ContainersTable struct {
	db         *Database
	containers []*Container
}

type Container struct {
	Id   string
	Name string
}

func (db *Database) InitContainersTable() error {
	db.ct = &ContainersTable{db: db, containers: []*Container{}}

	_, err := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS containers (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL
		);
	`)

	return err
}
