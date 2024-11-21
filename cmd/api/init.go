package main

import (
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"harry2an.com/throttler/internal/jsonlog"
)

func initDependencies(cfg config, logger *jsonlog.Logger) (buckets *redis.Client, users *redis.Client, conn *grpc.ClientConn, err error) {
	var bucketsErr, usersErr, grpcErr error
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		buckets, bucketsErr = openRedis(&cfg.redis, 0, 6)
		if bucketsErr != nil {
			logger.Error(fmt.Errorf("failed to connect to Redis for buckets: %v", bucketsErr), nil)
		} else {
			logger.Info("Successfully connected to Redis", nil)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		users, usersErr = openRedis(&cfg.redis, 1, 6)
		if usersErr != nil {
			logger.Error(fmt.Errorf("failed to connect to Redis for users: %v", usersErr), nil)
		} else {
			logger.Info("Successfully connected to Redis", nil)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		conn, grpcErr = openGRPC(cfg.grpcAddr)
		if grpcErr != nil {
			logger.Error(fmt.Errorf("failed to connect to the gRPC server: %v", grpcErr), nil)
		} else {
			logger.Info("Successfully connected to the gRPC server", nil)
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
