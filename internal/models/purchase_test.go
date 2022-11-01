package models

import (
	"github.com/go-playground/assert/v2"
	"testing"
)


func TestPurchase(t *testing.T) {
	first := NewPurchase("number", 0)
	second := Purchase{
		Number: "number",
		UserID: 0,
		Status: NewOrderStatus,
		Action: PurchaseOrderAction,
	}
	assert.Equal(t, first, second)
}
