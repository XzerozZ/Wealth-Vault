package mapper

import (
	"strconv"
	"strings"
	"wealth-vault/api-gateway/internal/domain"
	pb "wealth-vault/api-gateway/pkg/pb/proto/asset"

	"github.com/gofiber/fiber/v2"
)

func SafeMapEnum(valueMap map[string]int32, key string, prefix string) int32 {
	key = strings.ToUpper(strings.TrimSpace(key))
	if v, ok := valueMap[key]; ok {
		return v
	}

	if v, ok := valueMap[prefix+key]; ok {
		return v
	}

	return 0
}

func ToAssetDomain(p *pb.Asset) *domain.Asset {
	if p == nil {
		return nil
	}

	var domainFiles []domain.FileInfo
	if len(p.Files) > 0 {
		domainFiles = make([]domain.FileInfo, len(p.Files))
		for i, f := range p.Files {
			domainFiles[i] = domain.FileInfo{
				ID:       f.Id,
				URL:      f.Url,
				FileType: f.FileType,
			}
		}
	}

	detailObj := ExtractDetail(p)

	return &domain.Asset{
		ID:                   p.Id,
		Name:                 p.Name,
		Amount:               p.Amount,
		IsIncludedInNetWorth: p.IsIncludedInNetWorth,
		Description:          p.Description,
		Type:                 p.Type.String(),
		Details:              detailObj,
		UserID:               p.CreatedBy,
		Files:                domainFiles,
		CreatedAt:            p.CreatedAt.AsTime(),
		UpdatedAt:            p.UpdatedAt.AsTime(),
	}
}

func ToAssetList(pbList []*pb.Asset) []domain.Asset {
	if pbList == nil {
		return []domain.Asset{}
	}

	entities := make([]domain.Asset, 0, len(pbList))

	for _, pbItem := range pbList {
		if item := ToAssetDomain(pbItem); item != nil {
			entities = append(entities, *item)
		}
	}

	return entities
}

func ExtractDetail(p *pb.Asset) interface{} {
	switch v := p.Detail.(type) {

	case *pb.Asset_BankDetail:
		return domain.BankDetail{
			BankName:      v.BankDetail.BankName,
			AccountNumber: v.BankDetail.AccountNumber,
			AccountType:   v.BankDetail.AccountType.String(),
		}

	case *pb.Asset_InvestmentDetail:
		return domain.InvestmentDetail{
			SubType:   int(v.InvestmentDetail.SubType),
			Symbol:    v.InvestmentDetail.Symbol,
			Broker:    v.InvestmentDetail.BrokerName,
			Quantity:  v.InvestmentDetail.Quantity,
			CostPrice: v.InvestmentDetail.CostPrice,
		}

	case *pb.Asset_RealEstateDetail:
		loc := domain.Location{}
		if v.RealEstateDetail.Location != nil {
			loc = domain.Location{
				Address:     v.RealEstateDetail.Location.Address,
				SubDistrict: v.RealEstateDetail.Location.SubDistrict,
				District:    v.RealEstateDetail.Location.District,
				Province:    v.RealEstateDetail.Location.Province,
				PostalCode:  v.RealEstateDetail.Location.PostalCode,
			}
		}
		return domain.RealEstateDetail{
			PropertyType: int(v.RealEstateDetail.PropertyType),
			DeedNumber:   v.RealEstateDetail.DeedNumber,
			AreaSqm:      v.RealEstateDetail.AreaSqm,
			Location:     loc,
		}

	case *pb.Asset_InsuranceDetail:
		return domain.InsuranceDetail{
			SubType:        int(v.InsuranceDetail.SubType),
			PolicyNumber:   v.InsuranceDetail.PolicyNumber,
			CompanyName:    v.InsuranceDetail.CompanyName,
			PlanName:       v.InsuranceDetail.PlanName,
			CoverageAmount: v.InsuranceDetail.CoverageAmount,
			Premium:        v.InsuranceDetail.Premium,
		}

	default:
		return nil
	}
}

func BuildCreateDetail(c *fiber.Ctx, assetType pb.AssetType) interface{} {
	switch assetType {
	case pb.AssetType_ASSET_TYPE_BANK:
		return &pb.CreateAssetRequest_BankDetail{
			BankDetail: MapBankToProto(c),
		}
	case pb.AssetType_ASSET_TYPE_INVESTMENT:
		return &pb.CreateAssetRequest_InvestmentDetail{
			InvestmentDetail: MapInvestmentToProto(c),
		}
	default:
		return nil
	}
}

func BuildUpdateDetail(c *fiber.Ctx, assetType pb.AssetType) interface{} {
	switch assetType {
	case pb.AssetType_ASSET_TYPE_BANK:
		if c.FormValue("bank_detail_bank_name") != "" || c.FormValue("bank_detail_account_number]") != "" {
			return &pb.UpdateAssetRequest_BankDetail{
				BankDetail: MapBankToProto(c),
			}
		}
	case pb.AssetType_ASSET_TYPE_INVESTMENT:
		if c.FormValue("investment_detail_symbol") != "" {
			return &pb.UpdateAssetRequest_InvestmentDetail{
				InvestmentDetail: MapInvestmentToProto(c),
			}
		}
	}

	return nil
}

func MapBankToProto(c *fiber.Ctx) *pb.BankDetail {
	return &pb.BankDetail{
		BankName:      c.FormValue("bank_detail_bank_name"),
		AccountNumber: c.FormValue("bank_detail_account_number"),
		AccountType:   pb.BankAccountType(SafeMapEnum(pb.BankAccountType_value, c.FormValue("bank_detail_account_type"), "BANK_ACC_TYPE_")),
	}
}

func MapInvestmentToProto(c *fiber.Ctx) *pb.InvestmentDetail {
	qty, _ := strconv.ParseFloat(c.FormValue("investment_detail_quantity"), 64)
	cost, _ := strconv.ParseFloat(c.FormValue("investment_detail_cost_price"), 64)

	return &pb.InvestmentDetail{
		SubType:    pb.InvestmentType(SafeMapEnum(pb.InvestmentType_value, c.FormValue("investment_detail_invest_type"), "INVEST_TYPE_")),
		Symbol:     c.FormValue("investment_detail_symbol"),
		BrokerName: c.FormValue("investment_detail_broker"),
		Quantity:   qty,
		CostPrice:  cost,
	}
}

func MapRealEstateToProto(d *domain.RealEstateDetailDTO) *pb.RealEstateDetail {
	var loc *pb.Location
	if d.Location != nil {
		loc = &pb.Location{
			Address:     d.Location.Address,
			SubDistrict: d.Location.SubDistrict,
			District:    d.Location.District,
			Province:    d.Location.Province,
			PostalCode:  d.Location.PostalCode,
		}
	}
	return &pb.RealEstateDetail{
		PropertyType: pb.RealEstateType(pb.RealEstateType_value[d.PropertyType]),
		DeedNumber:   d.DeedNumber,
		AreaSqm:      d.AreaSqm,
		Location:     loc,
	}
}

func MapInsuranceToProto(d *domain.InsuranceDetailDTO) *pb.InsuranceDetail {
	return &pb.InsuranceDetail{
		PolicyNumber:   d.PolicyNumber,
		CompanyName:    d.CompanyName,
		PlanName:       d.PlanName,
		CoverageAmount: d.CoverageAmount,
		Premium:        d.Premium,
		SubType:        pb.InsuranceType(pb.InsuranceType_value[d.SubType]),
	}
}
