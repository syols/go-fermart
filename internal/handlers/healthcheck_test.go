package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthcheckHandler(t *testing.T) {
	router := gin.Default()
	router.Use(UserMiddleware())
	router.GET("/healthcheck", Healthcheck)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthcheck", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
