package domain

type SendEmailRequest struct {
	ToEmail    string
	AssetName  string
	AssetType  string
	ItemDetail map[string]string
}
