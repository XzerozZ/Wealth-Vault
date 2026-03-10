package mock

import (
	"context"

	pb "wealth-vault/auth-service/pkg/pb/proto/user"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockUserClient struct {
	pb.UserServiceClient
	mock.Mock
}

func (m *MockUserClient) CreateUser(ctx context.Context, in *pb.CreateUserRequest, opts ...grpc.CallOption) (*pb.UserResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.UserResponse), args.Error(1)
}
