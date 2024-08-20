package main

import (
	"log"
	"os"
	"sync"

	"harry2an.com/throttler/internal/data"
)

type application struct {
	config config
	wg     sync.WaitGroup
	logger *log.Logger
	models data.Models
}

func main() {
	var cfg config
	loadConfig(&cfg)

	l := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	client, err := initDependencies(cfg, l)
	if err != nil {
		l.Fatalln(err)
	}
	defer client.Close()

	app := application{
		config: cfg,
		logger: l,
		models: data.New(client),
	}

	err = app.serve()
	if err != nil {
		app.logger.Fatal(err)
	}
}
