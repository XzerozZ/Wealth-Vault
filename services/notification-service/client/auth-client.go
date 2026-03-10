package client

import (
	pb "wealth-vault/notification-service/pkg/pb/proto/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewAuthClient(host, port string) (pb.AuthServiceClient, error) {
	addr := host + ":" + port
	cred := insecure.NewCredentials()
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(cred))
	if err != nil {
		return nil, err
	}

	return pb.NewAuthServiceClient(conn), nil
}
