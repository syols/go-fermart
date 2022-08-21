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

type Database struct {
	databaseUrl string
	scripts     map[string]string
}

func NewConnection(config config.Config) (connection Database, err error) {
	connection = Database{
		databaseUrl: config.DatabaseURL,
		scripts:     map[string]string{},
	}

	m, err := migrate.New(MigrationPath, connection.databaseUrl)
	if err != nil {
		return connection, err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return connection, err
	}
	return
}

func (d Database) Execute(ctx context.Context, filename string, model any) (*sqlx.Rows, error) {
	script, err := d.script(filename)
	if err != nil {
		return nil, err
	}

	db, err := sqlx.ConnectContext(ctx, "postgres", d.databaseUrl)
	if err != nil {
		return nil, err
	}

	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	return db.NamedQuery(script, model)
}

func (d Database) script(filename string) (string, error) {
	script, isOk := d.scripts[filename]
	if !isOk {
		bytes, err := ioutil.ReadFile(filepath.Join(ScriptPath, filename))
		if err != nil {
			return "", err
		}

		script = string(bytes)
		d.scripts[filename] = script
	}
	return script, nil
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
