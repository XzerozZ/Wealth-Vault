package grpc

import (
	"context"
	"wealth-vault/auth-service/internal/domain"
	usecase "wealth-vault/auth-service/internal/usecase/interface"
	pb "wealth-vault/auth-service/pkg/pb/proto/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthGRPCHandler struct {
	pb.UnimplementedAuthServiceServer
	usecase usecase.AuthUsecase
}

func NewAuthGRPCHandler(u usecase.AuthUsecase) *AuthGRPCHandler {
	return &AuthGRPCHandler{usecase: u}
}

func (h *AuthGRPCHandler) Register(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	input := &domain.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
		Provider: "local",
	}

	output, err := h.usecase.Register(ctx, input)
	if err != nil {
		if err.Error() == "email already registered" {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "registration failed: %v", err)
	}

	return &pb.AuthResponse{
		UserId:       output.UserID,
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	}, nil
}

func (h *AuthGRPCHandler) Login(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	input := &domain.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := h.usecase.Login(ctx, input)
	if err != nil {
		if err.Error() == "invalid email" || err.Error() == "invalid password" {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "login failed: %v", err)
	}

	return &pb.AuthResponse{
		Success:      true,
		UserId:       output.UserID,
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	}, nil
}

func (h *AuthGRPCHandler) LoginGoogle(ctx context.Context, req *pb.GoogleRequest) (*pb.AuthResponse, error) {
	output, err := h.usecase.LoginWithGoogle(ctx, req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "login failed: %v", err)
	}

	return &pb.AuthResponse{
		Success:      true,
		UserId:       output.UserID,
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	}, nil
}

func (h *AuthGRPCHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.AuthResponse, error) {
	output, err := h.usecase.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "refresh token failed: %v", err)
	}

	return &pb.AuthResponse{
		Success:      true,
		UserId:       output.UserID,
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	}, nil
}

func (h *AuthGRPCHandler) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*pb.ForgotPasswordResponse, error) {
	res, err := h.usecase.ForgotPassword(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AuthGRPCHandler) VerifyForgotPasswordOTP(ctx context.Context, req *pb.VerifyOTPRequest) (*pb.VerifyOTPResponse, error) {
	res, err := h.usecase.VerifyForgotPasswordOTP(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AuthGRPCHandler) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	res, err := h.usecase.ResetPassword(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
