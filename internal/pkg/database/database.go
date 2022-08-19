package database

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/syols/go-devops/config"
)

const ScriptPath = "internal/pkg/database/scripts/"
const MigrationPath = "file://internal/pkg/database/migrations"

type RelativePath string

const (
	UserLoginQuery        RelativePath = "user/login.sql"
	UserRegisterQuery     RelativePath = "user/register.sql"
	UserSelectQuery       RelativePath = "user/select.sql"
	UserOrdersSelectQuery RelativePath = "order/user.sql"
	OrderCreateQuery      RelativePath = "order/create.sql"
	OrderUpdateQuery      RelativePath = "order/update.sql"
	OrderSelectQuery      RelativePath = "order/select.sql"
)

type Connection struct {
	string
}

func NewConnection(config config.Config) (db Connection, err error) {
	db = Connection{
		config.DatabaseConnectionString,
	}

	m, err := migrate.New(MigrationPath, config.DatabaseConnectionString)
	if err != nil {
		return db, err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return db, err
	}
	return db, nil
}

func (c Connection) Execute(ctx context.Context, script RelativePath, model any) (*sqlx.Rows, error) {
	query, err := ioutil.ReadFile(filepath.Join(ScriptPath, string(script)))
	if err != nil {
		return nil, err
	}

	db, err := sqlx.ConnectContext(ctx, "postgres", c.string)
	if err != nil {
		return nil, err
	}

	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	return db.NamedQuery(string(query), model)
}

func ScanOne[T any](rows sqlx.Rows) (*T, error) {
	var value T
	if rows.Next() {
		if err := rows.StructScan(&value); err != nil {
			return nil, err
		}
		return &value, nil
	}
	return nil, nil
}

func ScanAll[T any](rows sqlx.Rows) (*[]T, error) {
	var values []T
	for rows.Next() {
		var value T
		if err := rows.StructScan(&value); err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return &values, nil
}
