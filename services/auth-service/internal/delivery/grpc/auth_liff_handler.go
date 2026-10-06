package grpc

import (
	"context"
	"errors"
	"wealth-vault/auth-service/internal/domain"
	line "wealth-vault/auth-service/internal/usecase"
	pb "wealth-vault/auth-service/pkg/pb/proto/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func lineError(err error) error {
	switch {
	case errors.Is(err, line.ErrLineRejected):
		return status.Error(codes.Unauthenticated, "LINE verification or challenge rejected")
	case errors.Is(err, line.ErrLineNotConfigured):
		return status.Error(codes.FailedPrecondition, "LINE login not configured")
	case errors.Is(err, line.ErrLineUnavailable):
		return status.Error(codes.Unavailable, "LINE verification unavailable")
	default:
		return status.Error(codes.Internal, "LINE authentication failed")
	}
}

func lineResponse(result *domain.LineAuthResult) *pb.LineAuthResponse {
	return &pb.LineAuthResponse{Status: result.Status, ChallengeId: result.ChallengeID, UserId: result.UserID,
		DisplayName: result.DisplayName, AccessToken: result.AccessToken, RefreshToken: result.RefreshToken,
		AccessTokenExpiresIn: result.AccessTokenExpiresIn}
}

func (h *AuthGRPCHandler) ExchangeLine(ctx context.Context, req *pb.LineExchangeRequest) (*pb.LineAuthResponse, error) {
	result, err := h.usecase.ExchangeLine(ctx, req.IdToken, req.ChannelId)
	if err != nil {
		return nil, lineError(err)
	}
	return lineResponse(result), nil
}

func (h *AuthGRPCHandler) CompleteLineChallenge(ctx context.Context, req *pb.LineChallengeRequest) (*pb.LineAuthResponse, error) {
	result, err := h.usecase.CompleteLineChallenge(ctx, domain.LineChallengeInput{
		ChallengeID: req.ChallengeId, ChannelID: req.ChannelId, Action: req.Action, Email: req.Email,
		Password: req.Password, Username: req.Username, FirstName: req.FirstName, LastName: req.LastName,
	})
	if err != nil {
		return nil, lineError(err)
	}
	return lineResponse(result), nil
}

func (h *AuthGRPCHandler) RevokeSession(ctx context.Context, req *pb.RevokeSessionRequest) (*pb.RevokeSessionResponse, error) {
	if err := h.usecase.RevokeRefreshSession(ctx, req.RefreshToken); err != nil {
		return nil, lineError(err)
	}
	return &pb.RevokeSessionResponse{Success: true}, nil
}
