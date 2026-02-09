package mapper

import (
	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToCashProto(d *domain.Cash) *pb.Cash {
	res := &pb.Cash{
		Id:          d.ID.String(),
		UserId:      d.UserID.String(),
		Name:        d.Name,
		Amount:      d.Amount,
		Description: d.Description,
		Files:       ToPbFiles(d.Files),
		CreatedAt:   timestamppb.New(d.CreatedAt),
		UpdatedAt:   timestamppb.New(d.UpdatedAt),
		DeletedAt:   timestamppb.New(d.DeletedAt.Time),
	}

	return res
}
