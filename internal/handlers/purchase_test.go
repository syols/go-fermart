package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/syols/go-devops/internal/pkg"
)

func purchasesDatabase() (*pkg.Database, error) {
	mockDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, err
	}

	mock.ExpectQuery(SelectQuery).WillReturnRows(sqlmock.NewRows([]string{"user_id", "number", "score", "status", "ctime"}).AddRow(0, 1, 2, "NEW", time.Now()))
	mock.ExpectQuery(InsertQuery).WillReturnRows(sqlmock.NewRows([]string{"user_id", "number", "score", "status", "action"}))
	db, err := pkg.NewDatabase(pkg.NewSQLConnection(mockDb, "sqlmock"))
	if err != nil {
		return nil, err
	}
	db.Scripts["user_orders.sql"] = SelectQuery
	db.Scripts["order_create.sql"] = InsertQuery
	return &db, err
}

func TestPurchasesHandler(t *testing.T) {
	database, err := purchasesDatabase()
	assert.NoError(t, err)

	router := gin.Default()
	router.Use(UserMiddleware())
	router.GET("/orders", Purchases(*database))
	router.POST("/order", CreatePurchase(*database, nil))

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/orders", nil)
	assert.NoError(t, err)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
