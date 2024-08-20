package main

import (
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

func initDependencies(cfg config, logger *log.Logger) (*redis.Client, error) {
	var client *redis.Client
	var rdErr error

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		client, rdErr = openRedis(&cfg.redis, 6)
		if rdErr != nil {
			logger.Printf("Failed to connect to Redis: %v", rdErr)
		} else {
			logger.Println("Successfully connected to Redis")
		}
	}()

	wg.Wait()

	if rdErr != nil {
		return nil, rdErr
	}

	return client, nil
}
