package usecase_test

import (
	"context"
	"testing"
	"time"
	"wealth-vault/user-service/internal/domain"
	usecase "wealth-vault/user-service/internal/usecase/group"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
	mock_event "wealth-vault/user-service/test/mock/event"
	mock_repo "wealth-vault/user-service/test/mock/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

func TestGroupUsecase_DispatchEvents(t *testing.T) {
	ctx := context.Background()

	t.Run("DispatchGroupCreated - should publish event when members exist", func(t *testing.T) {
		mockPublisher := new(mock_event.MockEventPublisher)
		mockUserRepo := new(mock_repo.MockUserRepository)
		uc := usecase.NewGroupUsecase(nil, mockUserRepo, nil, nil, mockPublisher)

		creatorID := uuid.New()
		targetID := uuid.New()
		group := &domain.Group{
			ID:        uuid.New(),
			GroupName: "My Family",
			CreatedBy: creatorID,
		}
		req := &pb.CreateGroupRequest{
			CreatorId: creatorID.String(),
			Name:      "My Family",
			MemberIds: []string{targetID.String()},
		}

		mockUserRepo.On("GetUser", mock.Anything, creatorID).
			Return(&domain.User{Username: "OwnerName"}, nil).Once()

		mockPublisher.On("Publish", "noti.group.created", mock.MatchedBy(func(evt domain.GroupCreatedEvent) bool {
			return evt.GroupName == "My Family" &&
				evt.SenderName == "OwnerName" &&
				len(evt.TargetUserIDs) == 1
		})).Return(nil).Once()

		uc.DispatchGroupCreated(ctx, group, req)
		time.Sleep(50 * time.Millisecond)
		mockPublisher.AssertExpectations(t)
	})

	t.Run("DispatchMemberRemoved - success", func(t *testing.T) {
		mockPublisher := new(mock_event.MockEventPublisher)
		uc := usecase.NewGroupUsecase(nil, nil, nil, nil, mockPublisher)

		gid := uuid.New().String()
		targetID := uuid.New().String()
		group := &domain.Group{GroupName: "Work Group"}
		req := &pb.RemoveMemberRequest{
			GroupId:        gid,
			TargetMemberId: targetID,
		}
		adminName := "Boss"

		mockPublisher.On("Publish", "noti.group.member.removed", mock.MatchedBy(func(evt domain.MemberRemovedEvent) bool {
			return evt.GroupID == gid && evt.ActionBy == adminName && evt.TargetID == targetID
		})).Return(nil).Once()

		uc.DispatchMemberRemoved(group, req, adminName)
		time.Sleep(50 * time.Millisecond)
		mockPublisher.AssertExpectations(t)
	})

	t.Run("DispatchMemberLeft - success", func(t *testing.T) {
		mockPublisher := new(mock_event.MockEventPublisher)
		uc := usecase.NewGroupUsecase(nil, nil, nil, nil, mockPublisher)

		gid := uuid.New().String()
		uid := uuid.New().String()
		group := &domain.Group{GroupName: "Yoga Club"}
		req := &pb.LeaveGroupRequest{
			GroupId: gid,
			UserId:  uid,
		}
		userName := "Alice"

		mockPublisher.On("Publish", "noti.group.activity", mock.MatchedBy(func(evt domain.GroupActivityEvent) bool {
			return evt.ActivityType == "MEMBER_LEFT" &&
				evt.ActorID == uid &&
				evt.Payload == "Alice left group"
		})).Return(nil).Once()

		uc.DispatchMemberLeft(group, req, userName)
		time.Sleep(50 * time.Millisecond)
		mockPublisher.AssertExpectations(t)
	})
}
