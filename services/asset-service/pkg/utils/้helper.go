package utils

import (
	"encoding/json"
	"fmt"
	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
)

type StorageDeleter interface {
	Delete(url string) error
}

func DeleteFilesAsync(storage StorageDeleter, fileURLs []string) {
	if len(fileURLs) == 0 {
		return
	}

	go func(urls []string) {
		for _, url := range urls {
			if err := storage.Delete(url); err != nil {
				fmt.Printf("⚠️ [AsyncDelete] Failed to delete file %s: %v\n", url, err)
			}
		}
	}(fileURLs)
}

func MapProtoDetailToJSON(detail interface{}) ([]byte, error) {
	if detail == nil {
		return nil, nil
	}

	switch v := detail.(type) {
	case *pb.CreateAssetRequest_BankDetail:
		return MapBankToJSON(v.BankDetail)
	case *pb.UpdateAssetRequest_BankDetail:
		return MapBankToJSON(v.BankDetail)

	case *pb.CreateAssetRequest_InvestmentDetail:
		return MapInvestmentToJSON(v.InvestmentDetail)
	case *pb.UpdateAssetRequest_InvestmentDetail:
		return MapInvestmentToJSON(v.InvestmentDetail)

	case *pb.CreateAssetRequest_RealEstateDetail:
		return MapRealEstateToJSON(v.RealEstateDetail)
	case *pb.UpdateAssetRequest_RealEstateDetail:
		return MapRealEstateToJSON(v.RealEstateDetail)

	case *pb.CreateAssetRequest_InsuranceDetail:
		return MapInsuranceToJSON(v.InsuranceDetail)
	case *pb.UpdateAssetRequest_InsuranceDetail:
		return MapInsuranceToJSON(v.InsuranceDetail)

	default:
		return nil, nil
	}
}

func MapBankToJSON(src *pb.BankDetail) ([]byte, error) {
	if src == nil {
		return nil, nil
	}

	d := domain.BankDetail{
		BankName:      src.BankName,
		AccountNumber: src.AccountNumber,
		AccountType:   int(src.AccountType),
	}

	return json.Marshal(d)
}

func MapInvestmentToJSON(src *pb.InvestmentDetail) ([]byte, error) {
	if src == nil {
		return nil, nil
	}
	d := domain.InvestmentDetail{
		SubType:   int(src.SubType),
		Symbol:    src.Symbol,
		Broker:    src.BrokerName,
		Quantity:  src.Quantity,
		CostPrice: src.CostPrice,
	}

	return json.Marshal(d)
}

func MapRealEstateToJSON(src *pb.RealEstateDetail) ([]byte, error) {
	if src == nil {
		return nil, nil
	}

	loc := domain.Location{}
	if src.Location != nil {
		loc = domain.Location{
			Address:     src.Location.Address,
			SubDistrict: src.Location.SubDistrict,
			District:    src.Location.District,
			Province:    src.Location.Province,
			PostalCode:  src.Location.PostalCode,
		}
	}
	d := domain.RealEstateDetail{
		PropertyType: int(src.PropertyType),
		DeedNumber:   src.DeedNumber,
		AreaSqm:      src.AreaSqm,
		Location:     loc,
	}

	return json.Marshal(d)
}

func MapInsuranceToJSON(src *pb.InsuranceDetail) ([]byte, error) {
	if src == nil {
		return nil, nil
	}

	d := domain.InsuranceDetail{
		SubType:        int(src.SubType),
		PolicyNumber:   src.PolicyNumber,
		CompanyName:    src.CompanyName,
		PlanName:       src.PlanName,
		CoverageAmount: src.CoverageAmount,
		Premium:        src.Premium,
		ExpireDate:     src.ExpireDate.AsTime(),
		LinkedAssetID:  src.LinkedAssetId,
	}

	return json.Marshal(d)
}
