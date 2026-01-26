package domain

import "time"

type Cash struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	Name        string     `json:"name"`
	Amount      float64    `json:"amount"`
	Description string     `json:"description,omitempty"`
	Files       []FileInfo `json:"files,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateCashRequest struct {
	Name        string `json:"name" form:"name"`
	Amount      string `json:"amount" form:"amount"`
	Description string `json:"desc"    form:"description"`
}

type UpdateCashRequest struct {
	Name          string   `json:"name"    form:"name"    mask:"name"`
	Amount        string   `json:"amount" form:"amount"  mask:"amount"`
	Description   string   `json:"description"    form:"description"    mask:"description"`
	DeleteFileIDs []string `json:"delete_file_ids" form:"delete_file_ids"`
}
