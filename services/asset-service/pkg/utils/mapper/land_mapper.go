package mapper

import (
	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToLandProto(d *domain.Land) *pb.Land {
	res := &pb.Land{
		Id:          d.ID.String(),
		UserId:      d.UserID.String(),
		Name:        d.Name,
		DeedNum:     d.DeedNum,
		Area:        d.Area,
		Amount:      d.Amount,
		Description: d.Description,
		Location:    ToLocationProto(&d.Location),
		Files:       ToPbFiles(d.Files),
		Ref:         ToRefBuildingProto(d.Buildings),
		CreatedAt:   timestamppb.New(d.CreatedAt),
		UpdatedAt:   timestamppb.New(d.UpdatedAt),
		DeletedAt:   timestamppb.New(d.DeletedAt.Time),
	}

	return res
}
