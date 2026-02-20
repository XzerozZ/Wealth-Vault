package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
	"wealth-vault/auth-service/internal/domain"
	authPb "wealth-vault/auth-service/pkg/pb/proto/auth"
	"wealth-vault/auth-service/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

func (u *AuthUsecase) ForgotPassword(ctx context.Context, req *authPb.ForgotPasswordRequest) (*authPb.ForgotPasswordResponse, error) {
	account, err := u.authRepo.FindByEmailAndProvider(ctx, req.Email, ProviderLocal)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if account.Provider != ProviderLocal {
		return nil, errors.New("cannot reset password for social login account")
	}

	otpCode, err := utils.GenerateOTP(OTPLength)
	if err != nil {
		return nil, fmt.Errorf("failed to generate OTP: %w", err)
	}

	otp := &domain.AuthOTP{
		UserID:    account.UserID,
		OTP:       otpCode,
		ExpiredAt: time.Now().Add(OTPTTL),
	}

	otpEmail := domain.SendEmailRequest{
		ToEmail:   account.Email,
		OTP:       otpCode,
		ExpiredAt: "5 นาที",
	}

	if err = u.authRepo.SaveOTP(ctx, otp); err != nil {
		return nil, fmt.Errorf("failed to save OTP: %w", err)
	}

	go func() {
		bgCtx := context.Background()
		if err := u.mailclient.SendOTP(bgCtx, otpEmail); err != nil {
			log.Printf("failed to send OTP email to %s: %v", account.Email, err)
		}
	}()

	return &authPb.ForgotPasswordResponse{
		Success: true,
	}, nil
}

func (u *AuthUsecase) VerifyForgotPasswordOTP(ctx context.Context, req *authPb.VerifyOTPRequest) (*authPb.VerifyOTPResponse, error) {
	account, err := u.authRepo.FindByEmailAndProvider(ctx, req.Email, ProviderLocal)
	if err != nil {
		return nil, errors.New("user not found")
	}

	_, err = u.authRepo.GetValidOTP(ctx, account.UserID, req.Otp)
	if err != nil {
		return nil, errors.New("invalid or expired OTP")
	}

	resetToken, err := u.token.CreateToken(account.UserID.String(), account.Email, TokenTypeReset, ResetTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to create reset token: %w", err)
	}

	if err := u.authRepo.DeleteOTP(ctx, account.ID); err != nil {
		log.Printf("Warning: failed to delete OTP for user ID %d: %v", account.ID, err)
	}

	return &authPb.VerifyOTPResponse{
		Success:    true,
		ResetToken: resetToken,
	}, nil
}

func (u *AuthUsecase) ResetPassword(ctx context.Context, req *authPb.ResetPasswordRequest) (*authPb.ResetPasswordResponse, error) {
	claims, err := u.token.VerifyToken(req.ResetToken)
	if err != nil {
		return nil, fmt.Errorf("invalid reset token: %w", err)
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != TokenTypeReset {
		return nil, fmt.Errorf("invalid token type: expected reset token")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	account, err := u.authRepo.FindByEmailAndProvider(ctx, email, ProviderLocal)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if account.Password != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(req.NewPassword)); err == nil {
			return nil, errors.New("new password cannot be the same as the old password")
		}
	}

	hashedPwd, err := HashPassword(req.NewPassword)
	if err != nil {
		return nil, err
	}

	if err := u.authRepo.UpdatePassword(ctx, account.UserID, hashedPwd); err != nil {
		return nil, fmt.Errorf("failed to update password: %w", err)
	}

	return &authPb.ResetPasswordResponse{
		Success: true,
	}, nil
}

func (u *AuthUsecase) DeleteExpiredOTPs(ctx context.Context) error {
	if err := u.authRepo.DeleteExpiredOTPs(ctx); err != nil {
		return fmt.Errorf("failed to cleanup OTPs: %w", err)
	}
	return nil
}
