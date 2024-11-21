package main

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"harry2an.com/throttler/internal/jsonlog"
)

func openRedis(logger *jsonlog.Logger, cfg *rd, db int, maxRetries int) (*redis.Client, error) {
	var client *redis.Client
	var err error

	opt := redis.Options{
		Addr:     cfg.addr,
		Password: cfg.password,
		DB:       db,
	}

	client = redis.NewClient(&opt)

	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		_, err := client.Ping(ctx).Result()
		cancel()

		if err == nil {
			return client, nil
		}

		logger.Error(err, map[string]string{
			"message": "Failed to connect to Redis. Retrying in 5 seconds",
		})
		time.Sleep(5 * time.Second)
	}

	return nil, err
}
