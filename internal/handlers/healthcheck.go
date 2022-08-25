package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Healthcheck(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
