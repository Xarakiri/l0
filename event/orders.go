package event

import "github.com/Xarakiri/L0/schema"

type Order interface {
	Key() string
}

type OrderCreatedMessage struct {
	Order schema.Order
}

func (m *OrderCreatedMessage) Key() string {
	return "ORDERS.created"
}
