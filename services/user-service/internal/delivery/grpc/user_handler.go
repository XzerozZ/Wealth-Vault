package grpc

import (
	"context"
	"wealth-vault/user-service/internal/usecase"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
)

type UserGRPCHandler struct {
	pb.UnimplementedUserServiceServer
	usecase  usecase.UserUsecase
	gusecase usecase.GroupUsecase
}

func NewUserGRPCHandler(u usecase.UserUsecase, g usecase.GroupUsecase) *UserGRPCHandler {
	return &UserGRPCHandler{
		usecase:  u,
		gusecase: g,
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
