package mapper

import (
	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToLocationProto(d *domain.Location) *pb.Location {
	res := &pb.Location{
		Id:          d.ID.String(),
		Address:     d.Address,
		Subdistrict: d.Subdistrict,
		District:    d.District,
		Province:    d.Province,
		PostalCode:  d.PostalCode,
		CreatedAt:   timestamppb.New(d.CreatedAt),
		UpdatedAt:   timestamppb.New(d.UpdatedAt),
	}

	return res
}
