package main

import (
	"os"

	"github.com/JannisHajda/docker-backup/internal/db"
	"github.com/JannisHajda/docker-backup/internal/db/drivers"
	"github.com/JannisHajda/docker-backup/internal/utils"
)

func main() {
	err := utils.PrepareEnv()

	if err != nil {
		panic(err)
	}

	driver := drivers.PostgresDriver{
		User:     os.Getenv("PG_USER"),
		Password: os.Getenv("PG_PASSWORD"),
		Host:     os.Getenv("PG_HOST"),
		Port:     os.Getenv("PG_PORT"),
		Database: os.Getenv("PG_DATABASE"),
		Sslmode:  os.Getenv("PG_SSLMODE"),
	}

	db, err := db.Connect(driver)

	if err != nil {
		panic(err)
	}

	err = db.InitTables()

	if err != nil {
		panic(err)
	}

	err = db.AddProject("test2")

	if err != nil {
		panic(err)
	}

	err = db.AddProject("test")

	if err != nil {
		panic(err)
	}

	defer db.Close()

}
