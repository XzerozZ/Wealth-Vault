package domain

import "time"

type Liability struct {
	ID           string     `json:"id"`
	UserID       string     `json:"user_id"`
	Type         string     `json:"type"`
	Name         string     `json:"name"`
	Creditor     string     `json:"creditor"`
	Principal    float64    `json:"principal"`
	InterestRate float64    `json:"interest_rate"`
	Description  string     `json:"description,omitempty"`
	StartAt      *time.Time `json:"started_at"`
	EndAt        *time.Time `json:"ended_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Files        []FileInfo `json:"files,omitempty"`
}

type CreateLiabilityRequest struct {
	Type         string `json:"type" form:"type"`
	Name         string `json:"name" form:"name"`
	Creditor     string `json:"creditor" form:"creditor"`
	Principal    string `json:"principal" form:"principal"`
	InterestRate string `json:"interest_rate" form:"interest_rate"`
	Description  string `json:"desc"    form:"description"`
	StartAt      string `json:"started_at" form:"started_at"`
	EndAt        string `json:"ended_at" form:"ended_at"`
}

type UpdateLiabilityRequest struct {
	Name          string   `json:"name"    form:"name"    mask:"name"`
	Creditor      string   `json:"creditor" form:"creditor" mask:"creditor"`
	Principal     string   `json:"principal" form:"principal" mask:"principal"`
	InterestRate  string   `json:"interest_rate" form:"interest_rate"  mask:"interest_rate"`
	Type          string   `json:"type"    form:"type"`
	Description   string   `json:"description"    form:"description"    mask:"description"`
	StartAt       string   `json:"started_at" form:"started_at" mask:"started_at"`
	EndAt         string   `json:"ended_at" form:"ended_at" mask:"ended_at"`
	DeleteFileIDs []string `json:"delete_file_ids" form:"delete_file_ids"`
}
