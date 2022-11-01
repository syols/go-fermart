package handlers

import (
	"bou.ke/monkey"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/syols/go-devops/config"
	"github.com/syols/go-devops/internal/pkg"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const Header = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk0NjY5MjA2MSwiaXNz" +
	"IjoidXNlcm5hbWUiLCJVc2VybmFtZSI6InVzZXJuYW1lIn0.7Atr6d4ZpCmkGdqrE6yiBfnVktp7wrMURy4WUmRvdXI"

func authDatabase() (*pkg.Database, error) {
	mockDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, err
	}

	mock.ExpectQuery(SelectQuery).WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password"}).
		AddRow(0, "username", "password"))

	db, err := pkg.NewDatabase(pkg.NewSqlConnection(mockDb, "sqlmock"))
	if err != nil {
		return nil, err
	}
	db.Scripts["user_select.sql"] = SelectQuery
	return &db, err
}

func TestAuthMiddleware(t *testing.T) {
	monkey.Patch(time.Now, func() time.Time {
		return time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)
	})

	cfg := config.Config{Sign: "some_sign"}
	auth := pkg.NewAuthorizer(cfg)
	database, err := authDatabase()
	assert.NoError(t, err)

	router := gin.Default()
	router.Use(AuthMiddleware(*database, auth))
	router.GET("/healthcheck", Healthcheck)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthcheck", nil)
	req.Header.Set("Authorization", Header)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
