package domain

import "time"

type Investment struct {
	ID           string     `json:"id"`
	UserID       string     `json:"user_id"`
	Name         string     `json:"name"`
	Symbol       string     `json:"symbol"`
	Type         string     `json:"type"`
	BrokerName   string     `json:"broker_name"`
	Quantity     float64    `json:"quantity"`
	CostPerPrice float64    `json:"cost_per_price"`
	Description  string     `json:"description,omitempty"`
	Files        []FileInfo `json:"files,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type CreateInvestmentRequest struct {
	Type         string `json:"type" form:"type"`
	Name         string `json:"name" form:"name"`
	Symbol       string `json:"symbol" form:"symbol"`
	BrokerName   string `json:"broker_name" form:"broker_name"`
	Quantity     string `json:"quantity" form:"quantity"`
	CostPerPrice string `json:"cost_per_price" form:"cost_per_price"`
	Description  string `json:"desc"    form:"description"`
}

type UpdateInvestmentRequest struct {
	Name          string   `json:"name"    form:"name"    mask:"name"`
	Symbol        string   `json:"symbol" form:"symbol" mask:"symbol"`
	BrokerName    string   `json:"broker_name" form:"broker_name" mask:"broker_name"`
	Quantity      string   `json:"quantity" form:"quantity"  mask:"quantity"`
	CostPerPrice  string   `json:"cost_per_price" form:"cost_per_price"  mask:"cost_per_price"`
	Type          string   `json:"type"    form:"type" mask:"type"`
	Description   string   `json:"description"    form:"description"    mask:"description"`
	DeleteFileIDs []string `json:"delete_file_ids" form:"delete_file_ids"`
}
