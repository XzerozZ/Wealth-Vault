package domain

import "time"

type RefInfo struct {
	ID   string `json:"ref_id"`
	Name string `json:"ref_name"`
}

type Location struct {
	ID          string    `json:"location_id"`
	Address     string    `json:"address"`
	Subdistrict string    `json:"sub_district"`
	District    string    `json:"district"`
	Province    string    `json:"province"`
	PostalCode  string    `json:"postal_code"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LocationRequest struct {
	Address     string `json:"address" form:"address" mask:"address"`
	Subdistrict string `json:"sub_district" form:"sub_district" mask:"sub_district"`
	District    string `json:"district" form:"district" mask:"district"`
	Province    string `json:"province" form:"province" mask:"province"`
	PostalCode  string `json:"postal_code" form:"postal_code" mask:"postal_code"`
}
