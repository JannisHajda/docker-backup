package drivers

import "fmt"

type PostgresDriver struct {
	User     string
	Password string
	Host     string
	Port     string
	Sslmode  string
}

type PostgresOptions struct {
}

func (pg PostgresDriver) GetName() string {
	return "postgres"
}

func (pg PostgresDriver) GetConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=%s",
		pg.User,
		pg.Password,
		pg.Host,
		pg.Port,
		pg.Sslmode,
	)
}
