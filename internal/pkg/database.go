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

type Database struct {
	databaseURL string
	scripts     map[string]string
}

func NewDatabaseConnection(config config.Config) (connection Database, err error) {
	connection = Database{
		databaseURL: config.DatabaseURL,
		scripts:     map[string]string{},
	}

	m, err := migrate.New(MigrationPath, connection.databaseURL)
	if err != nil {
		return connection, err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return connection, err
	}
	return
}

func (d *Database) Execute(ctx context.Context, filename string, model interface{}) (*sqlx.Rows, error) {
	script, err := d.script(filename)
	if err != nil {
		return nil, err
	}

	db, err := sqlx.ConnectContext(ctx, "postgres", d.databaseURL)
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
