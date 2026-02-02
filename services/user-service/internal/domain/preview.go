package domain

type AccountPreview struct {
	ID          string
	Name        string
	BankName    string
	AccountMask string
	Amount      float64
}

type CashPreview struct {
	ID     string
	Name   string
	Amount float64
}

type InvestmentPreview struct {
	ID     string
	Name   string
	Type   string
	Symbol string
}

type BuildingPreview struct {
	ID       string
	Name     string
	Type     string
	Amount   float64
	Location string
}

type LandPreview struct {
	ID     string
	Name   string
	Area   float64
	Amount float64
}

type InsurancePreview struct {
	ID             string
	Type           string
	CompanyName    string
	PolicyNumber   string
	CoverageAmount float64
	ExpDate        string
}

type LiabilityPreview struct {
	ID        string
	Name      string
	Type      string
	Creditor  string
	Principal float64
}
