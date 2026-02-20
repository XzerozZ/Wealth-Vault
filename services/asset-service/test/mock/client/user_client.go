package mock

import (
	"context"

	pb "wealth-vault/asset-service/pkg/pb/proto/user"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockUserClient struct {
	mock.Mock
}

func (m *MockUserClient) AcceptFriend(ctx context.Context, in *pb.AcceptFriendRequest, opts ...grpc.CallOption) (*pb.FriendResponse, error) {
	panic("unimplemented")
}

func (m *MockUserClient) AddFriend(ctx context.Context, in *pb.FriendRequest, opts ...grpc.CallOption) (*pb.FriendResponse, error) {
	panic("unimplemented")
}

func (m *MockUserClient) AddGroupMember(ctx context.Context, in *pb.AddMemberRequest, opts ...grpc.CallOption) (*pb.ActionResponse, error) {
	panic("unimplemented")
}

func (m *MockUserClient) CreateGroup(ctx context.Context, in *pb.CreateGroupRequest, opts ...grpc.CallOption) (*pb.GroupResponse, error) {
	panic("unimplemented")
}

func (m *MockUserClient) CreateUser(ctx context.Context, in *pb.CreateUserRequest, opts ...grpc.CallOption) (*pb.UserResponse, error) {
	panic("unimplemented")
}

func (m *MockUserClient) DeleteGroup(ctx context.Context, in *pb.DeleteGroupRequest, opts ...grpc.CallOption) (*pb.ActionResponse, error) {
	panic("unimplemented")
}

func (m *MockUserClient) GetAllGroup(ctx context.Context, in *pb.AllGroupRequest, opts ...grpc.CallOption) (*pb.AllGroupResponse, error) {
	panic("unimplemented")
}

func (m *MockUserClient) GetCloseFriends(ctx context.Context, in *pb.GetCloseFriendsRequest, opts ...grpc.CallOption) (*pb.GetCloseFriendsResponse, error) {
	panic("unimplemented")
}

func (m *MockUserClient) GetFriendList(ctx context.Context, in *pb.GetUserByIDRequest, opts ...grpc.CallOption) (*pb.FriendListResponse, error) {
	panic("unimplemented")
}

func (m *MockUserClient) GetGroup(ctx context.Context, in *pb.GetGroupRequest, opts ...grpc.CallOption) (*pb.GroupResponse, error) {
	panic("unimplemented")
}

func (m *MockUserClient) GetGroupMessages(ctx context.Context, in *pb.GetGroupMessagesRequest, opts ...grpc.CallOption) (*pb.GetGroupMessagesResponse, error) {
	panic("unimplemented")
}

func (m *MockUserClient) GetItemSharedTargets(ctx context.Context, in *pb.GetItemSharedTargetsRequest, opts ...grpc.CallOption) (*pb.GetItemSharedTargetsResponse, error) {
	panic("unimplemented")
}

func (m *MockUserClient) GetItemsSharedByFriend(ctx context.Context, in *pb.GetItemsSharedByFriendRequest, opts ...grpc.CallOption) (*pb.GetItemsSharedByFriendResponse, error) {
	panic("unimplemented")
}

func (m *MockUserClient) GetMember(ctx context.Context, in *pb.GetGroupMembersRequest, opts ...grpc.CallOption) (*pb.GetGroupMembersResponse, error) {
	panic("unimplemented")
}

func (m *MockUserClient) GetPendingRequests(ctx context.Context, in *pb.GetUserByIDRequest, opts ...grpc.CallOption) (*pb.FriendListResponse, error) {
	panic("unimplemented")
}

func (m *MockUserClient) GetPrivateMessages(ctx context.Context, in *pb.GetPrivateMessagesRequest, opts ...grpc.CallOption) (*pb.GetPrivateMessagesResponse, error) {
	panic("unimplemented")
}

func (m *MockUserClient) GetSharedItem(ctx context.Context, in *pb.GetGroupItemsRequest, opts ...grpc.CallOption) (*pb.GetGroupItemsResponse, error) {
	panic("unimplemented")
}

func (m *MockUserClient) GetSharedItemIDs(ctx context.Context, in *pb.GetSharedItemIDsRequest, opts ...grpc.CallOption) (*pb.GetSharedItemIDsResponse, error) {
	panic("unimplemented")
}

func (m *MockUserClient) GetSharedIteminFriend(ctx context.Context, in *pb.GetFriendItemRequest, opts ...grpc.CallOption) (*pb.GetFriendItemsResponse, error) {
	panic("unimplemented")
}

func (m *MockUserClient) GetUser(ctx context.Context, in *pb.GetUserByIDRequest, opts ...grpc.CallOption) (*pb.UserResponse, error) {
	panic("unimplemented")
}

func (m *MockUserClient) GrantGroupItemAccess(ctx context.Context, in *pb.GrantAccessRequest, opts ...grpc.CallOption) (*pb.ActionResponse, error) {
	panic("unimplemented")
}

func (m *MockUserClient) LeaveGroup(ctx context.Context, in *pb.LeaveGroupRequest, opts ...grpc.CallOption) (*pb.ActionResponse, error) {
	panic("unimplemented")
}

func (m *MockUserClient) RemoveMember(ctx context.Context, in *pb.RemoveMemberRequest, opts ...grpc.CallOption) (*pb.ActionResponse, error) {
	panic("unimplemented")
}

func (m *MockUserClient) SetCloseFriend(ctx context.Context, in *pb.SetCloseFriendRequest, opts ...grpc.CallOption) (*pb.SetCloseFriendResponse, error) {
	panic("unimplemented")
}

func (m *MockUserClient) ShareItem(ctx context.Context, in *pb.ShareItemRequest, opts ...grpc.CallOption) (*pb.ShareItemResponse, error) {
	panic("unimplemented")
}

func (m *MockUserClient) UnsharedItem(ctx context.Context, in *pb.UnshareItemRequest, opts ...grpc.CallOption) (*pb.ShareItemResponse, error) {
	panic("unimplemented")
}

func (m *MockUserClient) UnsharedIteminFriend(ctx context.Context, in *pb.UnshareItemRequest, opts ...grpc.CallOption) (*pb.ShareItemResponse, error) {
	panic("unimplemented")
}

func (m *MockUserClient) UpdateGroup(ctx context.Context, in *pb.UpdateGroupRequest, opts ...grpc.CallOption) (*pb.GroupResponse, error) {
	panic("unimplemented")
}

func (m *MockUserClient) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest, opts ...grpc.CallOption) (*pb.UserResponse, error) {
	panic("unimplemented")
}

func (m *MockUserClient) DeleteAllReferencesByEntityID(
	ctx context.Context,
	req *pb.DeleteByEntityRequest,
	opts ...grpc.CallOption,
) (*pb.DeleteByEntityResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*pb.DeleteByEntityResponse), args.Error(1)
}
