package db

import (
	"context"
	"github.com/Xarakiri/L0/schema"
)

type Repository interface {
	Close()
	InsertOrder(ctx context.Context, order schema.Order) error
	GetOrder(ctx context.Context, orderUid string) (schema.Order, error)
	ListOrders(ctx context.Context) ([]schema.Order, error)
}

var impl Repository

func SetRepository(repository Repository) {
	impl = repository
}

func Close() {
	impl.Close()
}

func InsertOrder(ctx context.Context, order schema.Order) error {
	return impl.InsertOrder(ctx, order)
}

func GetOrder(ctx context.Context, orderUid string) (schema.Order, error) {
	return impl.GetOrder(ctx, orderUid)
}

func ListOrders(ctx context.Context) ([]schema.Order, error) {
	return impl.ListOrders(ctx)
}
