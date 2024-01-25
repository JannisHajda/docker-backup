package db_tests

import (
	"docker-backup/errors"
	"docker-backup/internal/db"
	db_mocks "docker-backup/mocks/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProjectByNameNonExisting(t *testing.T) {
	driver := db_mocks.NewDatabaseDriverMock()
	driver.On("Query", "SELECT id, name FROM projects WHERE id = $1;", []interface{}{int64(1)}).Return(nil, errors.NewItemNotFoundError(nil))

	client := &db.DatabaseClient{}
	pt := &db.ProjectsTable{}

	client.SetDriver(driver)
	client.SetProjectsTable(pt)
	pt.SetClient(client)

	_, err := client.GetProjectByID(1)
	assert.ErrorAs(t, err, &errors.ProjectNotFoundError{})
}
