package main

import (
	"encoding/json"
	"github.com/Xarakiri/L0/event"
	"github.com/Xarakiri/L0/schema"
	"github.com/Xarakiri/L0/util"
	"github.com/segmentio/ksuid"
	"log"
	"net/http"
)

func createOrderHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Message string `json:"message"`
	}

	// Read parameters
	var o schema.Order
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&o); err != nil {
		util.ResponseError(w, http.StatusBadRequest, "Failed to create order")
		return
	}

	// Setup delivery id
	o.Delivery.Id = ksuid.New().String()
	o.DeliveryId = o.Delivery.Id
	o.PaymentId = o.Payment.Transaction

	// Publish event
	if err := event.PublishOrderCreated(o); err != nil {
		log.Println(err)
	}

	// Return new order
	util.ResponseOk(w, response{Message: "Order send to NATS!"})
}
