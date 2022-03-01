package event

import (
	"bytes"
	"encoding/gob"
	"github.com/Xarakiri/L0/schema"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"log"
)

type NatsEventStore struct {
	nc                       stan.Conn
	orderCreatedSubscription stan.Subscription
	orderCreatedChan         chan OrderCreatedMessage
}

const durableName = "queue"

func NewNats(url string, clientId string) (*NatsEventStore, error) {
	// Connect to NATS
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	// Connect to NATS Streaming server
	sc, err := stan.Connect("test-cluster", clientId, stan.NatsConn(nc), stan.MaxPubAcksInflight(1),
		stan.SetConnectionLostHandler(func(_ stan.Conn, err error) {
			log.Fatalf("Connection lost. Reason: #{reason}")
		}))
	if err != nil {
		return nil, err
	}

	return &NatsEventStore{nc: sc}, nil
}

func (es *NatsEventStore) Close() {
	if es.nc != nil {
		es.nc.Close()
	}
	if es.orderCreatedSubscription != nil {
		es.orderCreatedSubscription.Unsubscribe()
	}
	close(es.orderCreatedChan)
}

func (es *NatsEventStore) PublishOrderCreated(order schema.Order) error {
	m := OrderCreatedMessage{schema.Order{
		OrderUid:          order.OrderUid,
		TrackNumber:       order.TrackNumber,
		Entry:             order.Entry,
		Delivery:          order.Delivery,
		Payment:           order.Payment,
		Items:             order.Items,
		Locale:            order.Locale,
		InternalSignature: order.InternalSignature,
		CustomerId:        order.CustomerId,
		DeliveryService:   order.DeliveryService,
		Shardkey:          order.Shardkey,
		SmId:              order.SmId,
		DateCreated:       order.DateCreated,
		OofShard:          order.OofShard,
		DeliveryId:        order.DeliveryId,
		PaymentId:         order.PaymentId,
	}}
	data, err := es.writeMessage(&m)
	if err != nil {
		return err
	}
	return es.nc.Publish(m.Key(), data)
}

func (es *NatsEventStore) writeMessage(o Order) ([]byte, error) {
	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(o)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (es *NatsEventStore) OnOrderCreated(f func(message OrderCreatedMessage)) (err error) {
	m := OrderCreatedMessage{}
	es.orderCreatedSubscription, err = es.nc.Subscribe(m.Key(), func(msg *stan.Msg) {
		if err := es.readMessage(msg.Data, &m); err != nil {
			log.Fatal(err)
		}
		f(m)
	}, stan.DurableName(durableName))
	return
}

func (es *NatsEventStore) readMessage(data []byte, m interface{}) error {
	b := bytes.Buffer{}
	b.Write(data)
	return gob.NewDecoder(&b).Decode(m)
}
