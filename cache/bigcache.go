package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Xarakiri/L0/db"
	"github.com/Xarakiri/L0/schema"
	"github.com/allegro/bigcache"
	"log"
	"time"
)

type BigcacheRepository struct {
	orders *bigcache.BigCache
}

var (
	errOrderNotInCache = errors.New("the order isn't in cache")
)

func NewBigcache() (*BigcacheRepository, error) {
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	if err != nil {
		return nil, err
	}
	return &BigcacheRepository{
		orders: cache,
	}, nil
}

func (r *BigcacheRepository) InsertOrder(order schema.Order) error {
	orderJson, err := json.Marshal(&order)
	if err != nil {
		return err
	}

	return r.orders.Set(order.OrderUid, orderJson)
}

func (r *BigcacheRepository) GetOrder(orderUid string) (schema.Order, error) {
	order, err := r.orders.Get(orderUid)
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return schema.Order{}, errOrderNotInCache
		}
		return schema.Order{}, err
	}

	var o schema.Order
	err = json.Unmarshal(order, &o)
	if err != nil {
		return schema.Order{}, err
	}
	return o, nil
}

func (r *BigcacheRepository) InitCacheFromDB() error {
	orders, err := db.ListOrders(context.Background())
	if err != nil {
		log.Fatalf("Eror while init cache from db: %s\n", err)
	}

	for _, order := range orders {
		err = r.InsertOrder(order)
		if err != nil {
			log.Fatalf("Error while insert order %s\n", err)
			return err
		}
	}

	return nil
}
