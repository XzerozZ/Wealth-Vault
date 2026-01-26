package mapper

import (
	"wealth-vault/api-gateway/internal/domain"
	pb "wealth-vault/api-gateway/pkg/pb/proto/asset"
)

func ToLocationDomain(p *pb.Location) *domain.Location {
	if p == nil {
		return nil
	}

	return &domain.Location{
		ID:          p.Id,
		Address:     p.Address,
		Subdistrict: p.Subdistrict,
		District:    p.District,
		Province:    p.Province,
		PostalCode:  p.PostalCode,
		CreatedAt:   p.CreatedAt.AsTime(),
		UpdatedAt:   p.UpdatedAt.AsTime(),
	}
}
