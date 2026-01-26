package domain

import "time"

type BankDetail struct {
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"`
	AccountType   int    `json:"account_type"`
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

type InsuranceDetail struct {
	SubType        int       `json:"sub_type"`
	PolicyNumber   string    `json:"policy_number"`
	CompanyName    string    `json:"company_name"`
	PlanName       string    `json:"plan_name"`
	CoverageAmount float64   `json:"coverage_amount"`
	Premium        float64   `json:"premium"`
	ExpireDate     time.Time `json:"expire_date"`
	LinkedAssetID  string    `json:"linked_asset_id,omitempty"`
}
