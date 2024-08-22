package main

import (
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

func initDependencies(cfg config, logger *log.Logger) (buckets *redis.Client, users *redis.Client, conn *grpc.ClientConn, err error) {
	var bucketsErr, usersErr, grpcErr error
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		buckets, bucketsErr = openRedis(&cfg.redis, 0, 6)
		if bucketsErr != nil {
			logger.Printf("Failed to connect to Redis for buckets: %v", bucketsErr)
		} else {
			logger.Println("Successfully connected to Redis")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		users, usersErr = openRedis(&cfg.redis, 1, 6)
		if usersErr != nil {
			logger.Printf("Failed to connect to Redis for users: %v", usersErr)
		} else {
			logger.Println("Successfully connected to Redis")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		conn, grpcErr = openGRPC(cfg.grpcAddr)
		if grpcErr != nil {
			log.Printf("Failed to connect to the gRPC server: %v", grpcErr)
		} else {
			logger.Println("Successfully connected to the gRPC server")
		}
	}()

	wg.Wait()

	if bucketsErr != nil {
		return nil, nil, nil, bucketsErr
	}
	if usersErr != nil {
		return nil, nil, nil, usersErr
	}
	if grpcErr != nil {
		return nil, nil, nil, grpcErr
	}

	return buckets, users, conn, nil
}
