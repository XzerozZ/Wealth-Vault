package mapper

import (
	"wealth-vault/api-gateway/internal/domain"
	pb "wealth-vault/api-gateway/pkg/pb/proto/user"
)

func ToGroupDomain(p *pb.Group) *domain.Group {
	if p == nil {
		return nil
	}

	return &domain.Group{
		ID:           p.Id,
		GroupName:    p.Name,
		GroupProfile: p.Profile,
		CreatedBy:    p.UserId,
		MemberCount:  p.MemberCount,
		CreatedAt:    p.CreatedAt.AsTime(),
		UpdatedAt:    p.UpdatedAt.AsTime(),
	}
}

func ToGroupList(pbList []*pb.Group) []domain.Group {
	if pbList == nil {
		return []domain.Group{}
	}

	entities := make([]domain.Group, 0, len(pbList))

	for _, pbItem := range pbList {
		if item := ToGroupDomain(pbItem); item != nil {
			entities = append(entities, *item)
		}
	}

	return entities
}
