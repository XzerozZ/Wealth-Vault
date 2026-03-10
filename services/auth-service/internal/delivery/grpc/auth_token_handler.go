package grpc

import (
	"context"
	pb "wealth-vault/auth-service/pkg/pb/proto/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
