package domain

import "time"

type Cash struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Amount      float64    `json:"amount"`
	Description string     `json:"desc"`
	CreatedBy   string     `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Files       []FileInfo `json:"files"`
}

type CreateCash struct {
	Name        string `json:"name"   form:"name"`
	Amount      string `json:"amount" form:"amount"     validate:"required"`
	Description string `json:"desc"   form:"description"`
	CreatedBy   string `json:"created_by"`
}

type UpdateCashRequest struct {
	Name          string   `json:"name"    form:"name"    mask:"name"`
	Amount        string   `json:"amount"  form:"amount"  mask:"value"`
	Description   string   `json:"desc"    form:"description"    mask:"description"`
	DeleteFileIDs []string `json:"delete_file_ids" form:"delete_file_ids"`
}
