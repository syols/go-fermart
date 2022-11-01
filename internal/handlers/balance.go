package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syols/go-devops/internal/models"
	"github.com/syols/go-devops/internal/pkg"
)

// Balance godoc
// @Tags Info
// @Summary Запрос баланса
// @ID userID
// @Accept  json
// @Produce json
// @Success 200 {object} Balance
// @Failure 500 {string} string "StatusInternalServerError"
// @Router /api/user/balance [get]
func Balance(db pkg.Database) gin.HandlerFunc {
	return func(context *gin.Context) {
		userID := context.GetInt("id")

		purchase, err := models.CalculateBalance(context, db, userID)
		if err != nil {
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		context.JSON(http.StatusOK, purchase)
	}
}
