package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syols/go-devops/internal/models"
	"github.com/syols/go-devops/internal/pkg"
)

func CreateWithdraw(db pkg.Database) gin.HandlerFunc {
	return func(context *gin.Context) {
		userID, isOk := context.Get("id")
		if !isOk {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		withdraw, err := bindWithdraw(context)
		if err != nil {
			context.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}
		withdraw.UserID = userID.(int)

		balance, err := models.CalculateBalance(context, db, withdraw.UserID)
		if err != nil {
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if balance.Current < withdraw.Score {
			context.AbortWithStatus(http.StatusPaymentRequired)
			return
		}

		if err := withdraw.Create(context, db); err != nil {
			context.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}
		context.Status(http.StatusOK)
	}
}

func Withdrawals(db pkg.Database) gin.HandlerFunc {
	return func(context *gin.Context) {
		userID, isOk := context.Get("id")
		if !isOk {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		withdraws, err := models.LoadWithdraw(context, db, userID.(int))
		if err != nil {
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if withdraws == nil || len(*withdraws) == 0 {
			context.AbortWithStatus(http.StatusNoContent)
			return
		}
		context.JSON(http.StatusOK, withdraws)
	}
}

func bindWithdraw(context *gin.Context) (*models.Withdraw, error) {
	var withdraw models.Withdraw
	if err := context.BindJSON(&withdraw); err != nil {
		return nil, err
	}

	withdraw.Status = models.ProcessedOrderStatus
	withdraw.Action = models.WithdrawOrderAction
	if err := pkg.Validate(withdraw); err != nil {
		return nil, err
	}
	return &withdraw, nil
}
