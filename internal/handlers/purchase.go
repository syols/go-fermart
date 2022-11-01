package handlers

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syols/go-devops/internal/models"
	"github.com/syols/go-devops/internal/pkg"
)

// CreatePurchase godoc
// @Summary Создание заказа
// @ID userID
// @Status: NewOrderStatus,
// @Action: PurchaseOrderAction,
// @Success 200 {string} string "OK"
// @Success 201 {string} string "StatusInternalServerError"
// @Success 409 {string} string "StatusUnprocessableEntity"
// @Failure 500 {string} string "StatusInternalServerError"
// @Router /orders [post]
func CreatePurchase(db pkg.Database, sess *pkg.Session) gin.HandlerFunc {
	return func(context *gin.Context) {
		bytes, err := ioutil.ReadAll(context.Request.Body)
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		userID, isOk := context.Get("id")
		if !isOk {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		purchase := models.NewPurchase(string(bytes), userID.(int))
		if err := pkg.Validate(purchase); err != nil {
			context.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		dbPurchase, err := models.LoadPurchase(context, db, purchase.Number)
		if err != nil {
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if dbPurchase != nil {
			if purchase.UserID == dbPurchase.UserID {
				context.Status(http.StatusOK)
				return
			}
			context.AbortWithStatus(http.StatusConflict)
			return
		}

		if err := purchase.Create(context, db); err != nil {
			context.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		if sess != nil {
			if _, err = sess.SendMessage(purchase.Number); err != nil {
				context.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}
		context.Status(http.StatusAccepted)
	}
}

// Purchases godoc
// @Summary Список заказов
// @ID userID
// @Status: NewOrderStatus,
// @Action: PurchaseOrderAction,
// @Success 200 {objects} Purchase
// @Success 204 {string} string "StatusNoContent"
// @Success 400 {string} string "StatusBadRequest"
// @Success 409 {string} string "StatusUnprocessableEntity"
// @Failure 500 {string} string "StatusInternalServerError"
// @Router /orders [get]
func Purchases(db pkg.Database) gin.HandlerFunc {
	return func(context *gin.Context) {
		userID, isOk := context.Get("id")
		if !isOk {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		purchases, err := models.LoadPurchases(context, db, userID.(int))
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
