package domain

type Authenticate struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type OAuth struct {
	Token string `json:"token" form:"token"`
}

type RefreshToken struct {
	RefreshToken string `json:"refreshtoken"`
}

type ForgetPassword struct {
	Email string `json:"email" form:"email"`
}

type OTP struct {
	Email string `json:"email" form:"email"`
	OTP   string `json:"otp" form:"otp"`
}

type ResetPassword struct {
	ResetToken string `json:"resettoken" form:"resettoken"`
	Password   string `json:"password" form:"password"`
}
