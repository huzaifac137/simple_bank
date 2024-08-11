package testdb

import (
	"database/sql"
	"errors"
)

type DBConfig struct {
	source string
	driver string
}

func ConnectTestDb() (db *sql.DB, err error) {

	dbConf := &DBConfig{
		source: "postgresql://postgres:qwerty7867@localhost:5432/simple_bank?sslmode=disable",
		driver: "postgres",
	}

	conn, err := sql.Open(dbConf.driver, dbConf.source)
	if err != nil {
		return nil, errors.New("cannot connect to test database")
	}

	return conn, nil
}
