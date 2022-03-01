package cache

import "github.com/Xarakiri/L0/schema"

type Repository interface {
	InsertOrder(order schema.Order) error
	GetOrder(orderUid string) (schema.Order, error)
}

var impl Repository

func SetRepository(repository Repository) {
	impl = repository
}

func InsertOrder(order schema.Order) error {
	return impl.InsertOrder(order)
}

func GetOrder(orderUid string) (schema.Order, error) {
	return impl.GetOrder(orderUid)
}
