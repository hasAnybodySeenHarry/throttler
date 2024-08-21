package main

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func openGRPC(addr string) (*grpc.ClientConn, error) {
	if addr == "" {
		addr = "development-expenses"
	}
	server := fmt.Sprintf("%s:50051", addr)
	conn, err := grpc.NewClient(server, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
