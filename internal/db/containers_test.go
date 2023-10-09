package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitContainersTable(t *testing.T) {
	d := getTestingDriver()
	db, err := Connect(d)

	if err != nil {
		t.Error(err)
	}

	defer db.Close()

	err = db.InitContainersTable()
	require.NoError(t, err)

	err = db.Close()
	require.NoError(t, err)
}

func TestAddContainerToTable(t *testing.T) {
	db, err := getTestingDb()
	require.NoError(t, err)

	defer db.Close()

	c, err := db.ct.Add("asdfasdf", "test")
	require.NoError(t, err)

	assert.Len(t, db.ct.containers, 1)

	assert.Equal(t, c.Id, db.ct.containers[0].Id)
	assert.Equal(t, c.Name, db.ct.containers[0].Name)
}

func TestAddDuplicateContainerToTable(t *testing.T) {
	db, err := getTestingDb()
	require.NoError(t, err)

	defer db.Close()

	c, err := db.ct.Add("asdfasdf", "test")
	require.NoError(t, err)

	assert.Len(t, db.ct.containers, 1)

	assert.Equal(t, c.Id, db.ct.containers[0].Id)
	assert.Equal(t, c.Name, db.ct.containers[0].Name)

	_, err = db.ct.Add("asdfasdf", "test")
	_, ok := err.(ContainerAlreadyExistsError)
	assert.True(t, ok)

	assert.Len(t, db.ct.containers, 1)
	assert.Equal(t, c.Id, db.ct.containers[0].Id)
	assert.Equal(t, c.Name, db.ct.containers[0].Name)
}

func TestGetContainerInTableByName(t *testing.T) {
	db, err := getTestingDb()
	require.NoError(t, err)

	defer db.Close()

	c, err := db.ct.Add("asdfasdf", "test")
	require.NoError(t, err)

	c2, err := db.ct.GetByName("test")
	require.NoError(t, err)

	assert.Equal(t, c.Id, c2.Id)
	assert.Equal(t, c.Name, c2.Name)
}

func TestGetNonExistingContainerInTableByName(t *testing.T) {
	db, err := getTestingDb()
	require.NoError(t, err)

	defer db.Close()

	_, err = db.ct.GetByName("test")
	require.Error(t, err)
	_, ok := err.(ContainerNotFoundError)
	assert.True(t, ok)
}

func TestGetContainerInTableById(t *testing.T) {
	db, err := getTestingDb()
	require.NoError(t, err)

	defer db.Close()

	c, err := db.ct.Add("asdfasdf", "test")
	require.NoError(t, err)

	c2, err := db.ct.GetById(c.Id)
	require.NoError(t, err)

	assert.Equal(t, c.Id, c2.Id)
	assert.Equal(t, c.Name, c2.Name)
}

func TestGetNonExistingContainerInTableById(t *testing.T) {
	db, err := getTestingDb()
	require.NoError(t, err)

	defer db.Close()

	_, err = db.ct.GetById("asdfasdf")
	require.Error(t, err)
	_, ok := err.(ContainerNotFoundError)
	assert.True(t, ok)
}
