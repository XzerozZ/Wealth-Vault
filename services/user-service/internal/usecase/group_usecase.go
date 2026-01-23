package usecase

import (
	"context"
	"errors"
	"wealth-vault/user-service/internal/domain"
	repo "wealth-vault/user-service/internal/repository/interface"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
	"wealth-vault/user-service/pkg/utils"
	"wealth-vault/user-service/pkg/utils/helper"

	"github.com/google/uuid"
)

type GroupUsecase struct {
	groupRepo repo.GroupRepository
	storage   *utils.StorageClient
}

func NewGroupUsecase(r repo.GroupRepository, s *utils.StorageClient) GroupUsecase {
	return GroupUsecase{
		groupRepo: r,
		storage:   s,
	}
}

func (u *GroupUsecase) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (*pb.GroupResponse, error) {
	userID, err := uuid.Parse(req.CreatorId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	group := &domain.Group{
		GroupName:    req.Name,
		GroupProfile: req.Profile,
		CreatedBy:    userID,
	}

	if err := u.groupRepo.CreateGroup(ctx, group, req.MemberIds); err != nil {
		return nil, err
	}

	protoGroup := utils.ToGroupProto(group)
	protoGroup.MemberCount = int64(len(req.MemberIds) + 1)

	return &pb.GroupResponse{
		Group: protoGroup,
	}, nil
}

func (u *GroupUsecase) GetMember(ctx context.Context, req *pb.GetGroupMembersRequest) (*pb.GetGroupMembersResponse, error) {
	groupID, err := uuid.Parse(req.GroupId)
	if err != nil {
		return nil, errors.New("invalid group id")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	isMember, err := u.groupRepo.IsUserMember(ctx, groupID, userID)
	if err != nil {
		return nil, errors.New("failed to check membership")
	}

	if !isMember {
		return nil, errors.New("access denied: you are not a member of this group")
	}

	users, total, err := u.groupRepo.GetMember(ctx, groupID)
	if err != nil {
		return nil, errors.New("failed to get group members")
	}

	var protoMembers []*pb.User
	for _, user := range users {
		protoMembers = append(protoMembers, utils.ToUserProto(user))
	}

	return &pb.GetGroupMembersResponse{
		Members: protoMembers,
		Total:   total,
	}, nil
}

func (u *GroupUsecase) GetGroup(ctx context.Context, req *pb.GetGroupRequest) (*pb.GroupResponse, error) {
	groupID, err := uuid.Parse(req.GroupId)
	if err != nil {
		return nil, errors.New("invalid group id")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	isMember, err := u.groupRepo.IsUserMember(ctx, groupID, userID)
	if err != nil {
		return nil, errors.New("failed to check membership")
	}

	if !isMember {
		return nil, errors.New("access denied: you are not a member of this group")
	}

	group, total, err := u.groupRepo.GetGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}

	protoGroup := utils.ToGroupProto(group)
	protoGroup.MemberCount = total

	return &pb.GroupResponse{
		Group: protoGroup,
	}, nil
}

func (u *GroupUsecase) UpdateGroup(ctx context.Context, req *pb.UpdateGroupRequest) (*pb.GroupResponse, error) {
	groupID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, errors.New("invalid group id")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	isMember, err := u.groupRepo.IsUserMember(ctx, groupID, userID)
	if err != nil {
		return nil, errors.New("failed to check membership")
	}

	if !isMember {
		return nil, errors.New("access denied: you are not a member of this group")
	}

	group, _, err := u.groupRepo.GetGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}

	updatedMask, err := helper.ApplyUpdateGroupFields(req, u.storage, group)
	if err != nil {
		return nil, err
	}

	updatedGroup, count, err := u.groupRepo.UpdateGroup(ctx, group, updatedMask)
	if err != nil {
		return nil, err
	}

	protoGroup := utils.ToGroupProto(updatedGroup)
	protoGroup.MemberCount = count

	return &pb.GroupResponse{
		Group: protoGroup,
	}, nil
}
