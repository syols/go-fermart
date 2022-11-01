package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/syols/go-devops/internal/models"
	"github.com/syols/go-devops/internal/pkg"
)

const SelectQuery = "SELECT (*)*"

func balanceDatabase() (*pkg.Database, error) {
	mockDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, err
	}

	mock.ExpectQuery(SelectQuery).WillReturnRows(sqlmock.NewRows([]string{"user_id", "current", "withdrawn"}).AddRow(0, 1, 2))
	db, err := pkg.NewDatabase(pkg.NewSQLConnection(mockDb, "sqlmock"))
	if err != nil {
		return nil, err
	}
	db.Scripts["user_balance.sql"] = SelectQuery
	return &db, err
}

func router(database pkg.Database) *gin.Engine {
	route := gin.Default()
	route.GET("/balance", Balance(database))
	return route
}

func TestBalanceHandler(t *testing.T) {
	database, err := balanceDatabase()
	assert.NoError(t, err)

	router := router(*database)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/balance", nil)
	router.ServeHTTP(w, req)

	var first, second models.Balance
	assert.NoError(t, json.Unmarshal([]byte("{\"current\": 1, \"withdrawn\": 2}"), &first))
	assert.NoError(t, json.Unmarshal([]byte(w.Body.String()), &second))

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, first, second)
}
