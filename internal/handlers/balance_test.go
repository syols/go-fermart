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

func database() (*pkg.Database, error) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		return nil, err
	}
	query := "SELECT orders.user_id, " +
		"SUM(CASE action WHEN 'PURCHASE' THEN score ELSE 0 END) - " +
		"SUM(CASE action WHEN 'WITHDRAW' THEN score ELSE 0 END) " +
		"AS current, SUM(CASE action WHEN 'WITHDRAW' " +
		"THEN score ELSE 0 END) AS withdrawn FROM orders WHERE user_id = 1 AND status = 'PROCESSED' GROUP BY user_id;"
	mock.ExpectQuery(query).
		WithArgs("login", "password").WillReturnRows(sqlmock.NewRows([]string{"user_id", "current", "withdrawn"}).
		AddRow(0, 1, 2))
	conn := pkg.NewSqlConnection(mockDb, "sqlmock")
	db, err := pkg.NewDatabase(conn)
	db.Scripts["user_balance.sql"] = query
	return &db, err
}

func router(database pkg.Database) *gin.Engine {
	route := gin.Default()
	route.GET("/balance", Balance(database))
	return route
}

func TestBalanceHandler(t *testing.T) {
	database, err := database()
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
