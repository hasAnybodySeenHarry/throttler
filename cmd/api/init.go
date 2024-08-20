package main

import (
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

func initDependencies(cfg config, logger *log.Logger) (*redis.Client, *grpc.ClientConn, error) {
	var client *redis.Client
	var conn *grpc.ClientConn
	var rdErr, grpcErr error

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

	if rdErr != nil {
		return nil, nil, rdErr
	}
	if grpcErr != nil {
		return nil, nil, grpcErr
	}

	return client, conn, nil
}
