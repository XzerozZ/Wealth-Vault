package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
	"wealth-vault/user-service/internal/domain"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
	"wealth-vault/user-service/pkg/utils"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (u *ShareItemUsecase) GetSharedIteminGroup(ctx context.Context, req *pb.GetGroupItemsRequest) (*pb.GetGroupItemsResponse, error) {
	groupID, err := utils.ParseUUID(req.GroupId)
	if err != nil {
		return nil, err
	}

	uid, err := utils.ParseUUID(req.UserId)
	if err != nil {
		return nil, err
	}

	items, err := u.itemRepo.GetSharedIteminGroup(ctx, groupID, uid)
	if err != nil {
		return nil, err
	}

	var summaries []domain.SharedItemSummary
	for _, item := range items {
		summaries = append(summaries, domain.SharedItemSummary{
			EntityID:   item.EntityID.String(),
			EntityType: item.EntityType,
		})
	}

	previewMap, err := u.FetchAssetPreviews(ctx, summaries)
	if err != nil {
		log.Printf("Failed to fetch asset previews: %v", err)
	}

	var responseItems []*pb.GroupItemDetail
	for _, item := range items {
		responseItems = append(responseItems, &pb.GroupItemDetail{
			GroupItemId: item.ID.String(),
			SharedBy:    item.OwnerID.String(),
			SharedAt:    timestamppb.New(item.ShareAt),
			Type:        item.EntityType,
			AssetDetail: previewMap[item.EntityID.String()],
		})
	}

	return &pb.GetGroupItemsResponse{Items: responseItems}, nil
}

func (u *ShareItemUsecase) UnsharedIteminGroup(ctx context.Context, req *pb.UnshareItemRequest) (*pb.ShareItemResponse, error) {
	itemID, err := utils.ParseUUID(req.ItemId)
	if err != nil {
		return nil, err
	}

	item, err := u.itemRepo.GetGroupItemByID(ctx, itemID)
	if err != nil {
		return nil, err
	}

	userID, err := utils.ParseUUID(req.UserId)
	if err != nil {
		return nil, err
	}

	if err := u.itemRepo.DeleteIteminGroup(ctx, itemID, userID); err != nil {
		return nil, err
	}

	go u.msgRepo.MarkAssetMessageAsDeletedinAssetService(context.Background(), item.EntityID)

	return &pb.ShareItemResponse{Finish: true}, nil
}

func (u *ShareItemUsecase) AddMemberToGroup(ctx context.Context, req *pb.AddMemberRequest) (*pb.ActionResponse, error) {
	groupID, err := utils.ParseUUID(req.GroupId)
	if err != nil {
		return nil, err
	}

	senderID, err := utils.ParseUUID(req.UserId)
	if err != nil {
		return nil, err
	}

	if len(req.TargetUserIds) == 0 {
		return nil, errors.New("no users specified")
	}

	senderName := "Unknown"
	if sender, err := u.userRepo.GetUser(ctx, senderID); err == nil {
		senderName = sender.Username
	}

	var newMembers []domain.GroupMember
	var targetUUIDs []uuid.UUID
	var addedNames []string

	for _, userIDStr := range req.TargetUserIds {
		if uid, err := uuid.Parse(userIDStr); err == nil {
			userName := "Someone"
			if user, err := u.userRepo.GetUser(ctx, uid); err == nil {
				userName = user.Username
			}
			addedNames = append(addedNames, userName)
			targetUUIDs = append(targetUUIDs, uid)
			newMembers = append(newMembers, domain.GroupMember{
				GroupID:  groupID,
				UserID:   uid,
				Role:     "member",
				JoinedAt: time.Now(),
			})
		}
	}

	if len(newMembers) > 0 {
		if err := u.itemRepo.AddMember(ctx, newMembers); err != nil {
			return nil, fmt.Errorf("failed to add members: %v", err)
		}
	}

	go func() {
		bgCtx := context.Background()
		logMetaJSON, _ := json.Marshal(map[string]interface{}{"action": "add_member", "target_ids": req.TargetUserIds, "added_count": len(newMembers)})
		u.groupRepo.CreateLog(bgCtx, &domain.GroupLog{
			GroupID: groupID, LogType: LogTypeSystem,
			Messages: fmt.Sprintf("%s เพิ่มสมาชิกใหม่ %d คน", senderName, len(newMembers)), Metadata: string(logMetaJSON), CreatedBy: senderID,
		})

		msgContent := fmt.Sprintf("%s เพิ่ม %s เข้ากลุ่ม", senderName, strings.Join(addedNames, ", "))
		u.msgRepo.CreateMessage(bgCtx, []domain.GroupMessage{{
			GroupID: groupID, SenderID: senderID, MsgType: MsgTypeSystemAlert, Content: msgContent, CreatedAt: time.Now(),
		}})

		u.publisher.Publish(TopicGroupActivity, domain.GroupActivityEvent{
			GroupID: req.GroupId, ActivityType: "MEMBER_ADDED", Payload: fmt.Sprintf("%s เพิ่มสมาชิกใหม่", senderName), ActorID: req.UserId, OccurredAt: time.Now().Unix(),
		})

		existingMembers, _, _ := u.groupRepo.GetMember(bgCtx, groupID)
		newNamesStr := strings.Join(addedNames, ", ")
		for _, m := range existingMembers {
			isNew := false
			for _, tid := range targetUUIDs {
				if m.ID == tid {
					isNew = true
					break
				}
			}
			if isNew {
				continue
			}

			promptMeta, _ := json.Marshal(map[string]interface{}{
				"is_action_required": true,
				"is_completed":       false,
				"target_user_ids":    req.TargetUserIds,
				"type":               "GRANT_ACCESS_PROMPT",
			})

			u.msgRepo.CreateMessage(bgCtx, []domain.GroupMessage{{
				GroupID:   groupID,
				SenderID:  m.ID,
				MsgType:   MsgTypeGrantAccess,
				Content:   fmt.Sprintf("คุณต้องการแชร์รายการของคุณให้ %s หรือไม่?", newNamesStr),
				Metadata:  string(promptMeta),
				CreatedAt: time.Now(),
			}})
		}
	}()

	return &pb.ActionResponse{Success: true}, nil
}

func (u *ShareItemUsecase) GrantAccess(ctx context.Context, req *pb.GrantAccessRequest) (*pb.ActionResponse, error) {
	ownerID, err := utils.ParseUUID(req.OwnerUserId)
	if err != nil {
		return nil, err
	}

	targetID, err := utils.ParseUUID(req.TargetUserId)
	if err != nil {
		return nil, err
	}

	groupID, err := utils.ParseUUID(req.GroupId)
	if err != nil {
		return nil, err
	}

	if len(req.GroupItemIds) == 0 {
		go u.updateMessageToCompleted(groupID, ownerID, targetID)
		return &pb.ActionResponse{Success: true}, nil
	}

	if isMember, err := u.itemRepo.IsGroupMember(ctx, groupID, targetID); err != nil || !isMember {
		return nil, errors.New("target user is not a member of this group")
	}

	validItemIDs, err := u.itemRepo.GetOwnedItemIDs(ctx, req.GroupItemIds, ownerID)
	if err != nil || len(validItemIDs) == 0 {
		return nil, errors.New("permission denied or no valid items owned")
	}

	var viewers []domain.GroupItemViewer
	for _, id := range validItemIDs {
		viewers = append(viewers, domain.GroupItemViewer{
			GroupItemID: id,
			ViewerID:    targetID,
			GrantedAt:   time.Now(),
		})
	}
	if err := u.itemRepo.BatchCreateViewers(ctx, viewers); err != nil {
		return nil, err
	}

	go u.updateMessageToCompleted(groupID, ownerID, targetID)

	return &pb.ActionResponse{Success: true}, nil
}

func (u *ShareItemUsecase) updateMessageToCompleted(groupID, ownerID, targetID uuid.UUID) {
	newMeta, _ := json.Marshal(map[string]interface{}{
		"is_action_required": true,
		"is_completed":       true,
		"target_user_id":     targetID.String(),
		"type":               "GRANT_ACCESS_PROMPT",
		"completed_at":       time.Now().Unix(),
	})

	err := u.msgRepo.UpdateGrantMessageStatus(context.Background(), groupID, ownerID, targetID, string(newMeta))
	if err != nil {
		log.Printf("Failed to update message metadata: %v", err)
	}
}
