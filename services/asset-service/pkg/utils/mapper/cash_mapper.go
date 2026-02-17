package mapper

import (
	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToCashDomain(req *pb.CreateCashRequest, userID uuid.UUID) *domain.Cash {
	return &domain.Cash{
		UserID:      userID,
		Name:        req.Name,
		Amount:      req.Amount,
		Description: req.Description,
		Files:       ToDomainFiles(req.NewFiles, userID),
	}
}

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

func ToCashProtoSlice(cashes []*domain.Cash) []*pb.Cash {
	if len(cashes) == 0 {
		return []*pb.Cash{}
	}

	res := make([]*pb.Cash, len(cashes))
	for i, c := range cashes {
		res[i] = ToCashProto(c)
	}
	return res
}
