package client

import (
	pb "wealth-vault/user-service/pkg/pb/proto/asset"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewAssetClient(host, port string) (pb.AssetServiceClient, error) {
	addr := host + ":" + port
	cred := insecure.NewCredentials()
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(cred))
	if err != nil {
		return nil, err
	}

	return pb.NewAssetServiceClient(conn), nil
}
