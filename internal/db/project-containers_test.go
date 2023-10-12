package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitProjectContainersTable(t *testing.T) {
	d := getTestingDriver()
	db, err := Connect(d)

	if err != nil {
		t.Error(err)
	}

	defer db.Close()

	err = db.InitProjectContainersTable()
	require.NoError(t, err)

	assert.Empty(t, db.pct.projectContainers)

	err = db.Close()
	require.NoError(t, err)
}

func TestAddContainerToProjectInTable(t *testing.T) {
	db, err := getTestingDb()
	require.NoError(t, err)

	defer db.Close()

	p, err := db.pt.Add("test")
	require.NoError(t, err)

	c, err := db.ct.Add("asdfasdf", "test")
	require.NoError(t, err)

	pc, err := db.pct.Add(p.Id, c.Id)
	require.NoError(t, err)

	assert.Len(t, db.pct.projectContainers, 1)

	assert.Equal(t, pc.ProjectId, db.pct.projectContainers[0].ProjectId)
	assert.Equal(t, pc.ContainerId, db.pct.projectContainers[0].ContainerId)
}

func TestAddContainerTwiceToProjectInTable(t *testing.T) {
	db, err := getTestingDb()
	require.NoError(t, err)

	defer db.Close()

	p, err := db.pt.Add("test")
	require.NoError(t, err)

	c, err := db.ct.Add("asdfasdf", "test")
	require.NoError(t, err)

	pc, err := db.pct.Add(p.Id, c.Id)
	require.NoError(t, err)

	assert.Len(t, db.pct.projectContainers, 1)

	assert.Equal(t, pc.ProjectId, db.pct.projectContainers[0].ProjectId)

	_, err = db.pct.Add(p.Id, c.Id)
	_, ok := err.(ProjectContainerAlreadyExistsError)
	assert.True(t, ok)

	assert.Len(t, db.pct.projectContainers, 1)
	assert.Equal(t, pc.ProjectId, db.pct.projectContainers[0].ProjectId)
	assert.Equal(t, pc.ContainerId, db.pct.projectContainers[0].ContainerId)
}

func TestAddNonExistingContainerToProjectInTable(t *testing.T) {
	db, err := getTestingDb()
	require.NoError(t, err)

	defer db.Close()

	p, err := db.pt.Add("test")
	require.NoError(t, err)

	_, err = db.pct.Add(p.Id, "asdfasdf")
	_, ok := err.(ContainerNotFoundError)
	assert.True(t, ok)

	assert.Len(t, db.pct.projectContainers, 0)
}

func TestAddContainerToNonExistingProjectInTable(t *testing.T) {
	db, err := getTestingDb()
	require.NoError(t, err)

	defer db.Close()

	c, err := db.ct.Add("asdfasdf", "test")
	require.NoError(t, err)

	_, err = db.pct.Add(1, c.Id)
	_, ok := err.(ProjectNotFoundError)
	assert.True(t, ok)

	assert.Len(t, db.pct.projectContainers, 0)
}
