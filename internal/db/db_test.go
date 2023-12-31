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

func getTestingDb() (*Database, error) {
	d := getTestingDriver()
	db, err := Connect(d)

	if err != nil {
		return nil, err
	}

	err = db.InitTables()

	if err != nil {
		return nil, err
	}

	return db, nil
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

func TestAddProject(t *testing.T) {
	db, err := getTestingDb()
	require.NoError(t, err)

	err = db.AddProject("test")
	require.NoError(t, err)

	assert.Len(t, db.projects, 1)

	db.projects = []*Project{}

	err = db.AddProject("test")
	require.Error(t, err)

	_, ok := err.(ProjectAlreadyExistsError)
	assert.True(t, ok)

	assert.Len(t, db.projects, 1)
	assert.Equal(t, "test", db.projects[0].Name)

	err = db.Close()
}

func TestGetAllProjects(t *testing.T) {
	db, err := getTestingDb()
	require.NoError(t, err)

	err = db.AddProject("test")
	require.NoError(t, err)

	err = db.AddProject("test2")
	require.NoError(t, err)

	projects, err := db.GetAllProjects()
	require.NoError(t, err)

	assert.Len(t, projects, 2)

	assert.Equal(t, "test", projects[0].Name)
	assert.Equal(t, "test2", projects[1].Name)

	err = db.Close()
	require.NoError(t, err)
}
