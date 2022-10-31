package models

import (
	"context"

	"github.com/syols/go-devops/internal/pkg"
)

type Purchase struct {
	Number   string    `json:"number" db:"number" validate:"luhn"`
	Score    *float32  `json:"accrual,omitempty" db:"score"`
	Uploaded OrderTime `json:"uploaded_at" db:"ctime"`

	UserID int         `json:"-" db:"user_id"`
	Status OrderStatus `json:"status" db:"status" validate:"oneof=REGISTERED NEW INVALID PROCESSING PROCESSED"`
	Action OrderAction `json:"-" db:"action" validate:"oneof=PURCHASE"`
}

func NewPurchase(number string, userID int) Purchase {
	return Purchase{
		Number: number,
		UserID: userID,
		Status: NewOrderStatus,
		Action: PurchaseOrderAction,
	}
}

func (p *Purchase) Create(ctx context.Context, db pkg.Database) error {
	rows, err := db.Execute(ctx, "order_create.sql", p)
	if err := rows.Err(); err != nil {
		return err
	}
	return err
}

func (p *Purchase) Update(ctx context.Context, db pkg.Database) error {
	rows, err := db.Execute(ctx, "order_update.sql", p)
	if err := rows.Err(); err != nil {
		return err
	}
	return err
}

func LoadPurchase(ctx context.Context, db pkg.Database, number string) (*Purchase, error) {
	purchase := Purchase{
		Number: number,
	}

	rows, err := db.Execute(ctx, "order_select.sql", purchase)
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

func LoadPurchases(ctx context.Context, db pkg.Database, userID int) (*[]Purchase, error) {
	purchase := Purchase{
		UserID: userID,
		Action: PurchaseOrderAction,
	}

	rows, err := db.Execute(ctx, "user_orders.sql", purchase)
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
