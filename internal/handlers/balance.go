package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syols/go-devops/internal/models"
	"github.com/syols/go-devops/internal/pkg/storage"
)

func Balance(connection storage.Database) gin.HandlerFunc {
	return func(context *gin.Context) {
		UserID := context.GetInt("id")

		purchase, err := models.CalculateBalance(context, connection, UserID)
		if err != nil {
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		context.JSON(http.StatusOK, purchase)
	}
}
