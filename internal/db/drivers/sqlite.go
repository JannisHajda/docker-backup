package drivers

import "fmt"

type SqliteDriver struct {
}

func (sq SqliteDriver) GetName() string {
	return "sqlite3"
}

func (sq SqliteDriver) GetConnectionString() string {
	return fmt.Sprintf("file::memory:?cache=shared")
}
