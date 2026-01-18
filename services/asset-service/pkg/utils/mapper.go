package utils

import (
	"encoding/json"
	"strings"
	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToPbFiles(files []domain.FileAssociate) []*pb.FileInfo {
	if len(files) == 0 {
		return []*pb.FileInfo{}
	}

	pbFiles := make([]*pb.FileInfo, len(files))
	for i, f := range files {
		pbFiles[i] = &pb.FileInfo{
			Id:       f.ID.String(),
			Url:      f.Link,
			FileType: f.FileType,
		}
	}

	return pbFiles
}

func MapAssetTypeToProto(t domain.AssetType) pb.AssetType {
	switch t {
	case domain.AssetTypeCash:
		return pb.AssetType_ASSET_TYPE_CASH
	case domain.AssetTypeBank:
		return pb.AssetType_ASSET_TYPE_BANK
	case domain.AssetTypeInvestment:
		return pb.AssetType_ASSET_TYPE_INVESTMENT
	case domain.AssetTypeRealEstate:
		return pb.AssetType_ASSET_TYPE_REAL_ESTATE
	case domain.AssetTypeInsurance:
		return pb.AssetType_ASSET_TYPE_INSURANCE
	default:
		return pb.AssetType_ASSET_TYPE_UNSPECIFIED
	}
}

func ToAssetProto(d *domain.Asset) *pb.Asset {
	rawType := strings.ToUpper(strings.TrimSpace(string(d.Type)))
	var assetTypeEnum pb.AssetType
	if v, ok := pb.AssetType_value[rawType]; ok {
		assetTypeEnum = pb.AssetType(v)
	} else {
		if v, ok := pb.AssetType_value["ASSET_TYPE_"+rawType]; ok {
			assetTypeEnum = pb.AssetType(v)
		} else {
			assetTypeEnum = pb.AssetType_ASSET_TYPE_UNSPECIFIED
		}
	}

	res := &pb.Asset{
		Id:                   d.ID.String(),
		Name:                 d.Name,
		Amount:               d.Amount,
		IsIncludedInNetWorth: d.IsIncludedInNetWorth,
		Type:                 assetTypeEnum,
		Description:          d.Description,
		CreatedBy:            d.UserID.String(),
		CreatedAt:            timestamppb.New(d.CreatedAt),
		UpdatedAt:            timestamppb.New(d.UpdatedAt),
		Files:                ToPbFiles(d.Files),
	}

	if len(d.Details) > 0 && string(d.Details) != "null" && string(d.Details) != "{}" {
		switch assetTypeEnum {
		case pb.AssetType_ASSET_TYPE_BANK:
			var info domain.BankDetail
			if err := json.Unmarshal(d.Details, &info); err == nil {
				res.Detail = &pb.Asset_BankDetail{
					BankDetail: &pb.BankDetail{
						BankName:      info.BankName,
						AccountNumber: info.AccountNumber,
						AccountType:   pb.BankAccountType(info.AccountType),
					},
				}
			}

		case pb.AssetType_ASSET_TYPE_INVESTMENT:
			var info domain.InvestmentDetail
			if err := json.Unmarshal(d.Details, &info); err == nil {
				res.Detail = &pb.Asset_InvestmentDetail{
					InvestmentDetail: &pb.InvestmentDetail{
						SubType:    pb.InvestmentType(info.SubType),
						Symbol:     info.Symbol,
						BrokerName: info.Broker,
						Quantity:   info.Quantity,
						CostPrice:  info.CostPrice,
					},
				}
			}

		case pb.AssetType_ASSET_TYPE_REAL_ESTATE:
			var info domain.RealEstateDetail
			if err := json.Unmarshal(d.Details, &info); err == nil {
				var loc *pb.Location
				if info.Location.Address != "" {
					loc = &pb.Location{
						Address:     info.Location.Address,
						SubDistrict: info.Location.SubDistrict,
						District:    info.Location.District,
						Province:    info.Location.Province,
						PostalCode:  info.Location.PostalCode,
					}
				}
				res.Detail = &pb.Asset_RealEstateDetail{
					RealEstateDetail: &pb.RealEstateDetail{
						PropertyType: pb.RealEstateType(info.PropertyType),
						DeedNumber:   info.DeedNumber,
						AreaSqm:      info.AreaSqm,
						Location:     loc,
					},
				}
			}

		case pb.AssetType_ASSET_TYPE_INSURANCE:
			var info domain.InsuranceDetail
			if err := json.Unmarshal(d.Details, &info); err == nil {
				res.Detail = &pb.Asset_InsuranceDetail{
					InsuranceDetail: &pb.InsuranceDetail{
						SubType:        pb.InsuranceType(info.SubType),
						PolicyNumber:   info.PolicyNumber,
						CompanyName:    info.CompanyName,
						PlanName:       info.PlanName,
						CoverageAmount: info.CoverageAmount,
						Premium:        info.Premium,
						ExpireDate:     timestamppb.New(info.ExpireDate),
						LinkedAssetId:  info.LinkedAssetID,
					},
				}
			}
		}
	}

	return res
}
