package mapper

import (
	"strings"
	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var domainToProtoBuildType = map[string]pb.BuildingType{
	"CONDO":      pb.BuildingType_BUILDING_TYPE_CONDO,
	"HOUSE":      pb.BuildingType_BUILDING_TYPE_HOUSE,
	"TOWNHOME":   pb.BuildingType_BUILDING_TYPE_TOWNHOME,
	"COMMERCIAL": pb.BuildingType_BUILDING_TYPE_COMMERCIAL,
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
