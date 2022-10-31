package pkg

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/syols/go-devops/config"
)

type DatabaseConnectionCreator interface {
	Url() string
	Create(ctx context.Context) (*sqlx.DB, error)
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

func (c UrlConnection) Url() string {
	return c.databaseURL
}

func (c SqlConnection) Create(_ context.Context) (*sqlx.DB, error) {
	dbx := sqlx.NewDb(c.db, c.driverName)
	return dbx, nil
}

func (c SqlConnection) Url() string {
	return "mock"
}
