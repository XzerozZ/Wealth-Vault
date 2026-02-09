package domain

type InsuranceExpiringEvent struct {
	UserID        string `json:"user_id"`
	InsuranceID   string `json:"insurance_id"`
	InsuranceName string `json:"insurance_name"`
	DaysLeft      int    `json:"days_left"`
	ExpDate       string `json:"exp_date"`
}
