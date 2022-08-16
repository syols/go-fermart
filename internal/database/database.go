package database

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Database struct {
	connectionString string
}

func NewDatabase(connectionString string) (Database, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return Database{}, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return Database{}, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///internal/database/migrations",
		"postgres", driver)
	if err := m.Up(); err != nil {
		return Database{}, err
	}

	return Database{
		connectionString: connectionString,
	}, nil
}
