package db_tests

import (
	"docker-backup/internal/db"
	db_mocks "docker-backup/mocks/db"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetContainer(t *testing.T) {
	driver := db_mocks.NewDatabaseDriverMock()
	client, err := db.NewDatabaseClient(driver)
	assert.NoError(t, err)

	fmt.Print(client.GetProjects())

	// Given
	// When
	// Then
	assert.Fail(t, "Not implemented")
}
