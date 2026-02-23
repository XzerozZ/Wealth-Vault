package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
	"wealth-vault/user-service/internal/domain"
	assetPb "wealth-vault/user-service/pkg/pb/proto/asset"
	pb "wealth-vault/user-service/pkg/pb/proto/user"

	"github.com/google/uuid"
)

type ShareAggregator struct {
	groupItems  []domain.GroupItem
	friendItems []domain.FriendItem
	emailItems  []domain.EmailItem
	emailsNow   []domain.EmailItem
	groupLogs   []domain.GroupLog
	friendLogs  []domain.FriendLog
	groupActs   []domain.GroupActivityEvent
	groupMsgs   []domain.GroupMessage
	privateMsgs []domain.PrivateMessage
}

func (u *ShareItemUsecase) ShareItem(ctx context.Context, req *pb.ShareItemRequest) (*pb.ShareItemResponse, error) {
	if len(req.ItemIds) == 0 {
		return nil, errors.New("no items to share")
	}
	if len(req.ItemIds) != len(req.ItemTypes) {
		return nil, errors.New("mismatch between item_ids and types length")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %v", err)
	}

	senderName := "Unknown"
	if userProfile, err := u.userRepo.GetUser(ctx, userID); err == nil {
		senderName = userProfile.Username
	}

	now := time.Now()
	agg := &ShareAggregator{}

	for i, idStr := range req.ItemIds {
		entityType := req.ItemTypes[i]
		entityID, err := uuid.Parse(idStr)
		if err != nil {
			log.Printf("⚠️ skip invalid entity id %s: %v", idStr, err)
			continue
		}

		res, err := u.assetClient.CheckAssetExists(ctx, &assetPb.CheckAssetRequest{
			Id:     idStr,
			UserId: req.UserId,
			Type:   entityType,
		})
		if err != nil || !res.Exists {
			return nil, fmt.Errorf("asset not found: %s", idStr)
		}

		assetDisplayName := fmt.Sprintf("%s shared item", entityType)
		assetNotifyMap := make(map[string]bool)

		u.prepareGroupShares(ctx, req.Groups, userID, senderName, entityID, entityType, assetDisplayName, now, agg, assetNotifyMap)
		u.prepareFriendShares(ctx, req.Friends, userID, entityID, entityType, assetDisplayName, now, agg, assetNotifyMap)
		u.prepareEmailShares(req.Emails, userID, entityID, entityType, now, agg)

		if len(assetNotifyMap) > 0 {
			var targetIDs []string
			for uid := range assetNotifyMap {
				targetIDs = append(targetIDs, uid)
			}
			go u.publisher.Publish(TopicItemShared, domain.ItemSharedEvent{
				SenderID:      req.UserId,
				SenderName:    senderName,
				AssetID:       idStr,
				TargetUserIDs: targetIDs,
				OccurredAt:    now.Unix(),
			})
		}
	}

	return u.executeShareTransaction(ctx, agg)
}

func (u *ShareItemUsecase) prepareGroupShares(ctx context.Context, targets []*pb.ShareTarget, userID uuid.UUID, senderName string, entityID uuid.UUID, entityType string, assetDisplayName string, now time.Time, agg *ShareAggregator, notifyMap map[string]bool) {
	for _, target := range targets {
		groupID, err := uuid.Parse(target.Id)
		if err != nil {
			continue
		}

		shareTime := now
		if target.ShareAt != nil {
			shareTime = target.ShareAt.AsTime()
		}

		exist, _ := u.itemRepo.IsItemSharedtoGroup(ctx, groupID, entityID, entityType)
		if exist {
			continue
		}

		newItem := domain.GroupItem{
			GroupID:    groupID,
			EntityType: entityType,
			EntityID:   entityID,
			OwnerID:    userID,
			ShareAt:    shareTime,
		}

		if members, _, err := u.groupRepo.GetMember(ctx, groupID); err == nil {
			for _, member := range members {
				newItem.Viewers = append(newItem.Viewers, domain.GroupItemViewer{
					GroupItemID: newItem.ID,
					ViewerID:    member.ID,
				})
				if member.ID != userID {
					notifyMap[member.ID.String()] = true
				}
			}
		}

		agg.groupItems = append(agg.groupItems, newItem)

		logMeta, _ := json.Marshal(map[string]string{
			"action":    "share_item",
			"item_type": entityType,
			"item_id":   entityID.String(),
			"shared_at": shareTime.Format(time.RFC3339),
		})

		agg.groupLogs = append(agg.groupLogs, domain.GroupLog{
			GroupID:   groupID,
			LogType:   LogTypeActivity,
			Messages:  fmt.Sprintf("%s ได้แชร์ %s รายการใหม่เข้ากลุ่ม", senderName, entityType),
			Metadata:  string(logMeta),
			CreatedBy: userID,
		})

		agg.groupActs = append(agg.groupActs, domain.GroupActivityEvent{
			GroupID:      target.Id,
			ActivityType: "ITEM_SHARED",
			Payload:      fmt.Sprintf("%s แชร์ %s ใหม่", senderName, entityType),
			ActorID:      userID.String(),
			OccurredAt:   now.Unix(),
		})

		cardMeta, _ := json.Marshal(map[string]interface{}{
			"asset_id":       entityID.String(),
			"asset_type":     entityType,
			"snapshot_title": assetDisplayName,
			"action_url":     fmt.Sprintf("/asset/%s/%s", entityType, entityID),
		})

		agg.groupMsgs = append(agg.groupMsgs, domain.GroupMessage{
			GroupID:   groupID,
			SenderID:  userID,
			MsgType:   MsgTypeAssetCard,
			Content:   "",
			Metadata:  string(cardMeta),
			CreatedAt: now,
		})
	}
}

func (u *ShareItemUsecase) prepareFriendShares(ctx context.Context, targets []*pb.ShareTarget, userID uuid.UUID, entityID uuid.UUID, entityType string, assetDisplayName string, now time.Time, agg *ShareAggregator, notifyMap map[string]bool) {
	for _, target := range targets {
		friendID, err := uuid.Parse(target.Id)
		if err != nil {
			continue
		}

		shareTime := now
		if target.ShareAt != nil {
			shareTime = target.ShareAt.AsTime()
		}

		exist, _ := u.itemRepo.IsItemSharedtoFriend(ctx, friendID, entityID, entityType)
		if exist {
			continue
		}

		agg.friendItems = append(agg.friendItems, domain.FriendItem{
			OwnerID:    userID,
			FriendID:   friendID,
			EntityType: entityType,
			EntityID:   entityID,
			ShareAt:    shareTime,
		})

		logMeta, _ := json.Marshal(map[string]string{
			"action":    "share_to_friend",
			"item_type": entityType,
			"item_id":   entityID.String(),
		})

		agg.friendLogs = append(agg.friendLogs, domain.FriendLog{
			OwnerID:   userID,
			FriendID:  friendID,
			LogType:   LogTypeActivity,
			Messages:  fmt.Sprintf("คุณได้แชร์ %s รายการใหม่ให้กับเพื่อน", entityType),
			Metadata:  string(logMeta),
			CreatedBy: userID,
		})

		cardMeta, _ := json.Marshal(map[string]interface{}{
			"asset_id":       entityID.String(),
			"asset_type":     entityType,
			"snapshot_title": assetDisplayName,
		})

		agg.privateMsgs = append(agg.privateMsgs, domain.PrivateMessage{
			SenderID:   userID,
			ReceiverID: friendID,
			MsgType:    MsgTypeAssetCard,
			Metadata:   string(cardMeta),
			CreatedAt:  now,
		})

		notifyMap[target.Id] = true
	}
}

func (u *ShareItemUsecase) prepareEmailShares(targets []*pb.ShareTarget, userID uuid.UUID, entityID uuid.UUID, entityType string, now time.Time, agg *ShareAggregator) {
	for _, target := range targets {
		shareTime := now
		if target.ShareAt != nil {
			shareTime = target.ShareAt.AsTime()
		}

		shouldSendNow := shareTime.IsZero() || shareTime.Before(now.Add(1*time.Minute))
		emailItem := domain.EmailItem{
			OwnerID:    userID,
			Email:      target.Id,
			EntityType: entityType,
			EntityID:   entityID,
			ShareAt:    shareTime,
			IsSent:     shouldSendNow,
		}

		agg.emailItems = append(agg.emailItems, emailItem)
		if shouldSendNow {
			agg.emailsNow = append(agg.emailsNow, emailItem)
		}
	}
}

func (u *ShareItemUsecase) executeShareTransaction(ctx context.Context, agg *ShareAggregator) (*pb.ShareItemResponse, error) {
	if len(agg.groupItems) > 0 {
		if err := u.itemRepo.ShareItemtoGroup(ctx, agg.groupItems); err != nil {
			return nil, err
		}

		go u.asyncSaveGroupLogs(agg.groupLogs)
		go u.asyncBroadcastGroupActivities(agg.groupActs)
		go u.asyncSaveGroupMessages(agg.groupMsgs)
	}

	if len(agg.friendItems) > 0 {
		if err := u.itemRepo.ShareItemtoFriend(ctx, agg.friendItems); err != nil {
			return nil, err
		}

		go u.asyncSaveFriendLogs(agg.friendLogs)
		go u.asyncSavePrivateMessages(agg.privateMsgs)
	}

	if len(agg.emailItems) > 0 {
		if err := u.itemRepo.ShareItemtoEmail(ctx, agg.emailItems); err != nil {
			return nil, err
		}
		if len(agg.emailsNow) > 0 {
			go u.SendEmailInvitations(agg.emailsNow)
		}
	}

	return &pb.ShareItemResponse{Finish: true}, nil
}
