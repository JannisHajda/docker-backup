package drivers

import "testing"

func newTestPostgresDriver() PostgresDriver {
	return PostgresDriver{
		User:    "test",
		Host:    "localhost",
		Port:    "5432",
		Sslmode: "disable",
	}
}

func TestPostgresGetName(t *testing.T) {
	d := newTestPostgresDriver()
	if d.GetName() != "postgres" {
		t.Error("Expected postgres, got ", d.GetName())
	}
}

func TestPostgresGetConnectionString(t *testing.T) {
	d := newTestPostgresDriver()
	expected := "postgres://test:@localhost:5432/postgres?sslmode=disable"
	if d.GetConnectionString() != expected {
		t.Error("Expected ", expected, ", got ", d.GetConnectionString())
	}
}
