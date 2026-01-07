package domain

type CreateUser struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}
