package pkg

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

const ScriptPath = "scripts/query/"
const MigrationPath = "file://scripts/migrations/"

type RelativePath string

type DatabaseConnectionCreator interface {
	create(ctx context.Context) (*sqlx.DB, error)
}

type SqlConnection struct {
	databaseURL string
}

type Database struct {
	scripts    map[string]string
	connection DatabaseConnectionCreator
}

func NewDatabase(config config.Config) (db Database, err error) {
	var conn DatabaseConnectionCreator = SqlConnection{
		databaseURL: config.DatabaseURL,
	}
	db = Database{
		scripts:     map[string]string{},
		connection: conn,
	}

	m, err := migrate.New(MigrationPath, config.DatabaseURL)
	if err != nil {
		return db, err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return db, err
	}
	return
}

func (d *Database) Execute(ctx context.Context, filename string, model interface{}) (*sqlx.Rows, error) {
	script, err := d.script(filename)
	if err != nil {
		return nil, err
	}

	db, err := d.connection.create(ctx)
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

func (c SqlConnection) create(ctx context.Context) (*sqlx.DB, error) {
	return sqlx.ConnectContext(ctx, "postgres", c.databaseURL)
}

func (d *Database) script(filename string) (string, error) {
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
