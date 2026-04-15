package usecase

import (
	"context"
	"wealth-vault/user-service/internal/domain"
	"wealth-vault/user-service/internal/event"
	storage "wealth-vault/user-service/internal/infra/storage"
	repo "wealth-vault/user-service/internal/repository/interface"
	usecase "wealth-vault/user-service/internal/usecase/interface"
	assetPb "wealth-vault/user-service/pkg/pb/proto/asset"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
	"wealth-vault/user-service/pkg/utils"
	helper "wealth-vault/user-service/pkg/utils/helper"
)

type UserUsecase struct {
	userRepo    repo.UserRepository
	itemUC      usecase.ShareItemUsecase
	storage     storage.SupabaseStorage
	publisher   event.EventPublisher
	assetClient assetPb.AssetServiceClient
}

func NewUserUsecase(r repo.UserRepository, i usecase.ShareItemUsecase, s storage.SupabaseStorage, e event.EventPublisher, assetClient assetPb.AssetServiceClient) *UserUsecase {
	return &UserUsecase{
		userRepo:    r,
		itemUC:      i,
		storage:     s,
		publisher:   e,
		assetClient: assetClient,
	}
}

func (u *UserUsecase) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	user := &domain.User{
		Email:    req.Email,
		Username: req.Username,
	}

	if err := u.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return &pb.UserResponse{
		Success: true,
		User:    utils.ToUserProto(user),
	}, nil
}

func (u *UserUsecase) GetUser(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.UserResponse, error) {
	id, err := utils.ParseUUID(req.Id)
	if err != nil {
		return nil, err
	}

	user, err := u.userRepo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{
		Success: true,
		User:    utils.ToUserProto(user),
	}, nil
}

func (u *UserUsecase) GetUsersByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.UserInfoResponse, error) {
	users, err := u.userRepo.GetUsersByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	currentUserID, _ := utils.ParseUUID(req.Id)
	var userProtos []*pb.User
	for _, user := range users {
		proto := utils.ToUserProto(user)

		if user.ID == currentUserID {
			proto.IsFriend = false
		} else {
			isFriend, status, _ := u.userRepo.CheckFriendship(ctx, currentUserID, user.ID)
			proto.IsFriend = isFriend && status == "ACCEPTED"
		}

		userProtos = append(userProtos, proto)
	}

	return &pb.UserInfoResponse{
		Success: true,
		User:    userProtos,
	}, nil
}

func (u *UserUsecase) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	id, err := utils.ParseUUID(req.Id)
	if err != nil {
		return nil, err
	}

	user, err := u.userRepo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	updateMask, err := helper.ApplyUpdateUserFields(req, u.storage, user)
	if err != nil {
		return nil, err
	}

	updatedUser, err := u.userRepo.UpdateUser(ctx, user, updateMask)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{
		Success: true,
		User:    utils.ToUserProto(updatedUser),
	}, nil
}
