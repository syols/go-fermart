package models

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joeljunstrom/go-luhn"
	"github.com/syols/go-devops/internal/pkg/database"
)

type Order struct {
	Number   string    `json:"number" db:"order_id" validate:"luhn"`
	UserId   int       `json:"-" db:"user_id"`
	Accrual  *int      `json:"accrual,omitempty" db:"accrual"`
	Status   string    `json:"status" db:"order_status" validate:"oneof=REGISTERED NEW INVALID PROCESSING PROCESSED"`
	Uploaded time.Time `json:"uploaded_at" db:"uploaded"`
}

func NewOrder(number string, userId int) Order {
	return Order{
		Number: number,
		UserId: userId,
		Status: "NEW",
	}
}

func (order *Order) Validate() error {
	validate := validator.New()
	err := validate.RegisterValidation("luhn", func(fl validator.FieldLevel) bool {
		number, ok := fl.Field().Interface().(string)
		if ok {
			return luhn.Valid(number)
		}
		return false
	})

	if err != nil {
		return err
	}
	return validate.Struct(order)
}

func (order Order) Create(ctx context.Context, connection database.Connection) error {
	_, err := connection.Execute(ctx, database.OrderCreateQuery, order)
	return err
}

func (order Order) Update(ctx context.Context, connection database.Connection) error {
	_, err := connection.Execute(ctx, database.OrderUpdateQuery, order)
	return err
}

func (order Order) Select(ctx context.Context, connection database.Connection) (*Order, error) {
	rows, err := connection.Execute(ctx, database.OrderSelectQuery, order)
	if err != nil {
		return nil, err
	}
	return database.ScanOne[Order](*rows)
}

func (order Order) UserOrders(ctx context.Context, connection database.Connection) (*[]Order, error) {
	rows, err := connection.Execute(ctx, database.UserOrdersSelectQuery, order)
	if err != nil {
		return nil, err
	}
	return database.ScanAll[Order](*rows)
}
