package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syols/go-devops/internal/models"
	"github.com/syols/go-devops/internal/pkg"
)

func Balance(connection pkg.Database) gin.HandlerFunc {
	return func(context *gin.Context) {
		userID := context.GetInt("id")

		purchase, err := models.CalculateBalance(context, connection, userID)
		if err != nil {
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		context.JSON(http.StatusOK, purchase)
	}
}
