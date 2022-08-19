package handlers

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syols/go-devops/internal/models"
	"github.com/syols/go-devops/internal/pkg/database"
)

func SetUserOrder(connection database.Connection) gin.HandlerFunc {
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

		order := models.NewOrder(string(bytes), userId.(int))
		if err := order.Validate(); err != nil {
			context.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		dbOrder, err := order.Select(context, connection)
		if err != nil {
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if dbOrder != nil {
			if order.UserId == dbOrder.UserId {
				context.Status(http.StatusOK)
				return
			}
			context.AbortWithStatus(http.StatusConflict)
			return
		}

		if err := order.Create(context, connection); err != nil {
			context.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}
		context.Status(http.StatusAccepted)
	}
}

func Orders(connection database.Connection) gin.HandlerFunc {
	return func(context *gin.Context) {
		userId, isOk := context.Get("id")
		if !isOk {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}
		order := models.Order{UserId: userId.(int)}
		orders, err := order.UserOrders(context, connection)
		if err != nil {
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		context.JSON(http.StatusOK, orders)
	}
}
