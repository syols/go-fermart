package models

import (
	"context"

	"github.com/syols/go-devops/internal/pkg/storage"
)

type Purchase struct {
	Number   string    `json:"number" db:"number" validate:"luhn"`
	Score    *float32  `json:"accrual,omitempty" db:"score"`
	Uploaded OrderTime `json:"uploaded_at" db:"ctime"`

	UserID int         `json:"-" db:"user_id"`
	Status OrderStatus `json:"status" db:"status" validate:"oneof=REGISTERED NEW INVALID PROCESSING PROCESSED"`
	Action OrderAction `json:"-" db:"action" validate:"oneof=PURCHASE"`
}

func NewPurchase(number string, UserID int) Purchase {
	return Purchase{
		Number: number,
		UserID: UserID,
		Status: NewOrderStatus,
		Action: PurchaseOrderAction,
	}
}

func LoadPurchase(ctx context.Context, connection storage.Database, number string) (*Purchase, error) {
	purchase := Purchase{
		Number: number,
	}

	rows, err := connection.Execute(ctx, "order_select.sql", purchase)
	if err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var value Purchase
	if rows.Next() {
		if err := rows.StructScan(&value); err != nil {
			return nil, err
		}
		return &value, nil
	}
	return nil, nil
}

func LoadPurchases(ctx context.Context, connection storage.Database, UserID int) (*[]Purchase, error) {
	purchase := Purchase{
		UserID: UserID,
		Action: PurchaseOrderAction,
	}

	rows, err := connection.Execute(ctx, "user_orders.sql", purchase)
	if err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var values []Purchase
	for rows.Next() {
		var value Purchase
		if err := rows.StructScan(&value); err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return &values, nil
}

func (p Purchase) Create(ctx context.Context, connection storage.Database) error {
	rows, err := connection.Execute(ctx, "order_create.sql", p)
	if err := rows.Err(); err != nil {
		return err
	}
	return err
}

func (p Purchase) Update(ctx context.Context, connection storage.Database) error {
	rows, err := connection.Execute(ctx, "order_update.sql", p)
	if err := rows.Err(); err != nil {
		return err
	}
	return err
}
