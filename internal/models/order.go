package models

import (
	"fmt"
	"time"
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

func (t OrderTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02T15:04:05.999999-07:00"))
	return []byte(stamp), nil
}
