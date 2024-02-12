package helper

import (
	"docker-backup/interfaces"
	"docker-backup/internal/db"
	"docker-backup/internal/db/driver"
	"docker-backup/internal/dockerclient"
)

func MapToSlice(m map[string]string) []string {
	s := make([]string, 0, len(m))
	for _, v := range m {
		s = append(s, v)
	}
	return s
}

func GetDBClient() (interfaces.DatabaseClient, error) {
	driver, err := driver.NewPostgresDriver("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}

	client, err := db.NewDatabaseClient(driver)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func GetDockerClient() (interfaces.DockerClient, error) {
	client, err := dockerclient.NewDockerClient()
	if err != nil {
		return nil, err
	}

	return client, nil
}
