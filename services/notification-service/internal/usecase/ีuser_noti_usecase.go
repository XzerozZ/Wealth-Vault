package usecase

import (
	"fmt"
	"log"
	"time"
	"wealth-vault/notification-service/internal/domain"
	"wealth-vault/notification-service/pkg/utils"

	"github.com/google/uuid"
)

func (u *NotificationUsecase) notifyTarget(receiverID uuid.UUID, senderID *uuid.UUID, entityType string, entityID uuid.UUID, message string, occurredAt int64) {
	noti := &domain.Notification{
		ID:         uuid.New(),
		EntityType: entityType,
		EntityID:   entityID,
		Receiver:   receiverID,
		SenderID:   senderID,
		Message:    message,
		Channel:    "IN_APP",
		CreatedAt:  time.Unix(occurredAt, 0),
		IsRead:     false,
	}

	if err := u.repo.CreateNotification(noti); err != nil {
		log.Printf("❌ [Notification] Save DB Error: %v", err)
	}

	u.hub.Emit(receiverID.String(), map[string]interface{}{
		"type":    "NOTIFICATION",
		"payload": noti,
	})
}

func (u *NotificationUsecase) HandleGroupCreated(evt domain.GroupCreatedEvent) {
	groupUUID, err := uuid.Parse(evt.GroupID)
	if err != nil {
		log.Printf("⚠️ Invalid GroupID: %s", evt.GroupID)
		return
	}

	senderUUID := utils.ParseUUIDPtr(evt.SenderID)
	msg := fmt.Sprintf("👋 %s ได้เพิ่มคุณลงในกลุ่ม '%s'", evt.SenderName, evt.GroupName)

	for _, targetIDStr := range evt.TargetUserIDs {
		targetID, err := uuid.Parse(targetIDStr)
		if err != nil {
			continue
		}
		u.notifyTarget(targetID, senderUUID, "GROUP_INVITE", groupUUID, msg, evt.OccurredAt)
	}
}

func (u *NotificationUsecase) HandleGroupMemberAdded(evt domain.GroupMemberAddedEvent) {
	groupUUID, err := uuid.Parse(evt.GroupID)
	if err != nil {
		return
	}

	senderUUID := utils.ParseUUIDPtr(evt.SenderID)
	msg := "มีสมาชิกใหม่ถูกเชิญเข้ากลุ่ม มาเริ่มแชร์รายการกันเถอะ 💸"

	for _, targetIDStr := range evt.TargetUserIDs {
		receiverID, err := uuid.Parse(targetIDStr)
		if err != nil {
			continue
		}
		if senderUUID != nil && *senderUUID == receiverID {
			continue
		}
		u.notifyTarget(receiverID, senderUUID, "GROUP", groupUUID, msg, evt.OccurredAt)
	}

	u.hub.BroadcastToGroup(evt.GroupID, map[string]interface{}{
		"type":  "DATA_UPDATE",
		"event": "MEMBER_ADDED",
		"payload": map[string]interface{}{
			"group_id":       evt.GroupID,
			"added_user_ids": evt.AddedUserIDs,
			"action_by":      evt.SenderID,
			"occurred_at":    evt.OccurredAt,
		},
	})
}

func (u *NotificationUsecase) HandleItemShared(evt domain.ItemSharedEvent) {
	assetUUID, err := uuid.Parse(evt.AssetID)
	if err != nil {
		return
	}

	senderUUID := utils.ParseUUIDPtr(evt.SenderID)
	msg := fmt.Sprintf("%s ได้แชร์รายการใหม่ให้กับคุณ", evt.SenderName)

	for _, targetIDStr := range evt.TargetUserIDs {
		receiverID, err := uuid.Parse(targetIDStr)
		if err != nil {
			continue
		}
		if senderUUID != nil && *senderUUID == receiverID {
			continue
		}
		u.notifyTarget(receiverID, senderUUID, "ASSET", assetUUID, msg, evt.OccurredAt)
	}
}

func (u *NotificationUsecase) HandleFriendRequest(evt domain.FriendRequestEvent) {
	targetUUID, err := uuid.Parse(evt.TargetID)
	if err != nil {
		return
	}

	msg := fmt.Sprintf("👋 %s ได้ส่งคำขอเป็นเพื่อนกับคุณ", evt.SenderName)
	senderUUID := utils.ParseUUIDPtr(evt.SenderID)
	u.notifyTarget(targetUUID, senderUUID, "FRIEND_REQUEST", uuid.Nil, msg, evt.OccurredAt)
}

func (u *NotificationUsecase) HandleFriendAccepted(evt domain.FriendAcceptedEvent) {
	requesterUUID, err := uuid.Parse(evt.RequesterID)
	if err != nil {
		return
	}

	senderUUID := utils.ParseUUIDPtr(evt.AccepterID)
	msg := fmt.Sprintf("✅ %s ได้ตอบรับคำขอเป็นเพื่อนของคุณแล้ว", evt.AccepterName)
	u.notifyTarget(requesterUUID, senderUUID, "FRIEND_ACCEPTED", uuid.Nil, msg, evt.OccurredAt)
}

func (u *NotificationUsecase) HandleAccessGranted(evt domain.AccessGrantedEvent) {
	targetUUID, err := uuid.Parse(evt.TargetUserID)
	if err != nil {
		return
	}

	groupUUID, _ := uuid.Parse(evt.GroupID)
	senderUUID := utils.ParseUUIDPtr(evt.GrantorID)
	msg := fmt.Sprintf("คุณได้รับสิทธิ์เข้าถึงรายการทรัพย์สินย้อนหลัง %d รายการ จาก %s", evt.ItemCount, evt.GrantorName)
	u.notifyTarget(targetUUID, senderUUID, "ACCESS_GRANTED", groupUUID, msg, evt.OccurredAt)
	u.hub.Emit(evt.TargetUserID, map[string]interface{}{
		"type":  "DATA_UPDATE",
		"event": "ACCESS_GRANTED",
		"payload": map[string]interface{}{
			"group_id": evt.GroupID,
			"count":    evt.ItemCount,
		},
	})
}

func (u *NotificationUsecase) HandleMemberRemoved(evt domain.MemberRemovedEvent) {
	targetUUID, err := uuid.Parse(evt.TargetID)
	if err != nil {
		return
	}

	groupUUID, _ := uuid.Parse(evt.GroupID)
	msg := fmt.Sprintf("คุณถูกลบออกจากกลุ่ม '%s'", evt.GroupName)
	u.notifyTarget(targetUUID, nil, "GROUP_REMOVED", groupUUID, msg, evt.OccurredAt)
	u.hub.Emit(evt.TargetID, map[string]interface{}{
		"type":  "DATA_UPDATE",
		"event": "YOU_ARE_REMOVED",
		"payload": map[string]interface{}{
			"group_id": evt.GroupID,
		},
	})
}

func (u *NotificationUsecase) HandleGroupActivity(evt domain.GroupActivityEvent) {
	u.hub.BroadcastToGroup(evt.GroupID, map[string]interface{}{
		"type":    "DATA_UPDATE",
		"event":   evt.ActivityType,
		"payload": evt,
	})
}
