package usecase

import (
	"context"
	"errors"
	"fmt"
	"wealth-vault/user-service/internal/domain"
	"wealth-vault/user-service/internal/event"
	storage "wealth-vault/user-service/internal/infra/storage"
	repo "wealth-vault/user-service/internal/repository/interface"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
	"wealth-vault/user-service/pkg/utils"
	"wealth-vault/user-service/pkg/utils/helper"
)

type GroupUsecase struct {
	groupRepo repo.GroupRepository
	userRepo  repo.UserRepository
	storage   storage.SupabaseStorage
	publisher event.EventPublisher
	msgRepo   repo.MsgRepository
}

func NewGroupUsecase(r repo.GroupRepository, u repo.UserRepository, m repo.MsgRepository, s storage.SupabaseStorage, e event.EventPublisher) *GroupUsecase {
	return &GroupUsecase{
		groupRepo: r,
		userRepo:  u,
		storage:   s,
		publisher: e,
		msgRepo:   m,
	}
}

func (u *GroupUsecase) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (*pb.GroupResponse, error) {
	creatorID, err := utils.ParseUUID(req.CreatorId)
	if err != nil {
		return nil, err
	}

	group := &domain.Group{
		GroupName:    req.Name,
		GroupProfile: req.Profile,
		CreatedBy:    creatorID,
	}

	if err := u.groupRepo.CreateGroup(ctx, group, req.MemberIds); err != nil {
		return nil, err
	}

	u.DispatchGroupCreated(ctx, group, req)

	proto := utils.ToGroupProto(group)
	proto.MemberCount = int64(len(req.MemberIds) + 1)

	return &pb.GroupResponse{
		Group: proto,
	}, nil
}

func (u *GroupUsecase) GetGroup(ctx context.Context, req *pb.GetGroupRequest) (*pb.GroupResponse, error) {
	groupID, err := utils.ParseUUID(req.GroupId)
	if err != nil {
		return nil, err
	}

	userID, err := utils.ParseUUID(req.UserId)
	if err != nil {
		return nil, err
	}

	if err := u.ensureMember(ctx, groupID, userID); err != nil {
		return nil, err
	}

	group, total, err := u.groupRepo.GetGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}

	proto := utils.ToGroupProto(group)
	proto.MemberCount = total

	return &pb.GroupResponse{
		Group: proto,
	}, nil
}

func (u *GroupUsecase) GetMember(ctx context.Context, req *pb.GetGroupMembersRequest) (*pb.GetGroupMembersResponse, error) {
	groupID, err := utils.ParseUUID(req.GroupId)
	if err != nil {
		return nil, err
	}

	userID, err := utils.ParseUUID(req.UserId)
	if err != nil {
		return nil, err
	}

	if err := u.ensureMember(ctx, groupID, userID); err != nil {
		return nil, err
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

func (u *GroupUsecase) AllGetGroup(ctx context.Context, req *pb.AllGroupRequest) (*pb.AllGroupResponse, error) {
	userID, err := utils.ParseUUID(req.UserId)
	if err != nil {
		return nil, err
	}

	groups, err := u.groupRepo.AllGetGroup(ctx, userID)
	if err != nil {
		return nil, err
	}

	var pbGroups []*pb.Group

	for _, item := range groups {
		pbGroup := utils.ToGroupProto(&item.Group)
		pbGroup.MemberCount = item.MemberCount
		pbGroups = append(pbGroups, pbGroup)
	}

	return &pb.AllGroupResponse{
		Group: pbGroups,
	}, nil
}

func (u *GroupUsecase) UpdateGroup(ctx context.Context, req *pb.UpdateGroupRequest) (*pb.GroupResponse, error) {
	groupID, err := utils.ParseUUID(req.Id)
	if err != nil {
		return nil, err
	}

	userID, err := utils.ParseUUID(req.UserId)
	if err != nil {
		return nil, err
	}

	if err := u.ensureMember(ctx, groupID, userID); err != nil {
		return nil, err
	}

	group, _, err := u.groupRepo.GetGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}

	editorName := u.getUsernameSafe(ctx, userID)

	updatedMask, err := helper.ApplyUpdateGroupFields(req, u.storage, group)
	if err != nil {
		return nil, err
	}

	logEntry := u.buildSystemLog(
		groupID,
		userID,
		fmt.Sprintf("%s ได้แก้ไขข้อมูลกลุ่ม", editorName),
		nil,
	)

	updatedGroup, count, err := u.groupRepo.UpdateGroup(
		ctx,
		group,
		updatedMask,
		logEntry,
	)
	if err != nil {
		return nil, err
	}

	proto := utils.ToGroupProto(updatedGroup)
	proto.MemberCount = count

	return &pb.GroupResponse{
		Group: proto,
	}, nil
}

func (u *GroupUsecase) RemoveMember(ctx context.Context, req *pb.RemoveMemberRequest) (*pb.ActionResponse, error) {
	groupID, err := utils.ParseUUID(req.GroupId)
	if err != nil {
		return nil, err
	}

	adminID, err := utils.ParseUUID(req.UserId)
	if err != nil {
		return nil, err
	}

	targetID, err := utils.ParseUUID(req.TargetMemberId)
	if err != nil {
		return nil, err
	}

	group, _, err := u.groupRepo.GetGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}

	if group.CreatedBy != adminID {
		return nil, errors.New("only admin can remove members")
	}

	adminName := u.getUsernameSafe(ctx, adminID)
	targetName := u.getUsernameSafe(ctx, targetID)

	logEntry := u.buildSystemLog(
		groupID,
		adminID,
		fmt.Sprintf("%s ได้ลบ %s ออกจากกลุ่ม", adminName, targetName),
		map[string]interface{}{
			"action":    "remove_member",
			"target_id": targetID.String(),
		},
	)

	if err := u.groupRepo.RemoveMemberAndTheirSharedItems(ctx, groupID, targetID, logEntry); err != nil {
		return nil, err
	}

	u.DispatchMemberRemoved(group, req, adminName)

	return &pb.ActionResponse{
		Success: true,
	}, nil
}

func (u *GroupUsecase) LeaveGroup(ctx context.Context, req *pb.LeaveGroupRequest) (*pb.ActionResponse, error) {
	groupID, err := utils.ParseUUID(req.GroupId)
	if err != nil {
		return nil, err
	}

	userID, err := utils.ParseUUID(req.UserId)
	if err != nil {
		return nil, err
	}

	group, _, err := u.groupRepo.GetGroup(ctx, groupID)
	if err != nil {
		return nil, errors.New("group not found")
	}

	if group.CreatedBy == userID {
		return nil, errors.New("owner cannot leave group")
	}

	if err := u.ensureMember(ctx, groupID, userID); err != nil {
		return nil, err
	}

	userName := u.getUsernameSafe(ctx, userID)

	logEntry := u.buildSystemLog(
		groupID,
		userID,
		fmt.Sprintf("%s ได้ออกจากกลุ่มแล้ว", userName),
		map[string]interface{}{
			"action": "leave_group",
		},
	)

	if err := u.groupRepo.RemoveMemberAndTheirSharedItems(ctx, groupID, userID, logEntry); err != nil {
		return nil, err
	}

	u.DispatchMemberLeft(group, req, userName)

	return &pb.ActionResponse{
		Success: true,
	}, nil
}

func (u *GroupUsecase) DeleteGroup(ctx context.Context, req *pb.DeleteGroupRequest) (*pb.ActionResponse, error) {
	groupID, err := utils.ParseUUID(req.GroupId)
	if err != nil {
		return nil, err
	}

	userID, err := utils.ParseUUID(req.UserId)
	if err != nil {
		return nil, err
	}

	group, _, err := u.groupRepo.GetGroup(ctx, groupID)
	if err != nil {
		return nil, errors.New("group not found")
	}

	if group.CreatedBy != userID {
		return nil, errors.New("only creator can delete group")
	}

	if err := u.groupRepo.DeleteGroupCompletely(ctx, groupID); err != nil {
		return nil, err
	}

	return &pb.ActionResponse{
		Success: true,
	}, nil
}
