package models

import (
	"context"

	"github.com/syols/go-devops/internal/pkg/database"
)

type Withdraw struct {
	Number   string    `json:"order" db:"number" validate:"luhn"`
	Score    float32   `json:"sum" db:"score"`
	Uploaded OrderTime `json:"processed_at,omitempty" db:"ctime"`

	UserID int         `json:"-" db:"user_id"`
	Status OrderStatus `json:"-" db:"status" validate:"oneof=REGISTERED NEW INVALID PROCESSING PROCESSED"`
	Action OrderAction `json:"-" db:"action" validate:"oneof=PURCHASE WITHDRAW"`
}

func (w Withdraw) Create(ctx context.Context, connection database.Database) error {
	rows, err := connection.Execute(ctx, "order_create.sql", w)
	if err := rows.Err(); err != nil {
		return err
	}
	return err
}

func LoadWithdraw(ctx context.Context, connection database.Database, UserID int) (*[]Withdraw, error) {
	withdraw := Withdraw{
		UserID: UserID,
		Status: ProcessedOrderStatus,
		Action: WithdrawOrderAction,
	}

	rows, err := connection.Execute(ctx, "user_orders.sql", withdraw)
	if err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var values []Withdraw
	for rows.Next() {
		var value Withdraw
		if err := rows.StructScan(&value); err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return &values, nil
}
