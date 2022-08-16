package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Balance(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func Withdraw(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func Withdrawals(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
