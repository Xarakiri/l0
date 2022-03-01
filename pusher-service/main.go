package main

import (
	"fmt"
	"github.com/Xarakiri/L0/db"
	"github.com/Xarakiri/L0/event"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
	"log"
	"net/http"
	"time"
)

const clientID = "client-pusher"

type Config struct {
	PostgresDB       string `envconfig:"POSTGRES_DB"`
	PostgresUser     string `envconfig:"POSTGRES_USER"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD"`
	NatsAddress      string `envconfig:"NATS_ADDRESS"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to PostgreSQL
	retry.ForeverSleep(2*time.Second, func(attempt int) error {
		addr := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)
		repo, err := db.NewPostgres(addr)
		if err != nil {
			log.Println(err)
			return nil
		}
		db.SetRepository(repo)
		return nil
	})
	defer db.Close()

	// Connect to NATS
	retry.ForeverSleep(2*time.Second, func(_ int) error {
		es, err := event.NewNats(fmt.Sprintf("nats://%s", cfg.NatsAddress), clientID)
		if err != nil {
			log.Println(err)
			return err
		}

		err = es.OnOrderCreated(func(m event.OrderCreatedMessage) {
			log.Printf("Order recived: %v\n", m)
		})
		if err != nil {
			log.Println(err)
			return err
		}

		event.SetEventStore(es)
		return nil
	})
	defer event.Close()

	// Run HTTP server
	router := newRouter()
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}

}

func newRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/orders", createOrderHandler).
		Methods(http.MethodPost)
	router.Use(mux.CORSMethodMiddleware(router))
	return router
}
