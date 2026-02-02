package domain

import "time"

type Account struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	Name        string     `json:"name"`
	BankName    string     `json:"bank_name"`
	BankAccount string     `json:"bank_account"`
	Type        string     `json:"type"`
	Amount      float64    `json:"amount"`
	Description string     `json:"description,omitempty"`
	Files       []FileInfo `json:"files,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateAccountRequest struct {
	Type        string `json:"type" form:"type"`
	Name        string `json:"name" form:"name"`
	BankName    string `json:"bank_name" form:"bank_name"`
	BankAccount string `json:"bank_account" form:"bank_account"`
	Amount      string `json:"amount" form:"amount"`
	Description string `json:"desc"    form:"description"`
}

type UpdateAccountRequest struct {
	Name          string   `json:"name"    form:"name"    mask:"name"`
	BankName      string   `json:"bank_name" form:"bank_name" mask:"bank_name"`
	BankAccount   string   `json:"bank_account" form:"bank_account" mask:"bank_account"`
	Amount        string   `json:"amount" form:"amount"  mask:"amount"`
	Type          string   `json:"type"    form:"type" mask:"type"`
	Description   string   `json:"description"    form:"description"    mask:"description"`
	DeleteFileIDs []string `json:"delete_file_ids" form:"delete_file_ids"`
}
