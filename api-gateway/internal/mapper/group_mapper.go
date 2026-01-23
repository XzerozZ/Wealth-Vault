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
