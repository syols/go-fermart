package handlers

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/syols/go-devops/internal/pkg"
)

type MockConnection struct {
	db *sql.DB
}


func (c MockConnection) create(_ context.Context) (*sqlx.DB, error) {
	dbx := sqlx.NewDb(c.db, "sqlmock")
	return dbx, nil
}


func handlers(t *testing.T) http.Handler {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	mock.ExpectExec("INSERT INTO users").
		WithArgs("john", AnyTime{}).
		WillReturnResult(sqlmock.NewResult(1, 1))

	r := http.NewServeMux()
	db = pkg.Database{}
	r.HandleFunc("/balance", Balance(db))
	return r
}

func BalanceHandlerTest() {

}
