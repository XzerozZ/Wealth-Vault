package mapper

import (
	"strings"
	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	"wealth-vault/asset-service/pkg/utils"
	"wealth-vault/asset-service/pkg/utils/helper"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var domainToProtoBuildType = map[string]pb.BuildingType{
	"CONDO":      pb.BuildingType_BUILDING_TYPE_CONDO,
	"HOUSE":      pb.BuildingType_BUILDING_TYPE_HOUSE,
	"TOWNHOME":   pb.BuildingType_BUILDING_TYPE_TOWNHOME,
	"COMMERCIAL": pb.BuildingType_BUILDING_TYPE_COMMERCIAL,
}

func MapToBuildings(ids []uuid.UUID) []domain.Building {
	res := make([]domain.Building, len(ids))
	for i, id := range ids {
		res[i] = domain.Building{
			ID: id,
		}
	}
	return res
}

func ToBuildingDomain(req *pb.CreateBuildingRequest, userID uuid.UUID) *domain.Building {
	buildType := domain.BuildingTypeHouse
	if val, ok := helper.ProtoToDomainBuildingType[req.Type]; ok {
		buildType = val
	}

	return &domain.Building{
		UserID:      userID,
		Name:        req.Name,
		Type:        buildType,
		Area:        req.Area,
		Amount:      req.Amount,
		Description: req.Description,
		Location:    ToLocationDomain(req.Location),
		Lands:       MapToLands(utils.ParseUUIDs(req.LandIds)),
		Insurances:  MapToInsurances(utils.ParseUUIDs(req.InsIds)),
		Files:       ToDomainFiles(req.NewFiles, userID),
	}
}

func ToBuildingProto(d *domain.Building) *pb.Building {
	normalizedType := strings.ToUpper(strings.TrimSpace(string(d.Type)))
	buTypeEnum, ok := domainToProtoBuildType[normalizedType]
	if !ok {
		buTypeEnum = pb.BuildingType_BUILDING_TYPE_UNSPECIFIED
	}

	res := &pb.Building{
		Id:          d.ID.String(),
		UserId:      d.UserID.String(),
		Name:        d.Name,
		Type:        buTypeEnum,
		Area:        d.Area,
		Amount:      d.Amount,
		Description: d.Description,
		Location:    ToLocationProto(&d.Location),
		Files:       ToPbFiles(d.Files),
		Ref:         ToRefLandProto(d.Lands),
		Ins:         ToRefInsProto(d.Insurances),
		CreatedAt:   timestamppb.New(d.CreatedAt),
		UpdatedAt:   timestamppb.New(d.UpdatedAt),
		DeletedAt:   timestamppb.New(d.DeletedAt.Time),
	}

	return res
}

func ToBuildingProtoSlice(buidings []*domain.Building) []*pb.Building {
	if len(buidings) == 0 {
		return []*pb.Building{}
	}

	res := make([]*pb.Building, len(buidings))
	for i, b := range buidings {
		res[i] = ToBuildingProto(b)
	}
	return res
}
