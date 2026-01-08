package grpc

import (
	"context"
	"wealth-vault/auth-service/internal/domain"
	"wealth-vault/auth-service/internal/usecase"
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
		UserId:       output.UserID,
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	}, nil
}
