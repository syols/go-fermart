package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/syols/go-devops/config"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/syols/go-devops/internal/pkg"
)

func userRegisterDatabase() (*pkg.Database, error) {
	mockDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, err
	}
	mock.ExpectQuery(InsertQuery).WillReturnRows(sqlmock.NewRows([]string{"login", "password"}))
	db, err := pkg.NewDatabase(pkg.NewSQLConnection(mockDb, "sqlmock"))
	if err != nil {
		return nil, err
	}
	db.Scripts["user_register.sql"] = InsertQuery
	return &db, err
}

func userLoginDatabase() (*pkg.Database, error) {
	mockDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, err
	}
	mock.ExpectQuery(SelectQuery).WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password"}).AddRow(0, "login", "password"))
	db, err := pkg.NewDatabase(pkg.NewSQLConnection(mockDb, "sqlmock"))
	if err != nil {
		return nil, err
	}
	db.Scripts["user_login.sql"] = SelectQuery
	return &db, err
}

func TestUserRegisterHandler(t *testing.T) {
	auth := pkg.NewAuthorizer(config.Config(config.Config{}))
	database, err := userRegisterDatabase()
	assert.NoError(t, err)

	router := gin.Default()
	router.POST("/register", Register(*database, auth))

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/register", strings.NewReader(`{"login": "login", "password": "password"}`))
	assert.NoError(t, err)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestUserLoginHandler(t *testing.T) {
	auth := pkg.NewAuthorizer(config.Config(config.Config{}))
	database, err := userLoginDatabase()
	assert.NoError(t, err)

	router := gin.Default()
	router.GET("/login", Login(*database, auth))

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/login", strings.NewReader(`{"login": "login", "password": "password"}`))
	assert.NoError(t, err)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
