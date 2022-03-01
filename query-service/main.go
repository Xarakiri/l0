package main

import (
	"fmt"
	"github.com/Xarakiri/L0/cache"
	"github.com/Xarakiri/L0/db"
	"github.com/Xarakiri/L0/event"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
	"log"
	"net/http"
	"time"
)

type Config struct {
	PostgresDB       string `envconfig:"POSTGRES_DB"`
	PostgresUser     string `envconfig:"POSTGRES_USER"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD"`
	NatsAddress      string `envconfig:"NATS_ADDRESS"`
}

const clientID = "client-query"

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

	// Init cache
	repo, err := cache.NewBigcache()
	if err != nil {
		log.Println(err)
	}
	cache.SetRepository(repo)
	if err := repo.InitCacheFromDB(); err != nil {
		log.Fatal(err)
	}

	// Connect to Nats
	retry.ForeverSleep(2*time.Second, func(_ int) error {
		es, err := event.NewNats(fmt.Sprintf("nats://%s", cfg.NatsAddress), clientID)
		if err != nil {
			log.Println(err)
			return err
		}
		err = es.OnOrderCreated(onOrderCreated)
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
	log.Println("RUN SERVER")
}

func newRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/orders/{id}", getOrderHandler).
		Methods(http.MethodGet)
	router.Use(mux.CORSMethodMiddleware(router))
	return router
}
