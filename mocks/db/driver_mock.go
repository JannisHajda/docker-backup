package db_mocks

import (
	"database/sql"

	"github.com/stretchr/testify/mock"
)

type DatabaseDriverMock struct {
	mock.Mock
}

func NewDatabaseDriverMock() *DatabaseDriverMock {
	return &DatabaseDriverMock{}
}

func (m *DatabaseDriverMock) Exec(query string, args ...interface{}) (sql.Result, error) {
	mockArgs := m.Called(query, args)

	var result sql.Result
	if mockArgs.Get(0) != nil {
		result = mockArgs.Get(0).(sql.Result)
	}

	var err error
	if mockArgs.Get(1) != nil {
		err = mockArgs.Get(1).(error)
	}

	return result, err
}

func (m *DatabaseDriverMock) Query(query string, args ...interface{}) (*sql.Rows, error) {
	mockArgs := m.Called(query, args)

	var result *sql.Rows
	if mockArgs.Get(0) != nil {
		result = mockArgs.Get(0).(*sql.Rows)
	}

	var err error
	if mockArgs.Get(1) != nil {
		err = mockArgs.Get(1).(error)
	}

	return result, err
}
