package models

import (
	"context"

	"github.com/syols/go-devops/internal/pkg/database"
)

type Purchase struct {
	Number   string    `json:"order" db:"number" validate:"luhn"`
	Score    *int      `json:"accrual,omitempty" db:"score"`
	Uploaded OrderTime `json:"uploaded_at" db:"ctime"`

	UserId int         `json:"-" db:"user_id"`
	Status OrderStatus `json:"status" db:"status" validate:"oneof=REGISTERED NEW INVALID PROCESSING PROCESSED"`
	Action OrderAction `json:"-" db:"action" validate:"oneof=PURCHASE"`
}

func NewPurchase(number string, userId int) Purchase {
	return Purchase{
		Number: number,
		UserId: userId,
		Status: NewOrderStatus,
		Action: PurchaseOrderAction,
	}
}

func LoadPurchase(ctx context.Context, connection database.Database, userId int) (*Purchase, error) {
	purchase := Purchase{
		UserId: userId,
	}

	rows, err := connection.Execute(ctx, "order_select.sql", purchase)
	if err != nil {
		return nil, err
	}
	return database.ScanOne[Purchase](*rows)
}

func LoadPurchases(ctx context.Context, connection database.Database, userId int) (*[]Purchase, error) {
	purchase := Purchase{
		UserId: userId,
		Action: PurchaseOrderAction,
	}

	rows, err := connection.Execute(ctx, "user_orders.sql", purchase)
	if err != nil {
		return nil, err
	}
	return database.ScanAll[Purchase](*rows)
}

func (p Purchase) Create(ctx context.Context, connection database.Database) error {
	_, err := connection.Execute(ctx, "order_update.sql", p)
	return err
}

func (p Purchase) Update(ctx context.Context, connection database.Database) error {
	_, err := connection.Execute(ctx, "order_update.sql", p)
	return err
}
