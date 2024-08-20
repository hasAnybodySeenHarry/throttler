package main

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"harry2an.com/throttler/cmd/proto/users"
)

type clients struct {
	users users.UserServiceClient
}

func openGRPC(addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:50051", addr), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func New(conn *grpc.ClientConn) clients {
	return clients{
		users: users.NewUserServiceClient(conn),
	}
}
