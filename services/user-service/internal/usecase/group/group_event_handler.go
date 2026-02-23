package usecase

import (
	"context"
	"time"
	"wealth-vault/user-service/internal/domain"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
)

func (u *GroupUsecase) DispatchGroupCreated(
	ctx context.Context,
	group *domain.Group,
	req *pb.CreateGroupRequest,
) {

	targetIDs := filterOut(req.MemberIds, req.CreatorId)
	if len(targetIDs) == 0 {
		return
	}

	senderName := u.getUsernameSafe(ctx, group.CreatedBy)

	evt := domain.GroupCreatedEvent{
		GroupID:       group.ID.String(),
		GroupName:     group.GroupName,
		SenderID:      req.CreatorId,
		SenderName:    senderName,
		TargetUserIDs: targetIDs,
		OccurredAt:    time.Now().Unix(),
	}

	u.publishAsync("noti.group.created", evt)
}

func (u *GroupUsecase) DispatchMemberRemoved(
	group *domain.Group,
	req *pb.RemoveMemberRequest,
	adminName string,
) {

	evt := domain.MemberRemovedEvent{
		GroupID:    req.GroupId,
		GroupName:  group.GroupName,
		TargetID:   req.TargetMemberId,
		ActionBy:   adminName,
		OccurredAt: time.Now().Unix(),
	}

	u.publishAsync("noti.group.member.removed", evt)
}

func (u *GroupUsecase) DispatchMemberLeft(
	group *domain.Group,
	req *pb.LeaveGroupRequest,
	userName string,
) {

	evt := domain.GroupActivityEvent{
		GroupID:      req.GroupId,
		ActivityType: "MEMBER_LEFT",
		Payload:      userName + " left group",
		ActorID:      req.UserId,
		TargetID:     req.UserId,
		OccurredAt:   time.Now().Unix(),
	}

	u.publishAsync("noti.group.activity", evt)
}

func filterOut(ids []string, creator string) []string {
	var result []string
	for _, id := range ids {
		if id != creator {
			result = append(result, id)
		}
	}
	return result
}
