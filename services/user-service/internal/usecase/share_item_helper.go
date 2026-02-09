package usecase

import (
	"context"
	"log"
	"wealth-vault/user-service/internal/domain"
)

func (u *ShareItemUsecase) asyncSaveGroupLogs(logs []domain.GroupLog) {
	if len(logs) == 0 {
		return
	}
	bgCtx := context.Background()
	for _, l := range logs {
		if err := u.groupRepo.CreateLog(bgCtx, &l); err != nil {
			log.Printf("⚠️ Failed to save group log: %v", err)
		}
	}
}

func (u *ShareItemUsecase) asyncSaveFriendLogs(logs []domain.FriendLog) {
	if len(logs) == 0 {
		return
	}
	bgCtx := context.Background()
	for _, l := range logs {
		if err := u.userRepo.CreateFriendLog(bgCtx, &l); err != nil {
			log.Printf("⚠️ Failed to save friend log: %v", err)
		}
	}
}

func (u *ShareItemUsecase) asyncBroadcastGroupActivities(activities []domain.GroupActivityEvent) {
	if len(activities) == 0 {
		return
	}
	for _, act := range activities {
		if err := u.publisher.Publish("noti.group.activity", act); err != nil {
			log.Printf("⚠️ Failed to publish group activity: %v", err)
		}
	}
}

func (u *ShareItemUsecase) asyncSaveGroupMessages(msgs []domain.GroupMessage) {
	if len(msgs) == 0 {
		return
	}
	bgCtx := context.Background()
	if err := u.msgRepo.CreateMessage(bgCtx, msgs); err != nil {
		log.Printf("⚠️ Failed to save group messages: %v", err)
	}
}

func (u *ShareItemUsecase) asyncSavePrivateMessages(msgs []domain.PrivateMessage) {
	if len(msgs) == 0 {
		return
	}
	bgCtx := context.Background()
	if err := u.msgRepo.CreatePrivateMessage(bgCtx, msgs); err != nil {
		log.Printf("⚠️ Failed to save private messages: %v", err)
	}
}
