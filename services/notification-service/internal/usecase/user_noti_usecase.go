package usecase

import (
	"context"
	"fmt"
	"wealth-vault/notification-service/internal/domain"
	"wealth-vault/notification-service/pkg/utils"

	"github.com/google/uuid"
)

func (u *NotificationUsecase) HandleGroupCreated(ctx context.Context, evt domain.GroupCreatedEvent) error {
	groupUUID, err := uuid.Parse(evt.GroupID)
	if err != nil {
		return err
	}

	senderUUID := utils.ParseUUIDPtr(evt.SenderID)
	msg := fmt.Sprintf("👋 %s ได้เพิ่มคุณลงในกลุ่ม '%s'", evt.SenderName, evt.GroupName)

	u.notifyMany(ctx, evt.TargetUserIDs, senderUUID, "GROUP_INVITE", groupUUID, msg, evt.OccurredAt)
	return nil
}

func (u *NotificationUsecase) HandleGroupMemberAdded(ctx context.Context, evt domain.GroupMemberAddedEvent) error {
	groupUUID, err := uuid.Parse(evt.GroupID)
	if err != nil {
		return err
	}

	senderUUID := utils.ParseUUIDPtr(evt.SenderID)
	msg := "มีสมาชิกใหม่ถูกเชิญเข้ากลุ่ม มาเริ่มแชร์รายการกันเถอะ 💸"

	u.notifyMany(ctx, evt.TargetUserIDs, senderUUID, "GROUP", groupUUID, msg, evt.OccurredAt)

	u.emitToGroup(evt.GroupID, "MEMBER_ADDED", evt)
	return nil
}

func (u *NotificationUsecase) HandleItemShared(ctx context.Context, evt domain.ItemSharedEvent) error {
	assetUUID, err := uuid.Parse(evt.AssetID)
	if err != nil {
		return err
	}

	senderUUID := utils.ParseUUIDPtr(evt.SenderID)
	msg := fmt.Sprintf("%s ได้แชร์รายการใหม่ให้กับคุณ", evt.SenderName)

	u.notifyMany(ctx, evt.TargetUserIDs, senderUUID, "ASSET", assetUUID, msg, evt.OccurredAt)
	return nil
}

func (u *NotificationUsecase) HandleFriendRequest(ctx context.Context, evt domain.FriendRequestEvent) error {
	targetUUID, err := uuid.Parse(evt.TargetID)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("👋 %s ได้ส่งคำขอเป็นเพื่อนกับคุณ", evt.SenderName)
	senderUUID := utils.ParseUUIDPtr(evt.SenderID)
	metadata := map[string]interface{}{
		"is_action_required": true,
		"is_completed":       false,
		"type":               "FRIEND_REQUEST",
	}

	u.notifyTarget(ctx, targetUUID, senderUUID, "FRIEND_REQUEST", uuid.Nil, msg, evt.OccurredAt, metadata)
	return nil
}

func (u *NotificationUsecase) HandleFriendAccepted(ctx context.Context, evt domain.FriendAcceptedEvent) error {
	requesterUUID, err := uuid.Parse(evt.RequesterID)
	if err != nil {
		return err
	}

	senderUUID := utils.ParseUUIDPtr(evt.AccepterID)
	accepterUUID := uuid.MustParse(evt.AccepterID)

	go func() {
		err := u.repo.UpdateNotificationMetadata(context.Background(), accepterUUID, requesterUUID, "FRIEND_REQUEST", map[string]interface{}{
			"is_completed": true,
			"status":       "ACCEPTED",
		})
		if err != nil {
			fmt.Printf("failed to update old notification: %v\n", err)
		}
	}()

	msg := fmt.Sprintf("✅ %s ได้ตอบรับคำขอเป็นเพื่อนของคุณแล้ว", evt.AccepterName)
	u.notifyTarget(ctx, requesterUUID, senderUUID, "FRIEND_ACCEPTED", uuid.Nil, msg, evt.OccurredAt, nil)
	return nil
}

func (u *NotificationUsecase) HandleFriendDecline(ctx context.Context, evt domain.FriendAcceptedEvent) error {
	requesterUUID, err := uuid.Parse(evt.RequesterID)
	if err != nil {
		return err
	}

	senderUUID := utils.ParseUUIDPtr(evt.AccepterID)
	accepterUUID := uuid.MustParse(evt.AccepterID)

	go func() {
		err := u.repo.UpdateNotificationMetadata(context.Background(), accepterUUID, requesterUUID, "FRIEND_REQUEST", map[string]interface{}{
			"is_completed": true,
			"status":       "DECLINE",
		})
		if err != nil {
			fmt.Printf("failed to update old notification: %v\n", err)
		}
	}()

	msg := fmt.Sprintf("✅ %s ได้ปฏิเสธคำขอเป็นเพื่อนของคุณแล้ว", evt.AccepterName)
	u.notifyTarget(ctx, requesterUUID, senderUUID, "FRIEND_DECLINE", uuid.Nil, msg, evt.OccurredAt, nil)
	return nil
}

func (u *NotificationUsecase) HandleAccessGranted(ctx context.Context, evt domain.AccessGrantedEvent) error {
	targetUUID, err := uuid.Parse(evt.TargetUserID)
	if err != nil {
		return err
	}

	groupUUID, err := uuid.Parse(evt.GroupID)
	if err != nil {
		return err
	}

	senderUUID := utils.ParseUUIDPtr(evt.GrantorID)
	msg := fmt.Sprintf("คุณได้รับสิทธิ์เข้าถึงรายการทรัพย์สินย้อนหลัง %d รายการ จาก %s", evt.ItemCount, evt.GrantorName)
	u.notifyTarget(ctx, targetUUID, senderUUID, "ACCESS_GRANTED", groupUUID, msg, evt.OccurredAt, nil)
	u.emitToUser(evt.TargetUserID, "ACCESS_GRANTED", map[string]interface{}{
		"group_id": evt.GroupID,
		"count":    evt.ItemCount,
	})

	return nil
}

func (u *NotificationUsecase) HandleMemberRemoved(ctx context.Context, evt domain.MemberRemovedEvent) error {
	targetUUID, err := uuid.Parse(evt.TargetID)
	if err != nil {
		return err
	}

	groupUUID, err := uuid.Parse(evt.GroupID)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("คุณถูกลบออกจากกลุ่ม '%s'", evt.GroupName)
	u.notifyTarget(ctx, targetUUID, nil, "GROUP_REMOVED", groupUUID, msg, evt.OccurredAt, nil)
	u.emitToUser(evt.TargetID, "YOU_ARE_REMOVED", map[string]interface{}{
		"group_id": evt.GroupID,
	})
	return nil
}

func (u *NotificationUsecase) HandleGroupActivity(ctx context.Context, evt domain.GroupActivityEvent) error {
	u.emitToGroup(evt.GroupID, evt.ActivityType, evt)
	return nil
}
