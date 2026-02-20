package mapper

import (
	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	"wealth-vault/asset-service/pkg/utils"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MapToLands(ids []uuid.UUID) []domain.Land {
	res := make([]domain.Land, len(ids))
	for i, id := range ids {
		res[i] = domain.Land{
			ID: id,
		}
	}
	return res
}

func ToLandDomain(req *pb.CreateLandRequest, userID uuid.UUID) *domain.Land {
	return &domain.Land{
		UserID:      userID,
		Name:        req.Name,
		DeedNum:     req.DeedNum,
		Area:        req.Area,
		Amount:      req.Amount,
		Description: req.Description,
		Location:    ToLocationDomain(req.Location),
		Buildings:   MapToBuildings(utils.ParseUUIDs(req.BuildingIds)),
		Files:       ToDomainFiles(req.NewFiles, userID),
	}
}

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

func ToLandProtoSlice(lands []*domain.Land) []*pb.Land {
	if len(lands) == 0 {
		return []*pb.Land{}
	}

	res := make([]*pb.Land, len(lands))
	for i, l := range lands {
		res[i] = ToLandProto(l)
	}
	return res
}
