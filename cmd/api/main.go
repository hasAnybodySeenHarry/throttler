package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/sony/gobreaker"
	"harry2an.com/throttler/internal/data"
	"harry2an.com/throttler/internal/metrics"
)

type application struct {
	config  config
	wg      sync.WaitGroup
	logger  *log.Logger
	models  data.Models
	clients data.Clients
	cb      *gobreaker.CircuitBreaker
	metrics *metrics.Metrics
}

func main() {
	var cfg config
	loadConfig(&cfg)

	l := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	buckets, users, conn, err := initDependencies(cfg, l)
	if err != nil {
		l.Fatalln(err)
	}
	defer buckets.Close()
	defer users.Close()
	defer conn.Close()

	settings := gobreaker.Settings{
		Name:        "UserClient",
		MaxRequests: 10,
		Interval:    10 * time.Second,
		Timeout:     5 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 5
		},
	}
	cb := gobreaker.NewCircuitBreaker(settings)

	app := application{
		config:  cfg,
		logger:  l,
		models:  data.NewModels(buckets, users),
		clients: data.NewClients(conn),
		cb:      cb,
		metrics: metrics.Register(),
	}

	err = app.serve()
	if err != nil {
		app.logger.Fatal(err)
	}
}
