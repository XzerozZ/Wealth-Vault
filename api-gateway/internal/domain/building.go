package domain

import "time"

type Building struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	Type        string     `json:"type"`
	Name        string     `json:"name"`
	Area        float64    `json:"area"`
	Amount      float64    `json:"amount"`
	Description string     `json:"description,omitempty"`
	Location    Location   `json:"location"`
	Files       []FileInfo `json:"files,omitempty"`
	Ref         []RefInfo  `json:"ref,omitempty"`
	Ins         []InsInfo  `json:"ins,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateBuildingRequest struct {
	Type         string          `json:"type" form:"type"`
	Name         string          `json:"name" form:"name"`
	Area         string          `json:"area" form:"area"`
	Amount       string          `json:"amount" form:"amount"`
	Description  string          `json:"desc" form:"description"`
	Location     LocationRequest `json:"location" form:"location"`
	ReferenceIDs []string        `json:"reference_ids" form:"reference_ids"`
	InsIDs       []string        `json:"ins_ids" form:"ins_ids"`
}

type UpdateBuildingRequest struct {
	Type               string          `json:"type"    form:"type" mask:"type"`
	Name               string          `json:"name"    form:"name"    mask:"name"`
	Area               string          `json:"area" form:"area" mask:"area"`
	Amount             string          `json:"amount" form:"amount"  mask:"amount"`
	Description        string          `json:"description"    form:"description"  mask:"description"`
	Location           LocationRequest `json:"location" form:"location"`
	ReferenceIDs       []string        `json:"reference_ids" form:"reference_ids"`
	DeleteReferenceIDs []string        `json:"delete_reference_ids" form:"delete_reference_ids"`
	InsIDs             []string        `json:"ins_ids" form:"ins_ids"`
	DeleteInsIDs       []string        `json:"delete_ins_ids" form:"delete_ins_ids"`
	DeleteFileIDs      []string        `json:"delete_file_ids" form:"delete_file_ids"`
}
