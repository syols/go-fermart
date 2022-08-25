package models

import (
	"context"

	"github.com/syols/go-devops/internal/pkg"
)

type Balance struct {
	UserID    int     `json:"-" db:"user_id"`
	Current   float32 `json:"current" db:"current"`
	Withdrawn float32 `json:"withdrawn" db:"withdrawn"`
}

func CalculateBalance(ctx context.Context, connection pkg.Database, userID int) (*Balance, error) {
	request := Balance{
		UserID: userID,
	}

	rows, err := connection.Execute(ctx, "user_balance.sql", request)
	if err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var value Balance
	if rows.Next() {
		if err := rows.StructScan(&value); err != nil {
			return nil, err
		}
	}

	return &value, nil
}
