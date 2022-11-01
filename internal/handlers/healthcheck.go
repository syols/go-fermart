package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Healthcheck godoc
// @Summary Healthcheck
// @ID userID
// @Success 200 {object} OK
// @Router /api/healthcheck [get]
func Healthcheck(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
