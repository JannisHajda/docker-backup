package drivers

import (
	"testing"
)

func newTestPostgresDriver() PostgresDriver {
	return PostgresDriver{
		User:     "test",
		Password: "randomPw",
		Host:     "localhost",
		Port:     "5432",
		Database: "TestDb",
		Sslmode:  "disable",
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
	expected := "postgres://test:randomPw@localhost:5432/TestDb?sslmode=disable"
	if d.GetConnectionString() != expected {
		t.Error("Expected ", expected, ", got ", d.GetConnectionString())
	}
}
