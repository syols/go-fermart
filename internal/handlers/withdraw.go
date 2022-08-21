package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syols/go-devops/internal/models"
	"github.com/syols/go-devops/internal/pkg/database"
	"github.com/syols/go-devops/internal/pkg/validator"
)

func CreateWithdraw(connection database.Database) gin.HandlerFunc {
	return func(context *gin.Context) {
		userId, isOk := context.Get("id")
		if !isOk {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		withdraw, err := bindWithdraw(context)
		if err != nil {
			context.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}
		withdraw.UserId = userId.(int)

		balance, err := models.CalculateBalance(context, connection, withdraw.UserId)
		if err != nil {
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if balance.Current < withdraw.Score {
			context.AbortWithStatus(http.StatusPaymentRequired)
			return
		}

		if err := withdraw.Create(context, connection); err != nil {
			context.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}
		context.Status(http.StatusAccepted)
	}
}

func Withdrawals(connection database.Database) gin.HandlerFunc {
	return func(context *gin.Context) {
		userId, isOk := context.Get("id")
		if !isOk {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		withdraws, err := models.LoadWithdraw(context, connection, userId.(int))
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
	if err := validator.Validate(withdraw); err != nil {
		return nil, err
	}
	return &withdraw, nil
}
