package drivers

import "fmt"

type SqliteDriver struct {
}

func (sq SqliteDriver) GetName() string {
	return "sqlite3"
}

func (sq SqliteDriver) GetConnectionString() string {
	return fmt.Sprintf("file::memory:?cache=shared&_foreign_keys=on")
}

func (sq SqliteDriver) NoRowsError() string {
	return "no rows in result set"
}

func (sq SqliteDriver) UniqueViolationError() string {
	return "UNIQUE constraint failed"
}
