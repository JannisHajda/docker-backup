package drivers

import "fmt"

type PostgresDriver struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
	Sslmode  string
}

type PostgresOptions struct {
}

func (pg PostgresDriver) GetName() string {
	return "postgres"
}

func (pg PostgresDriver) GetConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		pg.User,
		pg.Password,
		pg.Host,
		pg.Port,
		pg.Database,
		pg.Sslmode,
	)
}

func (pg PostgresDriver) NoRowsError() string {
	return "no rows in result set"
}

func (pg PostgresDriver) UniqueViolationError() string {
	return "duplicate key value violates unique constraint"
}
