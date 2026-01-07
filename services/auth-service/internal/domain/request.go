package domain

type RegisterInput struct {
	Email      string
	Password   string
	Username   string
	Provider   string
	ProviderID string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthOutput struct {
	UserID       string
	AccessToken  string
	RefreshToken string
}
