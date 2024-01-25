package main

import (
	"docker-backup/internal/db"
	"docker-backup/internal/db/driver"
	"docker-backup/internal/worker"
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

	worker, err := worker.NewWorker()
	if err != nil {
		fmt.Printf("Could not create worker: %s\n", err.Error())
		return
	}

	err = worker.BackupContainer("test-service")
	if err != nil {
		fmt.Printf("Could not backup container: %s\n", err.Error())
		return
	}
}
