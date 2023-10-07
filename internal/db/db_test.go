package db

import (
	"testing"

	"github.com/JannisHajda/docker-backup/internal/db/drivers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestingDriver() drivers.Driver {
	return drivers.SqliteDriver{}
}

func TestInitTables(t *testing.T) {
	d := getTestingDriver()
	db, err := Connect(d)

	if err != nil {
		t.Error(err)
	}

	err = db.InitTables()
	require.NoError(t, err)

	assert.Empty(t, db.projects)

	err = db.Close()
	require.NoError(t, err)
}
