package mapper

import (
	"strconv"
	"strings"
	"time"
	"wealth-vault/api-gateway/internal/domain"
	pb "wealth-vault/api-gateway/pkg/pb/proto/asset"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
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
			PropertyType:   int(v.RealEstateDetail.PropertyType),
			DeedNumber:     v.RealEstateDetail.DeedNumber,
			AreaSqm:        v.RealEstateDetail.AreaSqm,
			Location:       loc,
			LinkedAssetIDs: v.RealEstateDetail.LinkedAssetIds,
		}

	case *pb.Asset_InsuranceDetail:
		return domain.InsuranceDetail{
			SubType:        int(v.InsuranceDetail.SubType),
			PolicyNumber:   v.InsuranceDetail.PolicyNumber,
			CompanyName:    v.InsuranceDetail.CompanyName,
			PlanName:       v.InsuranceDetail.PlanName,
			CoverageAmount: v.InsuranceDetail.CoverageAmount,
			Premium:        v.InsuranceDetail.Premium,
			ExpireDate:     v.InsuranceDetail.ExpireDate.AsTime(),
			LinkedAssetID:  v.InsuranceDetail.LinkedAssetId,
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

	case pb.AssetType_ASSET_TYPE_REAL_ESTATE:
		return &pb.CreateAssetRequest_RealEstateDetail{
			RealEstateDetail: MapRealEstateToProto(c),
		}

	case pb.AssetType_ASSET_TYPE_INSURANCE:
		return &pb.CreateAssetRequest_InsuranceDetail{
			InsuranceDetail: MapInsuranceToProto(c),
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

	case pb.AssetType_ASSET_TYPE_REAL_ESTATE:
		if c.FormValue("location_detail_address") != "" || c.FormValue("real_estate_detail_area_sqm") != "" {
			return &pb.UpdateAssetRequest_RealEstateDetail{
				RealEstateDetail: MapRealEstateToProto(c),
			}
		}

	case pb.AssetType_ASSET_TYPE_INSURANCE:
		if c.FormValue("insurance_detail_policy_number") != "" {
			return &pb.UpdateAssetRequest_InsuranceDetail{
				InsuranceDetail: MapInsuranceToProto(c),
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

func MapRealEstateToProto(c *fiber.Ctx) *pb.RealEstateDetail {
	var loc *pb.Location
	if c.FormValue("location_detail_address") != "" {
		loc = &pb.Location{
			Address:     c.FormValue("location_detail_address"),
			SubDistrict: c.FormValue("location_detail_sub_district"),
			District:    c.FormValue("location_detail_district"),
			Province:    c.FormValue("location_detail_province"),
			PostalCode:  c.FormValue("location_detail_postal_code"),
		}
	}

	form, err := c.MultipartForm()
	var linkedIDs []string
	if err == nil && form.Value != nil {
		if values, ok := form.Value["real_estate_detail_linked_asset_ids"]; ok {
			linkedIDs = values
		}
	}

	area, _ := strconv.ParseFloat(c.FormValue("real_estate_detail_area_sqm"), 64)

	return &pb.RealEstateDetail{
		PropertyType:   pb.RealEstateType(SafeMapEnum(pb.RealEstateType_value, c.FormValue("real_estate_detail__type"), "REAL_ESTATE_TYPE_")),
		DeedNumber:     c.FormValue("real_estate_detail_deed_number"),
		AreaSqm:        area,
		Location:       loc,
		LinkedAssetIds: linkedIDs,
	}
}

func MapInsuranceToProto(c *fiber.Ctx) *pb.InsuranceDetail {
	coverage, _ := strconv.ParseFloat(c.FormValue("insurance_detail_coverage_amount"), 64)
	premium, _ := strconv.ParseFloat(c.FormValue("insurance_detail_premium"), 64)

	dateStr := c.FormValue("insurance_detail_expire_date")
	var expireDatePb *timestamppb.Timestamp
	if dateStr != "" {
		if t, err := time.Parse("2006-01-02", dateStr); err == nil {
			expireDatePb = timestamppb.New(t)
		}
	}

	return &pb.InsuranceDetail{
		PolicyNumber:   c.FormValue("insurance_detail_policy_number"),
		CompanyName:    c.FormValue("insurance_detail_company_name"),
		PlanName:       c.FormValue("insurance_detail_plan_name"),
		CoverageAmount: coverage,
		Premium:        premium,
		ExpireDate:     expireDatePb,
		SubType:        pb.InsuranceType(SafeMapEnum(pb.InsuranceType_value, c.FormValue("insurance_detail_insurance_type"), "INSURANCE_TYPE_")),
		LinkedAssetId:  c.FormValue("insurance_detail_linked_asset_id"),
	}
}
