package pkg

import (
	"context"
	"io/ioutil"
	"path/filepath"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const ScriptPath = "Scripts/query/"
const MigrationPath = "file://Scripts/migrations/"

type RelativePath string

type Database struct {
	Scripts    map[string]string
	connection DatabaseConnectionCreator
}

func NewDatabase(connectionCreator DatabaseConnectionCreator) (db Database, err error) {
	db = Database{
		Scripts:    map[string]string{},
		connection: connectionCreator,
	}
	err = connectionCreator.Migrate()
	return
}

func (d *Database) Execute(ctx context.Context, filename string, model interface{}) (*sqlx.Rows, error) {
	script, err := d.script(filename)
	if err != nil {
		return nil, err
	}

	db, err := d.connection.Create(ctx)
	if err != nil {
		return nil, err
	}
	defer d.connection.Close(db)
	return db.NamedQuery(script, model)
}

func (d *Database) script(filename string) (string, error) {
	script, isOk := d.Scripts[filename]
	if !isOk {
		bytes, err := ioutil.ReadFile(filepath.Join(ScriptPath, filename))
		if err != nil {
			return "", err
		}

		script = string(bytes)
		d.Scripts[filename] = script
	}
	return script, nil
}
