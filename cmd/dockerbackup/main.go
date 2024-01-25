package main

import (
	"docker-backup/internal/cli"
	"docker-backup/internal/db"
	"docker-backup/internal/db/driver"
	"fmt"
)

func main() {
	driver, err := driver.NewPostgresDriver("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}

	_, err = db.NewDatabaseClient(driver)
	if err != nil {
		fmt.Printf("Could not establish connection to database: %s\n", err.Error())
		return
	}

	err = cli.BackupProjectCmd("test")
}
