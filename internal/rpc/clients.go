package rpc

import (
	"google.golang.org/grpc"
	"harry2an.com/throttler/cmd/proto/users"
)

type Clients struct {
	Users users.UserServiceClient
}

func NewClients(conn *grpc.ClientConn) Clients {
	return Clients{
		Users: users.NewUserServiceClient(conn),
	}
}
