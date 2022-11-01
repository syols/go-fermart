package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/syols/go-devops/internal/pkg"
)

const InsertQuery = "INSERT *"

func withdrawalsDatabase() (*pkg.Database, error) {
	mockDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, err
	}

	mock.ExpectQuery(SelectQuery).WillReturnRows(sqlmock.NewRows([]string{"user_id", "current", "withdrawn"}).
		AddRow(0, 1, 2))

	mock.ExpectQuery(InsertQuery).WillReturnRows(sqlmock.NewRows([]string{"user_id", "number", "score", "status", "action"}))

	db, err := pkg.NewDatabase(pkg.NewSqlConnection(mockDb, "sqlmock"))
	if err != nil {
		return nil, err
	}
	db.Scripts["order_create.sql"] = InsertQuery
	db.Scripts["user_balance.sql"] = SelectQuery
	return &db, err
}

func UserMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set("username", "login")
		context.Set("id", 0)
		context.Next()
	}
}

func TestWithdrawHandler(t *testing.T) {
	database, err := withdrawalsDatabase()
	assert.NoError(t, err)

	router := gin.Default()
	router.Use(UserMiddleware())
	router.POST("/balance/withdraw", CreateWithdraw(*database))
	router.GET("/withdraw", Withdrawals(*database))

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/balance/withdraw", strings.NewReader(`{"order": "466417", "sum": 1}`))
	assert.NoError(t, err)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	req, err = http.NewRequest("GET", "/withdraw", nil)
	assert.NoError(t, err)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
