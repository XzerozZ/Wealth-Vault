package client

import (
	pb "wealth-vault/auth-service/pkg/pb/proto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewUserClient(host, port string) (pb.UserServiceClient, error) {
	addr := host + ":" + port
	cred := insecure.NewCredentials()
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(cred))
	if err != nil {
		return nil, err
	}

	return pb.NewUserServiceClient(conn), nil
}
