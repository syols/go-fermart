package pkg

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"

	"github.com/jmoiron/sqlx"
	"github.com/syols/go-devops/config"
)

type DatabaseConnectionCreator interface {
	Migrate() error
	Create(ctx context.Context) (*sqlx.DB, error)
	Close(*sqlx.DB)
}

type URLConnection struct {
	databaseURL string
}

type SQLConnection struct {
	db         *sql.DB
	driverName string
}

func NewDatabaseURLConnection(config config.Config) DatabaseConnectionCreator {
	return URLConnection{
		databaseURL: config.DatabaseURL,
	}
}

func NewSQLConnection(db *sql.DB, driverName string) DatabaseConnectionCreator {
	return SQLConnection{
		db:         db,
		driverName: driverName,
	}
}

func (c URLConnection) Create(ctx context.Context) (*sqlx.DB, error) {
	return sqlx.ConnectContext(ctx, "postgres", c.databaseURL)
}

func (c URLConnection) Migrate() error {
	m, err := migrate.New(MigrationPath, c.databaseURL)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

func (c URLConnection) Close(db *sqlx.DB) {
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)
}

func (c SQLConnection) Create(_ context.Context) (*sqlx.DB, error) {
	dbx := sqlx.NewDb(c.db, c.driverName)
	return dbx, nil
}

func (c SQLConnection) Migrate() error {
	return nil
}

func (c SQLConnection) Close(_ *sqlx.DB) {}
