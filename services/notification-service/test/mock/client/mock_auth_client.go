package mocks

import (
	"context"

	pb "wealth-vault/notification-service/pkg/pb/proto/auth"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockAuthClient struct {
	pb.AuthServiceClient
	mock.Mock
}

func (m *MockAuthClient) GetProviderAccount(ctx context.Context, in *pb.GetProviderAccountRequest, opts ...grpc.CallOption) (*pb.GetProviderAccountsResponse, error) {
	args := m.Called(ctx, in)
	var res *pb.GetProviderAccountsResponse
	if args.Get(0) != nil {
		res = args.Get(0).(*pb.GetProviderAccountsResponse)
	}
	return res, args.Error(1)
}
