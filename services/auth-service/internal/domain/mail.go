package domain

type SendEmailRequest struct {
	ToEmail   string
	OTP       string
	ExpiredAt string
}
