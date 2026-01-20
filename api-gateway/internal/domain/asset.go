package domain

import "time"

type Asset struct {
	ID                   string  `json:"id"`
	UserID               string  `json:"user_id"`
	Type                 string  `json:"type"`
	Name                 string  `json:"name"`
	Amount               float64 `json:"amount"`
	Description          string  `json:"description,omitempty"`
	IsIncludedInNetWorth bool    `json:"is_included_in_net_worth"`

	Details interface{} `json:"details,omitempty"`

	Files     []FileInfo `json:"files,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type BankDetail struct {
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"`
	AccountType   string `json:"account_type"`
}

type InvestmentDetail struct {
	SubType   int     `json:"sub_type"`
	Symbol    string  `json:"symbol"`
	Broker    string  `json:"broker"`
	Quantity  float64 `json:"quantity"`
	CostPrice float64 `json:"cost_price"`
}

type RealEstateDetail struct {
	PropertyType   int      `json:"property_type"`
	DeedNumber     string   `json:"deed_number"`
	AreaSqm        float64  `json:"area_sqm"`
	Location       Location `json:"location"`
	LinkedAssetIDs []string `json:"linked_asset_ids"`
}

type Location struct {
	Address     string `json:"address"`
	SubDistrict string `json:"sub_district"`
	District    string `json:"district"`
	Province    string `json:"province"`
	PostalCode  string `json:"postal_code"`
}

type InsuranceDetail struct {
	SubType        int       `json:"sub_type"`
	PolicyNumber   string    `json:"policy_number"`
	CompanyName    string    `json:"company_name"`
	PlanName       string    `json:"plan_name"`
	CoverageAmount float64   `json:"coverage_amount"`
	Premium        float64   `json:"premium"`
	ExpireDate     time.Time `json:"expire_date"`
	LinkedAssetID  string    `json:"linked_asset_id"`
}

type CreateAssetRequest struct {
	Type                 string               `json:"type" form:"type"`
	Name                 string               `json:"name" form:"name"`
	Amount               string               `json:"amount" form:"amount"`
	IsIncludedInNetWorth string               `json:"is_included_in_net_worth" form:"is_included_in_net_worth"`
	Description          string               `json:"desc"    form:"description"    mask:"description"`
	BankDetail           *BankDetailDTO       `json:"bank_detail" form:"bank_detail"`
	InvestmentDetail     *InvestmentDetailDTO `json:"investment_detail" form:"investment_detail"`
	RealEstateDetail     *RealEstateDetailDTO `json:"real_estate_detail" form:"real_estate_detail"`
	InsuranceDetail      *InsuranceDetailDTO  `json:"insurance_detail" form:"insurance_detail"`
}

type UpdateAssetRequest struct {
	Name             string               `json:"name"    form:"name"    mask:"name"`
	Amount           string               `json:"amount"  form:"amount"  mask:"amount"`
	Type             string               `json:"type"    form:"type"`
	Description      string               `json:"description"    form:"description"    mask:"description"`
	BankDetail       *BankDetailDTO       `json:"bank_detail" form:"bank_detail" mask:"detail"`
	InvestmentDetail *InvestmentDetailDTO `json:"investment_detail" form:"investment_detail" mask:"detail"`
	RealEstateDetail *RealEstateDetailDTO `json:"real_estate_detail" form:"real_estate_detail" mask:"detail"`
	InsuranceDetail  *InsuranceDetailDTO  `json:"insurance_detail" form:"insurance_detail" mask:"detail"`
	DeleteFileIDs    []string             `json:"delete_file_ids" form:"delete_file_ids"`
}

type BankDetailDTO struct {
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"`
	AccountType   string `json:"account_type"`
}

type InvestmentDetailDTO struct {
	SubType   string  `json:"sub_type"`
	Symbol    string  `json:"symbol"`
	Broker    string  `json:"broker"`
	Quantity  float64 `json:"quantity"`
	CostPrice float64 `json:"cost_price"`
}

type RealEstateDetailDTO struct {
	PropertyType  string       `json:"property_type"`
	DeedNumber    string       `json:"deed_number"`
	AreaSqm       float64      `json:"area_sqm"`
	Location      *LocationDTO `json:"location"`
	LinkedAssetID []string     `json:"linked_asset_id,omitempty"`
}

type InsuranceDetailDTO struct {
	SubType        string    `json:"sub_type"`
	PolicyNumber   string    `json:"policy_number"`
	CompanyName    string    `json:"company_name"`
	PlanName       string    `json:"plan_name"`
	CoverageAmount float64   `json:"coverage_amount"`
	Premium        float64   `json:"premium"`
	ExpireDate     time.Time `json:"expire_date"`
	LinkedAssetID  string    `json:"linked_asset_id,omitempty"`
}

type LocationDTO struct {
	Address     string  `json:"address"`
	SubDistrict string  `json:"sub_district"`
	District    string  `json:"district"`
	Province    string  `json:"province"`
	PostalCode  string  `json:"postal_code"`
	Latitude    float64 `json:"latitude,omitempty"`
	Longitude   float64 `json:"longitude,omitempty"`
}
