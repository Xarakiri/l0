package main

import (
	"context"
	"github.com/Xarakiri/L0/cache"
	"github.com/Xarakiri/L0/db"
	"github.com/Xarakiri/L0/event"
	"github.com/Xarakiri/L0/schema"
	"github.com/Xarakiri/L0/util"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func getOrderHandler(w http.ResponseWriter, r *http.Request) {
	// Read parameters
	vars := mux.Vars(r)
	orderId := vars["id"]

	// Find order in cache
	order, err := cache.GetOrder(orderId)
	if err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusNotFound, "Could not find the order")
		return
	}

	util.ResponseOk(w, order)
}

func onOrderCreated(m event.OrderCreatedMessage) {
	order := schema.Order{
		OrderUid:          m.Order.OrderUid,
		TrackNumber:       m.Order.TrackNumber,
		Entry:             m.Order.Entry,
		Delivery:          m.Order.Delivery,
		Payment:           m.Order.Payment,
		Items:             m.Order.Items,
		Locale:            m.Order.Locale,
		InternalSignature: m.Order.InternalSignature,
		CustomerId:        m.Order.CustomerId,
		DeliveryService:   m.Order.DeliveryService,
		Shardkey:          m.Order.Shardkey,
		SmId:              m.Order.SmId,
		DateCreated:       m.Order.DateCreated,
		OofShard:          m.Order.OofShard,
		DeliveryId:        m.Order.DeliveryId,
		PaymentId:         m.Order.PaymentId,
	}
	// Save order in db
	if err := db.InsertOrder(context.Background(), order); err != nil {
		log.Println(err)
		return
	}
	// Save order in cache
	if err := cache.InsertOrder(order); err != nil {
		log.Println(err)
		return
	}
	log.Println("order created")
	log.Println(order)
}
