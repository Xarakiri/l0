package event

import "github.com/Xarakiri/L0/schema"

type EventStore interface {
	Close()
	PublishOrderCreated(order schema.Order) error
	OnOrderCreated(f func(message OrderCreatedMessage)) error
}

var impl EventStore

func SetEventStore(es EventStore) {
	impl = es
}

func Close() {
	impl.Close()
}

func PublishOrderCreated(order schema.Order) error {
	return impl.PublishOrderCreated(order)
}

func OnOrderCreated(f func(message OrderCreatedMessage)) error {
	return impl.OnOrderCreated(f)
}
