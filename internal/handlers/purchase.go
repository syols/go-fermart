package handlers

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syols/go-devops/internal/models"
	"github.com/syols/go-devops/internal/pkg/database"
	"github.com/syols/go-devops/internal/pkg/validator"
)

func CreatePurchase(connection database.Database) gin.HandlerFunc {
	return func(context *gin.Context) {
		bytes, err := ioutil.ReadAll(context.Request.Body)
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		userId, isOk := context.Get("id")
		if !isOk {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		purchase := models.NewPurchase(string(bytes), userId.(int))
		if err := validator.Validate(purchase); err != nil {
			context.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		dbPurchase, err := models.LoadPurchase(context, connection, purchase.UserId)
		if err != nil {
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if dbPurchase != nil {
			if purchase.UserId == dbPurchase.UserId {
				context.Status(http.StatusOK)
				return
			}
			context.AbortWithStatus(http.StatusConflict)
			return
		}

		if err := purchase.Create(context, connection); err != nil {
			context.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}
		context.Status(http.StatusAccepted)
	}
}

func Purchases(connection database.Database) gin.HandlerFunc {
	return func(context *gin.Context) {
		userId, isOk := context.Get("id")
		if !isOk {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		purchases, err := models.LoadPurchases(context, connection, userId.(int))
		if err != nil {
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if purchases == nil || len(*purchases) == 0 {
			context.AbortWithStatus(http.StatusNoContent)
			return
		}
		context.JSON(http.StatusOK, purchases)
	}
}
