package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitProjectsTable(t *testing.T) {
	d := getTestingDriver()
	db, err := Connect(d)

	if err != nil {
		t.Error(err)
	}

	defer db.Close()

	err = db.InitProjectsTable()
	require.NoError(t, err)

	assert.Empty(t, db.projects)

	err = db.Close()
	require.NoError(t, err)
}

func TestAddProjectToTable(t *testing.T) {
	db, err := getTestingDb()
	require.NoError(t, err)

	defer db.Close()

	p, err := db.pt.Add("test")
	require.NoError(t, err)

	assert.Len(t, db.pt.projects, 1)

	assert.Equal(t, p.Id, db.pt.projects[0].Id)
	assert.Equal(t, p.Name, db.pt.projects[0].Name)
}

func TestGetProjectInTableByName(t *testing.T) {
	db, err := getTestingDb()
	require.NoError(t, err)

	defer db.Close()

	p, err := db.pt.Add("test")
	require.NoError(t, err)

	p2, err := db.pt.GetByName("test")
	require.NoError(t, err)

	assert.Equal(t, p.Id, p2.Id)
}
