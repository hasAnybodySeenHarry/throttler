package main

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func openRedis(cfg *rd, db int, maxRetries int) (*redis.Client, error) {
	var client *redis.Client
	var err error

	log.Printf("Redis server: %s%s", cfg.addr, cfg.password)

	opt := redis.Options{
		Addr:     cfg.addr,
		Password: cfg.password,
		DB:       db,
	}

	client = redis.NewClient(&opt)

	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := client.Ping(ctx).Result()
		if err == nil {
			return client, nil
		}

		log.Printf("Failed to connect to Redis: %v. Retrying in 5 seconds...", err)
		time.Sleep(5 * time.Second)
	}

	return nil, err
}
