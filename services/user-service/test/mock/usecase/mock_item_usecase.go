package mock

import (
	"context"
	"wealth-vault/user-service/internal/domain"
	pb "wealth-vault/user-service/pkg/pb/proto/user"

	"github.com/stretchr/testify/mock"
)

type MockItemUsecase struct {
	mock.Mock
}

// GetAllSharedItemIDsByUser implements usecase.ShareItemUsecase.
func (m *MockItemUsecase) GetAllSharedItemIDsByUser(ctx context.Context, req *pb.GetAllSharedItemIDsByUserRequest) (*pb.GetAllSharedItemIDsByUserResponse, error) {
	panic("unimplemented")
}

// AddMemberToGroup implements usecase.GroupItemUsecase.
func (m *MockItemUsecase) AddMemberToGroup(ctx context.Context, req *pb.AddMemberRequest) (*pb.ActionResponse, error) {
	panic("unimplemented")
}

// DeleteAllReferencesByEntityID implements usecase.GroupItemUsecase.
func (m *MockItemUsecase) DeleteAllReferencesByEntityID(ctx context.Context, req *pb.DeleteByEntityRequest) (*pb.DeleteByEntityResponse, error) {
	panic("unimplemented")
}

// GetItemSharedTargets implements usecase.GroupItemUsecase.
func (m *MockItemUsecase) GetItemSharedTargets(ctx context.Context, req *pb.GetItemSharedTargetsRequest) (*pb.GetItemSharedTargetsResponse, error) {
	panic("unimplemented")
}

// GetItemsSharedByFriend implements usecase.GroupItemUsecase.
func (m *MockItemUsecase) GetItemsSharedByFriend(ctx context.Context, req *pb.GetItemsSharedByFriendRequest) (*pb.GetItemsSharedByFriendResponse, error) {
	panic("unimplemented")
}

// GetSharedItemIDs implements usecase.GroupItemUsecase.
func (m *MockItemUsecase) GetSharedItemIDs(ctx context.Context, req *pb.GetSharedItemIDsRequest) (*pb.GetSharedItemIDsResponse, error) {
	panic("unimplemented")
}

// GetSharedIteminFriend implements usecase.GroupItemUsecase.
func (m *MockItemUsecase) GetSharedIteminFriend(ctx context.Context, req *pb.GetFriendItemRequest) (*pb.GetFriendItemsResponse, error) {
	panic("unimplemented")
}

// GetSharedIteminGroup implements usecase.GroupItemUsecase.
func (m *MockItemUsecase) GetSharedIteminGroup(ctx context.Context, req *pb.GetGroupItemsRequest) (*pb.GetGroupItemsResponse, error) {
	panic("unimplemented")
}

// GrantAccess implements usecase.GroupItemUsecase.
func (m *MockItemUsecase) GrantAccess(ctx context.Context, req *pb.GrantAccessRequest) (*pb.ActionResponse, error) {
	panic("unimplemented")
}

// ProcessScheduledEmails implements usecase.GroupItemUsecase.
func (m *MockItemUsecase) ProcessScheduledEmails(ctx context.Context) error {
	panic("unimplemented")
}

// ShareItemtoGroup implements usecase.GroupItemUsecase.
func (m *MockItemUsecase) ShareItem(ctx context.Context, req *pb.ShareItemRequest) (*pb.ShareItemResponse, error) {
	panic("unimplemented")
}

// UnsharedIteminFriend implements usecase.GroupItemUsecase.
func (m *MockItemUsecase) UnsharedIteminFriend(ctx context.Context, req *pb.UnshareItemRequest) (*pb.ShareItemResponse, error) {
	panic("unimplemented")
}

// UnsharedIteminGroup implements usecase.GroupItemUsecase.
func (m *MockItemUsecase) UnsharedIteminGroup(ctx context.Context, req *pb.UnshareItemRequest) (*pb.ShareItemResponse, error) {
	panic("unimplemented")
}

func (m *MockItemUsecase) BatchShareAssets(ctx context.Context, req domain.BatchShareRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}
