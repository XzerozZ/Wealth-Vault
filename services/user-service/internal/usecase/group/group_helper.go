package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"wealth-vault/user-service/internal/domain"

	"github.com/google/uuid"
)

func (u *GroupUsecase) ensureMember(
	ctx context.Context,
	groupID, userID uuid.UUID,
) error {
	isMember, err := u.groupRepo.IsUserMember(ctx, groupID, userID)
	if err != nil {
		return errors.New("failed to check membership")
	}
	if !isMember {
		return errors.New("access denied")
	}
	return nil
}

func (u *GroupUsecase) getUsernameSafe(
	ctx context.Context,
	userID uuid.UUID,
) string {
	user, err := u.userRepo.GetUser(ctx, userID)
	if err != nil || user == nil {
		return "Unknown"
	}
	return user.Username
}

func (u *GroupUsecase) buildSystemLog(
	groupID uuid.UUID,
	userID uuid.UUID,
	message string,
	meta map[string]interface{},
) *domain.GroupLog {

	var metaStr string
	if meta != nil {
		b, _ := json.Marshal(meta)
		metaStr = string(b)
	}

	return &domain.GroupLog{
		GroupID:   groupID,
		LogType:   "SYSTEM",
		Messages:  message,
		Metadata:  metaStr,
		CreatedBy: userID,
	}
}

func (u *GroupUsecase) publishAsync(topic string, payload interface{}) {
	go func() {
		_ = u.publisher.Publish(topic, payload)
	}()
}
