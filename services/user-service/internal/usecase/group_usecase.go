package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"wealth-vault/user-service/internal/domain"
	"wealth-vault/user-service/internal/event"
	repo "wealth-vault/user-service/internal/repository/interface"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
	"wealth-vault/user-service/pkg/utils"
	"wealth-vault/user-service/pkg/utils/helper"

	"github.com/google/uuid"
)

type GroupUsecase struct {
	groupRepo repo.GroupRepository
	userRepo  repo.UserRepository
	storage   *utils.StorageClient
	publisher *event.Publisher
	msgRepo   repo.MsgRepository
}

func NewGroupUsecase(r repo.GroupRepository, u repo.UserRepository, m repo.MsgRepository, s *utils.StorageClient, e *event.Publisher) GroupUsecase {
	return GroupUsecase{
		groupRepo: r,
		userRepo:  u,
		storage:   s,
		publisher: e,
		msgRepo:   m,
	}
}

func (u *GroupUsecase) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (*pb.GroupResponse, error) {
	userID, _ := uuid.Parse(req.CreatorId)
	group := &domain.Group{
		GroupName:    req.Name,
		GroupProfile: req.Profile,
		CreatedBy:    userID,
	}

	if err := u.groupRepo.CreateGroup(ctx, group, req.MemberIds); err != nil {
		return nil, err
	}

	var notifyTargetIDs []string
	for _, id := range req.MemberIds {
		if id != req.CreatorId {
			notifyTargetIDs = append(notifyTargetIDs, id)
		}
	}

	if len(notifyTargetIDs) > 0 {
		senderName := "Unknown"
		if creator, err := u.userRepo.GetUser(ctx, userID); err == nil {
			senderName = creator.Username
		}

		evt := domain.GroupCreatedEvent{
			GroupID:       group.ID.String(),
			GroupName:     group.GroupName,
			SenderID:      req.CreatorId,
			SenderName:    senderName,
			TargetUserIDs: notifyTargetIDs,
			OccurredAt:    time.Now().Unix(),
		}

		go u.publisher.Publish("noti.group.created", evt)
	}

	protoGroup := utils.ToGroupProto(group)
	protoGroup.MemberCount = int64(len(req.MemberIds) + 1)

	return &pb.GroupResponse{
		Group: protoGroup,
	}, nil
}

func (u *GroupUsecase) GetMember(ctx context.Context, req *pb.GetGroupMembersRequest) (*pb.GetGroupMembersResponse, error) {
	groupID, _ := uuid.Parse(req.GroupId)
	userID, _ := uuid.Parse(req.UserId)
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
	groupID, _ := uuid.Parse(req.GroupId)
	userID, _ := uuid.Parse(req.UserId)
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

func (u *GroupUsecase) AllGetGroup(ctx context.Context, req *pb.AllGroupRequest) (*pb.AllGroupResponse, error) {
	userID, _ := uuid.Parse(req.UserId)
	group, err := u.groupRepo.AllGetGroup(ctx, userID)
	if err != nil {
		return nil, err
	}

	var pbGroups []*pb.Group

	for _, item := range group {
		pbGroup := utils.ToGroupProto(&item.Group)
		pbGroup.MemberCount = item.MemberCount
		pbGroups = append(pbGroups, pbGroup)
	}

	return &pb.AllGroupResponse{
		Group: pbGroups,
	}, nil
}

func (u *GroupUsecase) UpdateGroup(ctx context.Context, req *pb.UpdateGroupRequest) (*pb.GroupResponse, error) {
	groupID, _ := uuid.Parse(req.Id)
	userID, _ := uuid.Parse(req.UserId)
	isMember, err := u.groupRepo.IsUserMember(ctx, groupID, userID)
	if err != nil {
		return nil, errors.New("failed to check membership")
	}

	if !isMember {
		return nil, errors.New("access denied: you are not a member of this group")
	}

	editor, _ := u.userRepo.GetUser(ctx, userID)
	group, _, err := u.groupRepo.GetGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}

	updatedMask, err := helper.ApplyUpdateGroupFields(req, u.storage, group)
	if err != nil {
		return nil, err
	}

	logEntry := &domain.GroupLog{
		GroupID:   groupID,
		LogType:   "SYSTEM",
		Messages:  fmt.Sprintf("%s ได้แก้ไขข้อมูลกลุ่ม", editor.Username),
		CreatedBy: userID,
	}

	updatedGroup, count, err := u.groupRepo.UpdateGroup(ctx, group, updatedMask, logEntry)
	if err != nil {
		return nil, err
	}

	protoGroup := utils.ToGroupProto(updatedGroup)
	protoGroup.MemberCount = count

	return &pb.GroupResponse{
		Group: protoGroup,
	}, nil
}

func (u *GroupUsecase) RemoveMember(ctx context.Context, req *pb.RemoveMemberRequest) (*pb.ActionResponse, error) {
	groupID, _ := uuid.Parse(req.GroupId)
	userID, _ := uuid.Parse(req.UserId)
	targetID, _ := uuid.Parse(req.TargetMemberId)

	group, _, err := u.groupRepo.GetGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}

	if group.CreatedBy != userID {
		return nil, errors.New("unauthorized: only admin can remove members")
	}

	adminUser, _ := u.userRepo.GetUser(ctx, userID)
	targetUser, _ := u.userRepo.GetUser(ctx, targetID)
	adminName := "Admin"
	targetName := "Member"
	if adminUser != nil {
		adminName = adminUser.Username
	}
	if targetUser != nil {
		targetName = targetUser.Username
	}

	metaMap := map[string]interface{}{
		"action":    "remove_member",
		"target_id": targetID.String(),
		"reason":    "admin_removed",
	}
	metaJSON, _ := json.Marshal(metaMap)

	logEntry := &domain.GroupLog{
		GroupID:   groupID,
		LogType:   "SYSTEM",
		Messages:  fmt.Sprintf("%s ได้ลบ %s ออกจากกลุ่ม", adminName, targetName),
		Metadata:  string(metaJSON),
		CreatedBy: userID,
	}

	if err := u.groupRepo.RemoveMemberAndTheirSharedItems(ctx, groupID, targetID, logEntry); err != nil {
		return nil, err
	}

	go func() {
		bgCtx := context.Background()
		sysMsg := domain.GroupMessage{
			GroupID:   groupID,
			SenderID:  userID,
			MsgType:   "SYSTEM_ALERT",
			Content:   fmt.Sprintf("%s ลบ %s ออกจากกลุ่ม", adminName, targetName),
			Metadata:  "{}",
			CreatedAt: time.Now(),
		}
		u.msgRepo.CreateMessage(bgCtx, []domain.GroupMessage{sysMsg})

		broadcastEvt := domain.GroupActivityEvent{
			GroupID:      req.GroupId,
			ActivityType: "MEMBER_REMOVED",
			Payload:      fmt.Sprintf("%s ลบ %s ออกจากกลุ่ม", adminName, targetName),
			ActorID:      req.UserId,
			TargetID:     req.TargetMemberId,
			OccurredAt:   time.Now().Unix(),
		}
		u.publisher.Publish("noti.group.activity", broadcastEvt)

		pushEvt := domain.MemberRemovedEvent{
			GroupID:    req.GroupId,
			GroupName:  group.GroupName,
			TargetID:   req.TargetMemberId,
			ActionBy:   adminName,
			OccurredAt: time.Now().Unix(),
		}
		u.publisher.Publish("noti.group.member.removed", pushEvt)
	}()

	return &pb.ActionResponse{Success: true}, nil
}

func (u *GroupUsecase) LeaveGroup(ctx context.Context, req *pb.LeaveGroupRequest) (*pb.ActionResponse, error) {
	groupID, _ := uuid.Parse(req.GroupId)
	userID, _ := uuid.Parse(req.UserId)

	group, _, err := u.groupRepo.GetGroup(ctx, groupID)
	if err != nil {
		return nil, errors.New("group not found")
	}

	if group.CreatedBy == userID {
		return nil, errors.New("owner cannot leave the group. please delete the group instead")
	}

	isMember, err := u.groupRepo.IsUserMember(ctx, groupID, userID)
	if err != nil || !isMember {
		return nil, errors.New("you are not a member of this group")
	}

	user, _ := u.userRepo.GetUser(ctx, userID)
	userName := "Member"
	if user != nil {
		userName = user.Username
	}

	metaMap := map[string]interface{}{
		"action": "leave_group",
		"reason": "user_left",
	}
	metaJSON, _ := json.Marshal(metaMap)

	logEntry := &domain.GroupLog{
		GroupID:   groupID,
		LogType:   "SYSTEM",
		Messages:  fmt.Sprintf("%s ได้ออกจากกลุ่มแล้ว", userName),
		Metadata:  string(metaJSON),
		CreatedBy: userID,
	}

	if err := u.groupRepo.RemoveMemberAndTheirSharedItems(ctx, groupID, userID, logEntry); err != nil {
		return nil, err
	}

	go func() {
		bgCtx := context.Background()

		sysMsg := domain.GroupMessage{
			GroupID:   groupID,
			SenderID:  userID,
			MsgType:   "SYSTEM_ALERT",
			Content:   fmt.Sprintf("%s ออกจากกลุ่ม", userName),
			Metadata:  "{}",
			CreatedAt: time.Now(),
		}
		u.msgRepo.CreateMessage(bgCtx, []domain.GroupMessage{sysMsg})

		broadcastEvt := domain.GroupActivityEvent{
			GroupID:      req.GroupId,
			ActivityType: "MEMBER_LEFT",
			Payload:      fmt.Sprintf("%s ออกจากกลุ่ม", userName),
			ActorID:      req.UserId,
			TargetID:     req.UserId,
			OccurredAt:   time.Now().Unix(),
		}
		u.publisher.Publish("noti.group.activity", broadcastEvt)
	}()

	return &pb.ActionResponse{Success: true}, nil
}

func (u *GroupUsecase) DeleteGroup(ctx context.Context, req *pb.DeleteGroupRequest) (*pb.ActionResponse, error) {
	groupID, _ := uuid.Parse(req.GroupId)
	userID, _ := uuid.Parse(req.UserId)

	group, _, err := u.groupRepo.GetGroup(ctx, groupID)
	if err != nil {
		return nil, errors.New("group not found")
	}

	if group.CreatedBy != userID {
		return nil, errors.New("permission denied: only creator can delete group")
	}

	if err := u.groupRepo.DeleteGroup(ctx, groupID); err != nil {
		return nil, err
	}

	return &pb.ActionResponse{
		Success: true,
	}, nil
}
