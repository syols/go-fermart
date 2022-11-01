package pkg

import (
	"context"
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/syols/go-devops/config"
)

type DatabaseConnectionCreator interface {
	Migrate() error
	Create(ctx context.Context) (*sqlx.DB, error)
	Close(*sqlx.DB)
}

type UrlConnection struct {
	databaseURL string
}

type SqlConnection struct {
	db         *sql.DB
	driverName string
}

func NewDatabaseUrlConnection(config config.Config) DatabaseConnectionCreator {
	return UrlConnection{
		databaseURL: config.DatabaseURL,
	}
}

func NewSqlConnection(db *sql.DB, driverName string) DatabaseConnectionCreator {
	return SqlConnection{
		db:         db,
		driverName: driverName,
	}
}

func (c UrlConnection) Create(ctx context.Context) (*sqlx.DB, error) {
	return sqlx.ConnectContext(ctx, "postgres", c.databaseURL)
}

func (c UrlConnection) Migrate() error {
	m, err := migrate.New(MigrationPath, c.databaseURL)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

func (c UrlConnection) Close(db *sqlx.DB) {
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)
}

func (c SqlConnection) Create(_ context.Context) (*sqlx.DB, error) {
	dbx := sqlx.NewDb(c.db, c.driverName)
	return dbx, nil
}

func (c SqlConnection) Migrate() error {
	return nil
}

func (c SqlConnection) Close(_ *sqlx.DB) {
	return
}
