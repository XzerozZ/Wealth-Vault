package grpc

import (
	"context"

	usecase "wealth-vault/user-service/internal/usecase/interface"

	pb "wealth-vault/user-service/pkg/pb/proto/user"
)

type UserGRPCHandler struct {
	pb.UnimplementedUserServiceServer
	usecase   usecase.UserUsecase
	gusecase  usecase.GroupUsecase
	giusecase usecase.ShareItemUsecase
	musecase  usecase.MessageUsecase
}

func NewUserGRPCHandler(u usecase.UserUsecase, g usecase.GroupUsecase, i usecase.ShareItemUsecase, m usecase.MessageUsecase) *UserGRPCHandler {
	return &UserGRPCHandler{
		usecase:   u,
		gusecase:  g,
		giusecase: i,
		musecase:  m,
	}
}

func (h *UserGRPCHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	res, err := h.usecase.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) GetUser(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.UserResponse, error) {
	res, err := h.usecase.GetUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) GetUsersByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.UserInfoResponse, error) {
	res, err := h.usecase.GetUsersByEmail(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	res, err := h.usecase.UpdateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) GetFriendList(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.FriendListResponse, error) {
	res, err := h.usecase.GetFriendList(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) AddFriend(ctx context.Context, req *pb.FriendRequest) (*pb.FriendResponse, error) {
	res, err := h.usecase.AddFriend(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (*pb.GroupResponse, error) {
	res, err := h.gusecase.CreateGroup(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) GetMember(ctx context.Context, req *pb.GetGroupMembersRequest) (*pb.GetGroupMembersResponse, error) {
	res, err := h.gusecase.GetMember(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) GetAllGroup(ctx context.Context, req *pb.AllGroupRequest) (*pb.AllGroupResponse, error) {
	res, err := h.gusecase.AllGetGroup(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) GetGroup(ctx context.Context, req *pb.GetGroupRequest) (*pb.GroupResponse, error) {
	res, err := h.gusecase.GetGroup(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) UpdateGroup(ctx context.Context, req *pb.UpdateGroupRequest) (*pb.GroupResponse, error) {
	res, err := h.gusecase.UpdateGroup(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) ShareItem(ctx context.Context, req *pb.ShareItemRequest) (*pb.ShareItemResponse, error) {
	res, err := h.giusecase.ShareItem(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) GetSharedItem(ctx context.Context, req *pb.GetGroupItemsRequest) (*pb.GetGroupItemsResponse, error) {
	res, err := h.giusecase.GetSharedIteminGroup(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) GetSharedIteminFriend(ctx context.Context, req *pb.GetFriendItemRequest) (*pb.GetFriendItemsResponse, error) {
	res, err := h.giusecase.GetSharedIteminFriend(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) UnsharedItem(ctx context.Context, req *pb.UnshareItemRequest) (*pb.ShareItemResponse, error) {
	res, err := h.giusecase.UnsharedIteminGroup(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) UnsharedIteminFriend(ctx context.Context, req *pb.UnshareItemRequest) (*pb.ShareItemResponse, error) {
	res, err := h.giusecase.UnsharedIteminFriend(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) DeleteAllReferencesByEntityID(ctx context.Context, req *pb.DeleteByEntityRequest) (*pb.DeleteByEntityResponse, error) {
	res, err := h.giusecase.DeleteAllReferencesByEntityID(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) AddGroupMember(ctx context.Context, req *pb.AddMemberRequest) (*pb.ActionResponse, error) {
	res, err := h.giusecase.AddMemberToGroup(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) GrantGroupItemAccess(ctx context.Context, req *pb.GrantAccessRequest) (*pb.ActionResponse, error) {
	res, err := h.giusecase.GrantAccess(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) GetPendingRequests(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.FriendListResponse, error) {
	res, err := h.usecase.GetPendingRequests(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) AcceptFriend(ctx context.Context, req *pb.AcceptFriendRequest) (*pb.FriendResponse, error) {
	res, err := h.usecase.AcceptFriend(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) SetCloseFriend(ctx context.Context, req *pb.SetCloseFriendRequest) (*pb.SetCloseFriendResponse, error) {
	res, err := h.usecase.SetCloseFriend(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) GetCloseFriends(ctx context.Context, req *pb.GetCloseFriendsRequest) (*pb.GetCloseFriendsResponse, error) {
	res, err := h.usecase.GetCloseFriends(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) DeleteFriend(ctx context.Context, req *pb.FriendRequest) (*pb.FriendResponse, error) {
	res, err := h.usecase.DeleteFriend(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) RemoveMember(ctx context.Context, req *pb.RemoveMemberRequest) (*pb.ActionResponse, error) {
	res, err := h.gusecase.RemoveMember(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) LeaveGroup(ctx context.Context, req *pb.LeaveGroupRequest) (*pb.ActionResponse, error) {
	res, err := h.gusecase.LeaveGroup(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) GetGroupMessages(ctx context.Context, req *pb.GetGroupMessagesRequest) (*pb.GetGroupMessagesResponse, error) {
	res, err := h.musecase.GetGroupMessages(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) GetPrivateMessages(ctx context.Context, req *pb.GetPrivateMessagesRequest) (*pb.GetPrivateMessagesResponse, error) {
	res, err := h.musecase.GetPrivateMessages(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) GetItemSharedTargets(ctx context.Context, req *pb.GetItemSharedTargetsRequest) (*pb.GetItemSharedTargetsResponse, error) {
	res, err := h.giusecase.GetItemSharedTargets(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) GetSharedItemIDs(ctx context.Context, req *pb.GetSharedItemIDsRequest) (*pb.GetSharedItemIDsResponse, error) {
	res, err := h.giusecase.GetSharedItemIDs(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) GetItemsSharedByFriend(ctx context.Context, req *pb.GetItemsSharedByFriendRequest) (*pb.GetItemsSharedByFriendResponse, error) {
	res, err := h.giusecase.GetItemsSharedByFriend(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) DeleteGroup(ctx context.Context, req *pb.DeleteGroupRequest) (*pb.ActionResponse, error) {
	res, err := h.gusecase.DeleteGroup(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) GetAllSharedItemIDsByUser(ctx context.Context, req *pb.GetAllSharedItemIDsByUserRequest) (*pb.GetAllSharedItemIDsByUserResponse, error) {
	res, err := h.giusecase.GetAllSharedItemIDsByUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
