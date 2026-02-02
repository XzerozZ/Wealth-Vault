package domain

import "time"

type InsInfo struct {
	ID   string `json:"ins_id"`
	Name string `json:"ins_name"`
}

type Insurance struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	Name           string     `json:"name"`
	PolicyNumber   string     `json:"policy_number"`
	Type           string     `json:"type"`
	CompanyName    string     `json:"company_name"`
	CoveragePeriod float64    `json:"coverage_period"`
	CoverageAmount float64    `json:"coverage_amount"`
	ConDate        *time.Time `json:"con_date"`
	ExpDate        *time.Time `json:"exp_date"`
	Description    string     `json:"description,omitempty"`
	Files          []FileInfo `json:"files,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type CreateInsuranceRequest struct {
	Type           string `json:"type" form:"type"`
	Name           string `json:"name" form:"name"`
	PolicyNumber   string `json:"policy_number" form:"policy_number"`
	CompanyName    string `json:"company_name" form:"company_name"`
	CoveragePeriod string `json:"coverage_period" form:"coverage_period"`
	CoverageAmount string `json:"coverage_amount" form:"coverage_amount"`
	ConDate        string `json:"con_date" form:"con_date"`
	ExpDate        string `json:"exp_date" form:"exp_date"`
	Description    string `json:"desc"    form:"description"`
}

type UpdateInsuranceRequest struct {
	Name           string   `json:"name"    form:"name"    mask:"name"`
	PolicyNumber   string   `json:"policy_number"    form:"policy_number"    mask:"policy_number"`
	CompanyName    string   `json:"company_name" form:"company_name" mask:"company_name"`
	CoveragePeriod string   `json:"coverage_period" form:"coverage_period" mask:"coverage_period"`
	CoverageAmount string   `json:"coverage_amount" form:"coverage_amount"  mask:"coverage_amount"`
	ConDate        string   `json:"con_date" form:"con_date"  mask:"con_date"`
	ExpDate        string   `json:"exp_date" form:"exp_date"  mask:"exp_date"`
	Type           string   `json:"type"    form:"type" mask:"type"`
	Description    string   `json:"description"    form:"description"    mask:"description"`
	DeleteFileIDs  []string `json:"delete_file_ids" form:"delete_file_ids"`
}
