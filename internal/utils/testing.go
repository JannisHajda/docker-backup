package utils

import (
	"github.com/JannisHajda/docker-backup/internal/db"
	"github.com/JannisHajda/docker-backup/internal/db/drivers"
)

func GetTestingDriver() drivers.Driver {
	return drivers.SqliteDriver{}
}

func GetTestingDatabse() (*db.Database, error) {
	d := GetTestingDriver()
	db, err := db.Connect(d)

	if err != nil {
		return nil, err
	}

	err = db.InitTables()

	if err != nil {
		return nil, err
	}

	return db, nil
}
