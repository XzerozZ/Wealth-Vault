package usecase

import (
	"context"
	"wealth-vault/auth-service/internal/domain"
	authPb "wealth-vault/auth-service/pkg/pb/proto/auth"
)

type AuthUsecase interface {
	Register(ctx context.Context, input *domain.RegisterInput) (*domain.AuthOutput, error)
	Login(ctx context.Context, input *domain.LoginInput) (*domain.AuthOutput, error)
	LoginWithGoogle(ctx context.Context, googleIDToken string) (*domain.AuthOutput, error)
	LinkLineAccount(ctx context.Context, userID string, lineUserID string) error
	RefreshToken(ctx context.Context, refreshToken string) (*domain.AuthOutput, error)
	ForgotPassword(ctx context.Context, req *authPb.ForgotPasswordRequest) (*authPb.ForgotPasswordResponse, error)
	VerifyForgotPasswordOTP(ctx context.Context, req *authPb.VerifyOTPRequest) (*authPb.VerifyOTPResponse, error)
	ResetPassword(ctx context.Context, req *authPb.ResetPasswordRequest) (*authPb.ResetPasswordResponse, error)
	CleanupSessions(ctx context.Context) error
	DeleteExpiredOTPs(ctx context.Context) error

	GetAllProviderAccounts(ctx context.Context, userID string) ([]domain.AuthAccount, error)
}
