package drivers

import "testing"

func newTestSqliteDriver() SqliteDriver {
	return SqliteDriver{}
}

func TestSqliteGetName(t *testing.T) {
	d := newTestSqliteDriver()
	if d.GetName() != "sqlite3" {
		t.Error("Expected sqlite3, got ", d.GetName())
	}
}

func TestSqliteGetConnectionString(t *testing.T) {
	d := newTestSqliteDriver()
	expected := "file::memory:?cache=shared&_foreign_keys=on"
	if d.GetConnectionString() != expected {
		t.Error("Expected ", expected, ", got ", d.GetConnectionString())
	}
}
