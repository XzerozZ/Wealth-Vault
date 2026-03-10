package mock

import (
	"context"

	pb "wealth-vault/asset-service/pkg/pb/proto/user"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockUserClient struct {
	pb.UserServiceClient
	mock.Mock
}

func (m *MockUserClient) DeleteAllReferencesByEntityID(
	ctx context.Context,
	req *pb.DeleteByEntityRequest,
	opts ...grpc.CallOption,
) (*pb.DeleteByEntityResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*pb.DeleteByEntityResponse), args.Error(1)
}
