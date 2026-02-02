package domain

import "time"

type ShareTarget struct {
	ID      string `json:"id" form:"id"`
	ShareAt string `json:"share_at" form:"share_at"`
}

type ShareItemRequest struct {
	ItemIDs   []string      `json:"item_ids" form:"item_ids"`
	ItemTypes []string      `json:"item_types" form:"item_types"`
	Groups    []ShareTarget `json:"groups" form:"groups"`
	Friends   []ShareTarget `json:"friends" form:"friends"`
	Emails    []ShareTarget `json:"emails" form:"emails"`
}

type GroupItemResponse struct {
	GroupItemID string      `json:"group_item_id"`
	SharedBy    string      `json:"shared_by"`
	SharedAt    time.Time   `json:"shared_at"`
	Type        string      `json:"type"`
	AssetDetail interface{} `json:"asset_detail"`
}

type FriendItemResponse struct {
	ItemID      string      `json:"shared_item_id"`
	SharedBy    string      `json:"shared_by"`
	SharedAt    time.Time   `json:"shared_at"`
	Type        string      `json:"type"`
	AssetDetail interface{} `json:"asset_detail"`
}

type BuildingDetail struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Amount       float64 `json:"amount"`
	LocationText string  `json:"location_text"`
	TypeName     string  `json:"type"`
}

type AccountDetail struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	BankName      string  `json:"bank_name"`
	AccountNumber string  `json:"account_number"`
	Amount        float64 `json:"amount"`
}

type LandDetail struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	DeedNum      string  `json:"deed_num"`
	Area         float64 `json:"area"`
	Amount       float64 `json:"amount"`
	LocationText string  `json:"location"`
}

type CashDetail struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}

type InsuranceDetail struct {
	ID             string  `json:"id"`
	CompanyName    string  `json:"company_name"`
	TypeName       string  `json:"type"`
	PolNum         string  `json:"pol_num"`
	CoverageAmount float64 `json:"coverage_amount"`
	ExpDateText    string  `json:"exp_date_text"`
}

type InvestmentDetail struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	TypeName string `json:"type_name"`
}

type LiabilityDetail struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Creditor  string  `json:"creditor"`
	Principal float64 `json:"principal"`
	TypeName  string  `json:"type"`
}
