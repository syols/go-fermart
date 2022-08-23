package models

import (
	"context"
	"fmt"
	"time"

	"github.com/syols/go-devops/internal/pkg/database"
)

type OrderStatus string
type OrderAction string

type OrderTime time.Time

const (
	RegisteredOrderStatus OrderStatus = "REGISTERED"
	NewOrderStatus        OrderStatus = "NEW"
	InvalidOrderStatus    OrderStatus = "INVALID"
	ProcessingOrderStatus OrderStatus = "PROCESSING"
	ProcessedOrderStatus  OrderStatus = "PROCESSED"
)

const (
	PurchaseOrderAction OrderAction = "PURCHASE"
	WithdrawOrderAction OrderAction = "WITHDRAW"
)

type Order struct {
	Number string      `json:"order" db:"number" validate:"luhn"`
	Score  *float32    `json:"accrual,omitempty" db:"score"`
	Status OrderStatus `json:"status" db:"status" validate:"oneof=REGISTERED NEW INVALID PROCESSING PROCESSED"`
}

func (p Order) Update(ctx context.Context, connection database.Database) error {
	rows, err := connection.Execute(ctx, "order_update.sql", p)
	if err := rows.Err(); err != nil {
		return err
	}
	return err
}

func (t OrderTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02T15:04:05.999999-07:00"))
	return []byte(stamp), nil
}
