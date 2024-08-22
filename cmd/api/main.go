package main

import (
	"log"
	"os"
	"sync"

	"harry2an.com/throttler/internal/data"
)

type application struct {
	config  config
	wg      sync.WaitGroup
	logger  *log.Logger
	models  data.Models
	clients data.Clients
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

	app := application{
		config:  cfg,
		logger:  l,
		models:  data.NewModels(buckets, users),
		clients: data.NewClients(conn),
	}

	err = app.serve()
	if err != nil {
		app.logger.Fatal(err)
	}
}
