package domain

type Authenticate struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
