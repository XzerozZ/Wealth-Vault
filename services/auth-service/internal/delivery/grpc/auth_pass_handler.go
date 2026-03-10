package grpc

import (
	"context"
	pb "wealth-vault/auth-service/pkg/pb/proto/auth"
)

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
