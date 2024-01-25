package helper

import (
	"database/sql"
	"docker-backup/internal/db/driver"
)

func MapToSlice(m map[string]string) []string {
	s := make([]string, 0, len(m))
	for _, v := range m {
		s = append(s, v)
	}
	return s
}

func GetDBConnection() (*sql.DB, error) {
	driver, err := driver.NewPostgresDriver("postgres://postgres:postgres@localhost:5432/docker_backup?sslmode=disable")
	if err != nil {
		return nil, err
	}

	conn, err := driver.Connect()
	if err != nil {
		return nil, err
	}

	return conn, nil
}
