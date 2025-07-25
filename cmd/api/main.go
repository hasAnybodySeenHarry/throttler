package main

import (
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/sony/gobreaker"
	"harry2an.com/throttler/internal/data"
	rpc "harry2an.com/throttler/internal/grpc"
	"harry2an.com/throttler/internal/jsonlog"
	"harry2an.com/throttler/internal/metrics"
)

type application struct {
	config  config
	wg      sync.WaitGroup
	logger  *jsonlog.Logger
	models  data.Models
	clients rpc.Clients
	cb      *gobreaker.CircuitBreaker
	metrics *metrics.Metrics
	randGen *rand.Rand
}

func main() {

	var cfg config
	loadConfig(&cfg)

	l := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	buckets, users, conn, err := initDependencies(cfg, l)
	if err != nil {
		l.Fatal(err, nil)
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
		clients: rpc.NewClients(conn),
		cb:      cb,
		metrics: metrics.Register(),
		randGen: rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	err = app.serve()
	if err != nil {
		app.logger.Fatal(err, nil)
	}
}
