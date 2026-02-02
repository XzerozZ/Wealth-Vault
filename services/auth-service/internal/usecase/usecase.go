package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"wealth-vault/auth-service/internal/domain"
	repo "wealth-vault/auth-service/internal/repository/interface"
	authPb "wealth-vault/auth-service/pkg/pb/proto/auth"
	pb "wealth-vault/auth-service/pkg/pb/proto/user"
	"wealth-vault/auth-service/pkg/token"
	"wealth-vault/auth-service/pkg/utils"
	"wealth-vault/auth-service/pkg/utils/mail"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	authRepo   repo.AuthRepository
	userClient pb.UserServiceClient
	token      token.Generator
	mailclient mail.NotificationClient
}

func NewAuthUsecase(
	r repo.AuthRepository,
	userClient pb.UserServiceClient,
	t token.Generator,
	mail mail.NotificationClient,
) AuthUsecase {
	return AuthUsecase{
		authRepo:   r,
		userClient: userClient,
		token:      t,
		mailclient: mail,
	}
}

func (u *AuthUsecase) Register(ctx context.Context, input *domain.RegisterInput) (*domain.AuthOutput, error) {
	existingUser, err := u.authRepo.FindByEmail(ctx, input.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already registered")
	}

	username := input.Username
	if username == "" {
		username = strings.Split(input.Email, "@")[0]
	}

	userRes, err := u.userClient.CreateUser(ctx, &pb.CreateUserRequest{
		Email:    input.Email,
		Username: username,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user profile: %v", err)
	}

	hashedPwd := ""
	if input.Provider == "local" {
		bytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		hashedPwd = string(bytes)

	}

	userID, err := uuid.Parse(userRes.User.Id)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	newAuth := &domain.AuthAccount{
		UserID:            userID,
		Email:             input.Email,
		Password:          hashedPwd,
		Provider:          input.Provider,
		ProviderAccountID: input.Email,
		IsEmailVerified:   false,
	}

	if err := u.authRepo.Register(ctx, newAuth); err != nil {
		return nil, fmt.Errorf("failed to register auth account: %w", err)
	}

	return &domain.AuthOutput{
		UserID:       userRes.User.Id,
		AccessToken:  "",
		RefreshToken: "",
	}, nil
}

func (u *AuthUsecase) Login(ctx context.Context, input *domain.LoginInput) (*domain.AuthOutput, error) {
	existingUser, err := u.authRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(input.Password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	accessToken, err := u.token.CreateToken(existingUser.UserID.String(), existingUser.Email, "access", 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	refreshToken, err := u.token.CreateToken(existingUser.UserID.String(), existingUser.Email, "refresh", 7*24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	session := &domain.AuthSession{
		UserID:           existingUser.UserID,
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ExpiresAt:        time.Now().Add(time.Hour * 1),
		RefreshExpiresAt: time.Now().Add(time.Hour * 24 * 7),
		Revoked:          false,
	}

	if err := u.authRepo.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &domain.AuthOutput{
		UserID:       existingUser.UserID.String(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *AuthUsecase) RefreshToken(ctx context.Context, refreshTokenStr string) (*domain.AuthOutput, error) {
	claims, err := u.token.VerifyToken(refreshTokenStr)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return nil, fmt.Errorf("invalid token type")
	}

	session, err := u.authRepo.GetSessionByRefreshToken(ctx, refreshTokenStr)
	if err != nil {
		return nil, fmt.Errorf("invalid or revoked refresh token")
	}

	if time.Now().After(session.RefreshExpiresAt) {
		return nil, fmt.Errorf("refresh token expired")
	}

	if err := u.authRepo.RevokeSession(ctx, refreshTokenStr); err != nil {
		return nil, err
	}

	user, err := u.authRepo.FindByID(ctx, session.UserID.String())
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	newAccess, err := u.token.CreateToken(user.UserID.String(), user.Email, "access", 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	newRefresh, err := u.token.CreateToken(user.UserID.String(), user.Email, "refresh", 7*24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	newSession := &domain.AuthSession{
		UserID:           session.UserID,
		AccessToken:      newAccess,
		RefreshToken:     newRefresh,
		ExpiresAt:        time.Now().Add(time.Hour * 1),
		RefreshExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}

	if err = u.authRepo.CreateSession(ctx, newSession); err != nil {
		return nil, err
	}

	return &domain.AuthOutput{
		UserID:       session.UserID.String(),
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	}, nil
}

func (u *AuthUsecase) CleanupSessions(ctx context.Context) error {
	if err := u.authRepo.DeleteExpiredSessions(ctx); err != nil {
		return fmt.Errorf("failed to cleanup sessions: %w", err)
	}

	return nil
}

func (u *AuthUsecase) DeleteExpiredOTPs(ctx context.Context) error {
	if err := u.authRepo.DeleteExpiredOTPs(ctx); err != nil {
		return fmt.Errorf("failed to cleanup sessions: %w", err)
	}

	return nil
}

func (u *AuthUsecase) ForgotPassword(ctx context.Context, req *authPb.ForgotPasswordRequest) (*authPb.ForgotPasswordResponse, error) {
	account, err := u.authRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if account.Provider != "local" {
		return nil, errors.New("cannot reset password for social login account")
	}

	otpCode, err := utils.GenerateOTP(6)
	if err != nil {
		return nil, err
	}

	otp := &domain.AuthOTP{
		UserID:    account.UserID,
		OTP:       otpCode,
		ExpiredAt: time.Now().Add(5 * time.Minute),
	}

	otpEmail := domain.SendEmailRequest{
		ToEmail:   account.Email,
		OTP:       otpCode,
		ExpiredAt: "5 นาที",
	}

	err = u.authRepo.SaveOTP(ctx, otp)
	if err != nil {
		return nil, err
	}

	go u.mailclient.SendOTP(ctx, otpEmail)

	return &authPb.ForgotPasswordResponse{
		Success: true,
	}, nil
}

func (u *AuthUsecase) VerifyForgotPasswordOTP(ctx context.Context, req *authPb.VerifyOTPRequest) (*authPb.VerifyOTPResponse, error) {
	account, err := u.authRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	_, err = u.authRepo.GetValidOTP(ctx, account.UserID, req.Otp)
	if err != nil {
		return nil, errors.New("invalid or expired OTP")
	}

	resetToken, err := u.token.CreateToken(account.UserID.String(), account.Email, "reset", 5*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to create reset token: %w", err)
	}

	_ = u.authRepo.DeleteOTP(ctx, account.ID)

	return &authPb.VerifyOTPResponse{
		Success:    true,
		ResetToken: resetToken,
	}, nil
}

func (u *AuthUsecase) ResetPassword(ctx context.Context, req *authPb.ResetPasswordRequest) (*authPb.ResetPasswordResponse, error) {
	claims, err := u.token.VerifyToken(req.ResetToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "reset" {
		return nil, fmt.Errorf("invalid token type: expected reset token")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	account, err := u.authRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if account.Password != "" {
		err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(req.NewPassword))
		if err == nil {
			return nil, errors.New("new password cannot be the same as the old password")
		}
	}

	hashedPwd := ""
	bytes, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	hashedPwd = string(bytes)

	if err := u.authRepo.UpdatePassword(ctx, account.UserID, hashedPwd); err != nil {
		return nil, err
	}

	return &authPb.ResetPasswordResponse{
		Success: true,
	}, nil
}
