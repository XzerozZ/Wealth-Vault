package usecase

import (
	"context"
	"errors"
	"log"
	"time"
	"wealth-vault/user-service/internal/domain"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
	"wealth-vault/user-service/pkg/utils"

	"github.com/google/uuid"
)

func (u *UserUsecase) AddFriend(ctx context.Context, req *pb.FriendRequest) (*pb.FriendResponse, error) {
	userID, err := utils.ParseUUID(req.Id)
	if err != nil {
		return nil, err
	}

	friendID, err := utils.ParseUUID(req.UserId)
	if err != nil {
		return nil, err
	}

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

	reverseExists, reverseStatus, _ := u.userRepo.CheckFriendship(ctx, friendID, userID)
	if reverseExists {
		if reverseStatus == "ACCEPTED" {
			return nil, errors.New("already friends")
		}
		if reverseStatus == "PENDING" {
			return nil, errors.New("this user has already sent you a friend request, please check your pending requests")
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
	currentUserID, err := utils.ParseUUID(req.UserId)
	if err != nil {
		return nil, err
	}

	requesterID, err := utils.ParseUUID(req.RequesterId)
	if err != nil {
		return nil, err
	}

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

	if req.Action == "DECLINE" {
		err := u.userRepo.RemoveFriend(ctx, currentUserID, requesterID)
		go u.publisher.Publish("noti.friend.decline", evt)
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
	go u.publisher.Publish("noti.friend.accepted", evt)

	return &pb.FriendResponse{
		Success: true,
	}, nil
}

func (u *UserUsecase) SetCloseFriend(ctx context.Context, req *pb.SetCloseFriendRequest) (*pb.SetCloseFriendResponse, error) {
	userID, err := utils.ParseUUID(req.UserId)
	if err != nil {
		return nil, err
	}

	friendID, err := uuid.Parse(req.FriendId)
	if err != nil {
		return nil, err
	}

	exists, status, err := u.userRepo.CheckFriendship(ctx, userID, friendID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("friendship not found")
	}

	if status != "ACCEPTED" {
		return nil, errors.New("cannot set close friend while friendship status is pending")
	}

	if err := u.userRepo.SetCloseFriendStatus(ctx, userID, friendID, req.IsClose); err != nil {
		return nil, err
	}

	return &pb.SetCloseFriendResponse{
		Success: true,
	}, nil
}

func (u *UserUsecase) GetCloseFriends(ctx context.Context, req *pb.GetCloseFriendsRequest) (*pb.GetCloseFriendsResponse, error) {
	userID, err := utils.ParseUUID(req.UserId)
	if err != nil {
		return nil, err
	}

	closefriendlists, err := u.userRepo.GetCloseFriends(ctx, userID)
	if err != nil {
		return nil, err
	}

	var friends []*pb.User
	for _, item := range closefriendlists {
		friendProto := utils.ToUserProto(&item.Friend)
		friendProto.IsCloseFriend = item.IsCloseFriend
		friends = append(friends, friendProto)
	}

	log.Println(friends[0].IsCloseFriend)
	return &pb.GetCloseFriendsResponse{
		Friends: friends,
	}, nil
}

func (u *UserUsecase) GetFriendList(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.FriendListResponse, error) {
	id, err := utils.ParseUUID(req.Id)
	if err != nil {
		return nil, err
	}

	friendLists, err := u.userRepo.GetFriendList(ctx, id)
	if err != nil {
		return nil, err
	}

	var friends []*pb.User
	for _, item := range friendLists {
		friendProto := utils.ToUserProto(&item.Friend)
		friendProto.IsCloseFriend = item.IsCloseFriend
		friends = append(friends, friendProto)
	}

	return &pb.FriendListResponse{
		Friends: friends,
	}, nil
}

func (u *UserUsecase) GetPendingRequests(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.FriendListResponse, error) {
	id, err := utils.ParseUUID(req.Id)
	if err != nil {
		return nil, err
	}

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

func (u *UserUsecase) DeleteFriend(ctx context.Context, req *pb.FriendRequest) (*pb.FriendResponse, error) {
	userID, err := utils.ParseUUID(req.Id)
	if err != nil {
		return nil, err
	}

	friendID, err := utils.ParseUUID(req.UserId)
	if err != nil {
		return nil, err
	}

	exists, _, err := u.userRepo.CheckFriendship(ctx, userID, friendID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("friendship not found")
	}

	if err := u.userRepo.RemoveFriendAndSharedItems(ctx, userID, friendID); err != nil {
		return &pb.FriendResponse{Success: false}, err
	}

	return &pb.FriendResponse{
		Success: true,
	}, nil
}

func (u *UserUsecase) VerifyFriendship(ctx context.Context, req *pb.CheckFriendshipRequest) (*pb.CheckFriendshipResponse, error) {
	userID, err := utils.ParseUUID(req.UserId)
	if err != nil {
		return nil, err
	}

	friendID, err := utils.ParseUUID(req.TargetId)
	if err != nil {
		return nil, err
	}

	if userID == friendID {
		return &pb.CheckFriendshipResponse{
			IsFriend: true,
		}, nil
	}

	isRelated, status, err := u.userRepo.CheckFriendship(ctx, userID, friendID)
	if err != nil {
		return &pb.CheckFriendshipResponse{
			IsFriend: false,
		}, err
	}

	if isRelated && status == "ACCEPTED" {
		return &pb.CheckFriendshipResponse{
			IsFriend: true,
		}, nil
	}

	return &pb.CheckFriendshipResponse{
		IsFriend: false,
	}, nil
}
