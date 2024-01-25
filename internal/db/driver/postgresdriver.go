package driver

import (
	"database/sql"
	"docker-backup/errors"
	"docker-backup/interfaces"

	"github.com/lib/pq"
)

type PostgresDriver struct {
	conn string
}

func (d *PostgresDriver) Connect() (*sql.DB, error) {
	conn, err := sql.Open("postgres", d.conn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func NewPostgresDriver(connection string) (interfaces.DatabaseDriver, error) {
	driver := &PostgresDriver{conn: connection}
	_, err := driver.Connect()
	if err != nil {
		return nil, err
	}

	return driver, nil
}

func (d *PostgresDriver) handleError(err error) error {
	pqErr, ok := err.(*pq.Error)
	if !ok {
		return err
	}

	switch pqErr.Code {
	// item not found
	case "23503":
		return errors.NewItemNotFoundError(pqErr)
	case "23505":
		return errors.NewDuplicateKeyError(pqErr)
	default:
		return err
	}
}

func (d *PostgresDriver) Exec(query string, args ...interface{}) (sql.Result, error) {
	conn, err := d.Connect()
	defer conn.Close()

	if err != nil {
		return nil, err
	}

	res, err := conn.Exec(query, args...)
	if err != nil {
		return nil, d.handleError(err)
	}

	return res, nil
}

func (d *PostgresDriver) Query(query string, args ...interface{}) (*sql.Rows, error) {
	conn, err := d.Connect()
	defer conn.Close()

	if err != nil {
		return nil, err
	}

	res, err := conn.Query(query, args...)
	if err != nil {
		return nil, d.handleError(err)
	}

	return res, nil
}
