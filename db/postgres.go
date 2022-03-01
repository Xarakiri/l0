package db

import (
	"context"
	"database/sql"
	"github.com/Xarakiri/L0/schema"
	_ "github.com/lib/pq"
	"log"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgres(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{
		db,
	}, nil
}

func (r *PostgresRepository) Close() {
	if err := r.db.Close(); err != nil {
		log.Fatal(err)
	}
}

func (r *PostgresRepository) InsertOrder(ctx context.Context, order schema.Order) error {
	// Create a helper function for preparing failure results.
	fail := func(err error) error {
		return err
	}

	// Get a Tx for making transaction requests.
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fail(err)
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	// 1. Insert delivery
	delivery := order.Delivery
	_, err1 := tx.ExecContext(ctx, "CALL set_delivery($1, $2, $3, $4, $5, $6, $7, $8)", delivery.Id, delivery.Name,
		delivery.Phone, delivery.Zip, delivery.City, delivery.Address, delivery.Region, delivery.Email)
	if err1 != nil {
		return fail(err1)
	}

	// 2. Insert payment
	payment := order.Payment
	_, err2 := tx.ExecContext(ctx, "CALL set_payment($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", payment.Transaction,
		payment.RequestID, payment.Currency, payment.Provider, payment.Amount, payment.PaymentDt, payment.Bank,
		payment.DeliveryCost, payment.GoodsTotal, payment.CustomFee)
	if err2 != nil {
		return fail(err2)
	}

	// 3. Set items
	for _, item := range order.Items {
		_, err3 := tx.ExecContext(ctx, "CALL set_item($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)", item.Rid, item.ChrtID,
			item.TrackNumber, item.Price, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand,
			item.Status)
		if err3 != nil {
			return fail(err3)
		}
	}

	// 4. Set order
	_, err4 := tx.ExecContext(ctx, "CALL set_order($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)", order.OrderUid,
		order.TrackNumber, order.Entry, order.Delivery.Id, order.Payment.Transaction, order.Locale,
		order.InternalSignature, order.CustomerId, order.DeliveryService, order.Shardkey, order.SmId,
		order.DateCreated, order.OofShard)
	if err4 != nil {
		return fail(err4)
	}

	// 5. Set order_item
	for _, item := range order.Items {
		_, err5 := tx.ExecContext(ctx, "CALL set_order_item($1, $2)", order.OrderUid, item.Rid)
		if err5 != nil {
			return fail(err5)
		}
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	return nil
}

func (r *PostgresRepository) getDelivery(ctx context.Context, deliveryId string) (schema.Delivery, error) {
	rows, err := r.db.Query(`SELECT * FROM delivery WHERE id=$1`, deliveryId)
	if err != nil {
		return schema.Delivery{}, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	delivery := schema.Delivery{}
	rows.Next()
	if err = rows.Scan(&delivery.Id, &delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City,
		&delivery.Address, &delivery.Region, &delivery.Email); err != nil {
		return schema.Delivery{}, err
	}
	if err = rows.Err(); err != nil {
		return schema.Delivery{}, err
	}
	return delivery, nil
}

func (r *PostgresRepository) getItem(ctx context.Context, itemId string) (schema.Item, error) {
	rows, err := r.db.Query(`SELECT * FROM item WHERE rid=$1`, itemId)
	if err != nil {
		return schema.Item{}, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	item := schema.Item{}
	rows.Next()
	if err = rows.Scan(&item.Rid, &item.ChrtID, &item.TrackNumber, &item.Price, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status); err != nil {
		return schema.Item{}, err
	}
	if err = rows.Err(); err != nil {
		return schema.Item{}, err
	}
	return item, nil
}

func (r *PostgresRepository) getPayment(ctx context.Context, paymentId string) (schema.Payment, error) {
	rows, err := r.db.Query(`SELECT * FROM payment WHERE transaction=$1`, paymentId)
	if err != nil {
		return schema.Payment{}, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	payment := schema.Payment{}
	rows.Next()
	if err = rows.Scan(&payment.Transaction, &payment.RequestID, &payment.Currency, &payment.Provider, &payment.Amount,
		&payment.PaymentDt, &payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee); err != nil {
		return schema.Payment{}, err
	}
	if err = rows.Err(); err != nil {
		return schema.Payment{}, err
	}
	return payment, nil
}

func (r *PostgresRepository) ListOrders(ctx context.Context) ([]schema.Order, error) {
	rows, err := r.db.Query(`SELECT * FROM "order"`)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Parse all rows into an array of Orders
	var orders []schema.Order

	for rows.Next() {
		order := schema.Order{}
		if err = rows.Scan(&order.OrderUid, &order.TrackNumber, &order.Entry, &order.DeliveryId, &order.PaymentId, &order.Locale, &order.InternalSignature, &order.CustomerId, &order.DeliveryService, &order.Shardkey, &order.SmId, &order.DateCreated, &order.OofShard); err == nil {
			// Get Delivery
			order.Delivery, err = r.getDelivery(ctx, order.DeliveryId)
			if err != nil {
				return nil, err
			}

			// Get Payment
			order.Payment, err = r.getPayment(ctx, order.PaymentId)
			if err != nil {
				return nil, err
			}

			// Get Items
			order.Items, err = r.listItems(ctx, order.OrderUid)
			if err != nil {
				return nil, err
			}

			orders = append(orders, order)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *PostgresRepository) listItems(ctx context.Context, orderUid string) ([]schema.Item, error) {
	// Get list of items
	var items []schema.Item
	rows, err := r.db.Query(`SELECT * FROM item i JOIN (SELECT * FROM order_item WHERE order_id=$1) t2 ON i.rid=t2.item_id;`, orderUid)
	for rows.Next() {
		item := schema.Item{}
		if err = rows.Scan(&item.Rid, &item.ChrtID, &item.TrackNumber, &item.Price, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status); err == nil {
			items = append(items, item)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *PostgresRepository) GetOrder(ctx context.Context, orderUid string) (schema.Order, error) {
	rows, err := r.db.Query(`SELECT * FROM "order" WHERE order_uid=$1`, orderUid)
	if err != nil {
		return schema.Order{}, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Parse
	order := schema.Order{}
	rows.Next()
	if err = rows.Scan(&order.OrderUid, &order.TrackNumber, &order.Entry, &order.DeliveryId, &order.PaymentId, &order.Locale, &order.InternalSignature, &order.CustomerId, &order.DeliveryService, &order.Shardkey, &order.SmId, &order.DateCreated, &order.OofShard); err != nil {
		return order, err
	}
	if err = rows.Err(); err != nil {
		return order, err
	}
	order.Delivery, err = r.getDelivery(ctx, order.DeliveryId)
	if err != nil {
		return order, err
	}
	order.Payment, err = r.getPayment(ctx, order.PaymentId)
	if err != nil {
		return order, err
	}

	order.Items, err = r.listItems(ctx, order.OrderUid)
	if err != nil {
		return schema.Order{}, err
	}

	return order, nil
}

func (r *PostgresRepository) insertOrderItem(ctx context.Context, orderId string, itemId string) error {
	_, err := r.db.Exec("CALL set_order_item($1, $2)", orderId, itemId)
	return err
}

func (r *PostgresRepository) insertDelivery(ctx context.Context, delivery schema.Delivery) error {
	_, err := r.db.Exec("CALL set_delivery($1, $2, $3, $4, $5, $6, $7, $8)", delivery.Id, delivery.Name,
		delivery.Phone, delivery.Zip, delivery.City, delivery.Address, delivery.Region, delivery.Email)
	return err
}

func (r *PostgresRepository) insertPayment(ctx context.Context, payment schema.Payment) error {
	_, err := r.db.Exec("CALL set_payment($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", payment.Transaction,
		payment.RequestID, payment.Currency, payment.Provider, payment.Amount, payment.PaymentDt, payment.Bank,
		payment.DeliveryCost, payment.GoodsTotal, payment.CustomFee)
	return err
}

func (r *PostgresRepository) insertItem(ctx context.Context, item schema.Item) error {
	_, err := r.db.Exec("CALL set_item($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)", item.Rid, item.ChrtID,
		item.TrackNumber, item.Price, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand,
		item.Status)
	return err
}
