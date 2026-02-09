package usecase

import (
	"context"
	"errors"
	"log"
	"time"
	"wealth-vault/user-service/internal/domain"
	"wealth-vault/user-service/internal/event"
	repo "wealth-vault/user-service/internal/repository/interface"
	assetPb "wealth-vault/user-service/pkg/pb/proto/asset"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
	"wealth-vault/user-service/pkg/utils"
	helper "wealth-vault/user-service/pkg/utils/helper"

	"github.com/google/uuid"
)

type UserUsecase struct {
	userRepo    repo.UserRepository
	itemUC      ShareItemUsecase
	storage     *utils.StorageClient
	publisher   *event.Publisher
	assetClient assetPb.AssetServiceClient
}

func NewUserUsecase(r repo.UserRepository, i ShareItemUsecase, s *utils.StorageClient, e *event.Publisher, assetClient assetPb.AssetServiceClient) UserUsecase {
	return UserUsecase{
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
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, errors.New("invalid user id")
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

func (u *UserUsecase) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, errors.New("invalid user id")
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

func (u *UserUsecase) GetFriendList(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.FriendListResponse, error) {
	id, _ := uuid.Parse(req.Id)
	friendLists, err := u.userRepo.GetFriendList(ctx, id)
	if err != nil {
		return nil, err
	}

	var friends []*pb.User
	for _, item := range friendLists {
		friends = append(friends, utils.ToUserProto(&item.Friend))
	}

	return &pb.FriendListResponse{
		Friends: friends,
	}, nil
}

func (u *UserUsecase) GetPendingRequests(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.FriendListResponse, error) {
	id, _ := uuid.Parse(req.Id)

	requests, err := u.userRepo.GetIncomingRequests(ctx, id)
	if err != nil {
		return nil, err
	}

	var pendingUsers []*pb.User
	for _, item := range requests {
		pendingUsers = append(pendingUsers, utils.ToUserProto(&item.User))
	}

	return &pb.FriendListResponse{
		Friends: pendingUsers,
	}, nil
}

func (u *UserUsecase) AddFriend(ctx context.Context, req *pb.FriendRequest) (*pb.FriendResponse, error) {
	userID, _ := uuid.Parse(req.Id)
	friendID, _ := uuid.Parse(req.UserId)

	if userID == friendID {
		return nil, errors.New("cannot add yourself")
	}

	exists, status, _ := u.userRepo.CheckFriendship(ctx, userID, friendID)
	if exists {
		if status == "ACCEPTED" {
			return nil, errors.New("already friends")
		}
		if status == "PENDING" {
			return nil, errors.New("friend request already sent")
		}
	}

	friendRequest := &domain.FriendList{
		UserID:   userID,
		FriendID: friendID,
		Status:   "PENDING",
	}

	if err := u.userRepo.AddFriend(ctx, friendRequest); err != nil {
		return nil, err
	}

	me, _ := u.userRepo.GetUser(ctx, userID)

	evt := domain.FriendRequestEvent{
		SenderID:   req.Id,
		SenderName: me.Username,
		TargetID:   req.UserId,
		OccurredAt: time.Now().Unix(),
	}

	go u.publisher.Publish("noti.friend.request", evt)

	return &pb.FriendResponse{
		Success: true,
	}, nil
}

func (u *UserUsecase) AcceptFriend(ctx context.Context, req *pb.AcceptFriendRequest) (*pb.FriendResponse, error) {
	currentUserID, _ := uuid.Parse(req.UserId)
	requesterID, _ := uuid.Parse(req.RequesterId)
	if req.Action == "DECLINE" {
		err := u.userRepo.RemoveFriend(ctx, currentUserID, requesterID)
		if err != nil {
			return &pb.FriendResponse{
				Success: false,
			}, err
		}
		return &pb.FriendResponse{
			Success: true,
		}, nil
	}

	if err := u.userRepo.UpdateFriendStatus(ctx, currentUserID, requesterID, "ACCEPTED"); err != nil {
		return nil, err
	}

	reverseFriend := &domain.FriendList{
		UserID:   currentUserID,
		FriendID: requesterID,
		Status:   "ACCEPTED",
	}

	_ = u.userRepo.CreateFriendship(ctx, reverseFriend)
	accepter, err := u.userRepo.GetUser(ctx, currentUserID)
	accepterName := "Unknown"
	if err == nil {
		accepterName = accepter.Username
	}

	evt := domain.FriendAcceptedEvent{
		AccepterID:   req.UserId,
		AccepterName: accepterName,
		RequesterID:  req.RequesterId,
		OccurredAt:   time.Now().Unix(),
	}

	go u.publisher.Publish("noti.friend.accepted", evt)

	return &pb.FriendResponse{
		Success: true,
	}, nil
}

func (u *UserUsecase) ProcessLegacyAutoShare(ctx context.Context) error {
	log.Println("⏰ [LEGACY] Start checking for eligible users...")

	users, err := u.userRepo.GetUsersReadyForAutoShare(ctx)
	if err != nil {
		return nil
	}

	log.Printf("🔎 Found %d users eligible for legacy transfer", len(users))
	for _, user := range users {
		u.ProcessSingleUserLegacy(ctx, user)
	}

	return nil
}

func (u *UserUsecase) ProcessSingleUserLegacy(ctx context.Context, user domain.User) error {
	if len(user.Friends) == 0 {
		return nil
	}

	assets, err := u.assetClient.GetAllAssetIDs(ctx, &assetPb.GetMyAssetsRequest{
		UserId: user.ID.String(),
	})
	if err != nil {
		return err
	}

	totalAssets := len(assets.AccountIds) + len(assets.BuildingIds) + len(assets.LandIds) +
		len(assets.CashIds) + len(assets.InsuranceIds) + len(assets.InvestmentIds) + len(assets.LiabilityIds)

	if totalAssets == 0 {
		u.userRepo.MarkAutoShareTriggered(ctx, user.ID)
		return nil
	}

	for _, friend := range user.Friends {
		err := u.itemUC.BatchShareAssets(ctx, domain.BatchShareRequest{
			OwnerID:       user.ID,
			TargetID:      friend.ID,
			AccountIDs:    assets.AccountIds,
			BuildingIDs:   assets.BuildingIds,
			LandIDs:       assets.LandIds,
			CashIDs:       assets.CashIds,
			InsuranceIDs:  assets.InsuranceIds,
			InvestmentIDs: assets.InvestmentIds,
			LiabilityIDs:  assets.LiabilityIds,
		})

		if err != nil {
			log.Printf("⚠️ Failed to share to %s: %v", friend.ID, err)
		}
	}

	u.userRepo.MarkAutoShareTriggered(ctx, user.ID)
	return nil
}

func (u *UserUsecase) SetCloseFriend(ctx context.Context, req *pb.SetCloseFriendRequest) (*pb.SetCloseFriendResponse, error) {
	userID, _ := uuid.Parse(req.UserId)
	friendID, _ := uuid.Parse(req.FriendId)
	if err := u.userRepo.SetCloseFriendStatus(ctx, userID, friendID, req.IsClose); err != nil {
		return nil, err
	}

	return &pb.SetCloseFriendResponse{
		Success: true,
	}, nil
}

func (u *UserUsecase) GetCloseFriends(ctx context.Context, req *pb.GetCloseFriendsRequest) (*pb.GetCloseFriendsResponse, error) {
	userID, _ := uuid.Parse(req.UserId)
	closefriendlists, err := u.userRepo.GetCloseFriends(ctx, userID)
	if err != nil {
		return nil, err
	}

	var friends []*pb.User
	for _, item := range closefriendlists {
		friends = append(friends, utils.ToUserProto(&item))
	}

	return &pb.GetCloseFriendsResponse{
		Friends: friends,
	}, nil
}
