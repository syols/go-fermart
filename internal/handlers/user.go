package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func Login(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func SetUserOrders(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func UserOrders(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
