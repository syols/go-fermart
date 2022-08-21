package models

import (
	"context"

	"github.com/syols/go-devops/internal/pkg/database"
)

type Withdraw struct {
	Number   string    `json:"order" db:"number" validate:"luhn"`
	Score    int       `json:"sum" db:"score"`
	Uploaded OrderTime `json:"processed_at,omitempty" db:"ctime"`

	UserId int         `json:"-" db:"user_id"`
	Status OrderStatus `json:"-" db:"status" validate:"oneof=REGISTERED NEW INVALID PROCESSING PROCESSED"`
	Action OrderAction `json:"-" db:"action" validate:"oneof=PURCHASE WITHDRAW"`
}

func (w Withdraw) Create(ctx context.Context, connection database.Database) error {
	_, err := connection.Execute(ctx, "order_update.sql", w)
	return err
}

func LoadWithdraw(ctx context.Context, connection database.Database, userId int) (*[]Withdraw, error) {
	withdraw := Withdraw{
		UserId: userId,
		Status: ProcessedOrderStatus,
		Action: WithdrawOrderAction,
	}

	rows, err := connection.Execute(ctx, "user_orders.sql", withdraw)
	if err != nil {
		return nil, err
	}
	return database.ScanAll[Withdraw](*rows)
}
