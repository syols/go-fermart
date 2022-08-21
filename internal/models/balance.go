package models

import (
	"context"

	"github.com/syols/go-devops/internal/pkg/database"
)

type Balance struct {
	UserId    int `json:"-" db:"user_id"`
	Current   int `json:"current" db:"current"`
	Withdrawn int `json:"withdrawn" db:"withdrawn"`
}

func CalculateBalance(ctx context.Context, connection database.Database, userId int) (*Balance, error) {
	request := Balance{
		UserId: userId,
	}

	rows, err := connection.Execute(ctx, "user_balance.sql", request)
	if err != nil {
		return nil, err
	}
	return database.ScanOne[Balance](*rows)
}
