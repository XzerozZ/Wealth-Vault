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

type FriendItemPreviewResponse struct {
	ItemID      string      `json:"item_id"`
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
	Name           string  `json:"name"`
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

type DeletedDetail struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

type SharedTargetsResponse struct {
	Groups  []SharedGroupTarget  `json:"groups"`
	Friends []SharedFriendTarget `json:"friends"`
	Emails  []SharedEmailTarget  `json:"emails"`
}

type SharedGroupTarget struct {
	GroupID     string    `json:"group_id"`
	GroupName   string    `json:"group_name"`
	GroupImage  string    `json:"group_image"`
	MemberCount int64     `json:"member_count"`
	SharedAt    time.Time `json:"shared_at"`
}

type SharedFriendTarget struct {
	FriendID     string    `json:"friend_id"`
	Username     string    `json:"username"`
	ProfileImage string    `json:"profile_image"`
	SharedAt     time.Time `json:"shared_at"`
}

type SharedEmailTarget struct {
	Email    string    `json:"email"`
	SharedAt time.Time `json:"shared_at"`
	IsSent   bool      `json:"is_sent"`
}

type AssetSelection struct {
	ID       string  `json:"id"`
	Type     string  `json:"type"`
	Name     string  `json:"name"`
	Value    float64 `json:"value"`
	IsShared bool    `json:"is_shared"`
}
