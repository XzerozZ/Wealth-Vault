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

func (h *AuthGRPCHandler) GetProviderAccount(ctx context.Context, req *pb.GetProviderAccountRequest) (*pb.GetProviderAccountsResponse, error) {
	supportedProviders := []string{"local", "google", "line"}
	existingAccounts, err := h.usecase.GetAllProviderAccounts(ctx, req.UserId)

	linkedMap := make(map[string]string)
	if err == nil {
		for _, acc := range existingAccounts {
			linkedMap[acc.Provider] = acc.ProviderAccountID
		}
	}

	var pbAccounts []*pb.ProviderAccount
	for _, provider := range supportedProviders {
		accountID, exists := linkedMap[provider]

		pbAccounts = append(pbAccounts, &pb.ProviderAccount{
			Provider:          provider,
			IsLinked:          exists,
			ProviderAccountId: accountID,
		})
	}

	return &pb.GetProviderAccountsResponse{
		Accounts: pbAccounts,
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

func (h *AuthGRPCHandler) LinkLineAccount(ctx context.Context, req *pb.LinkLineAccountRequest) (*pb.LinkAccountResponse, error) {
	err := h.usecase.LinkLineAccount(ctx, req.UserId, req.LineUserId)
	if err != nil {
		return &pb.LinkAccountResponse{
			Success: false,
		}, err
	}

	return &pb.LinkAccountResponse{
		Success: true,
	}, nil
}
