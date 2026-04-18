package usecase

import (
	"context"
	"errors"
	"log"
	"time"
	"wealth-vault/asset-service/internal/domain"
	repo "wealth-vault/asset-service/internal/repository/interface"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	userPb "wealth-vault/asset-service/pkg/pb/proto/user"
	"wealth-vault/asset-service/pkg/utils"
	helper "wealth-vault/asset-service/pkg/utils/helper"
	"wealth-vault/asset-service/pkg/utils/mapper"

	"github.com/google/uuid"
)

type InvestmentUsecase struct {
	inRepo      repo.InvestmentRepository
	assetHelper helper.AssetHelper
	userClient  userPb.UserServiceClient
}

func NewInvestmentUsecase(r repo.InvestmentRepository, ah helper.AssetHelper, uc userPb.UserServiceClient) *InvestmentUsecase {
	return &InvestmentUsecase{
		inRepo:      r,
		assetHelper: ah,
		userClient:  uc,
	}
}

func (u *InvestmentUsecase) CreateInvestment(ctx context.Context, req *pb.CreateInvestmentRequest) (*pb.InvestmentResponse, error) {
	uid, err := utils.ParseID(req.UserId)
	if err != nil {
		return nil, err
	}

	in := mapper.ToInvestmentDomain(req, uid)

	if err := u.inRepo.CreateInvestment(ctx, in); err != nil {
		return nil, err
	}

	return &pb.InvestmentResponse{
		Invest: mapper.ToInvestProto(in),
	}, nil
}

func (u *InvestmentUsecase) GetInvestment(ctx context.Context, req *pb.GetAssetRequest) (*pb.InvestmentArrayResponse, error) {
	uid, err := utils.ParseID(req.UserId)
	if err != nil {
		return nil, err
	}

	invests, err := u.inRepo.GetInvestment(ctx, uid)
	if err != nil {
		return nil, err
	}

	return &pb.InvestmentArrayResponse{
		Success: true,
		Invest:  mapper.ToInvestProtoSlice(invests),
	}, nil
}

func (u *InvestmentUsecase) GetInvestmentByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.InvestmentArrayResponse, error) {
	ids := utils.ParseUUIDs(req.Ids)

	invests, err := u.inRepo.GetInvestmentByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return &pb.InvestmentArrayResponse{
		Invest: mapper.ToInvestProtoSlice(invests),
	}, nil
}

func (u *InvestmentUsecase) GetBatchInvestmentByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.InvestmentArrayResponse, error) {
	ids := utils.ParseUUIDs(req.Ids)

	invests, err := u.inRepo.GetBatchInvestmentByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return &pb.InvestmentArrayResponse{
		Invest: mapper.ToInvestProtoSlice(invests),
	}, nil
}

func (u *InvestmentUsecase) GetInvestmentByID(ctx context.Context, req *pb.GetAssetByIDRequest) (*pb.InvestmentResponse, error) {
	id, err := utils.ParseID(req.Id)
	if err != nil {
		return nil, err
	}

	in, err := u.inRepo.GetInvestmentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.InvestmentResponse{
		Success: true,
		Invest:  mapper.ToInvestProto(in),
	}, nil
}

func (u *InvestmentUsecase) UpdateInvestment(ctx context.Context, req *pb.UpdateInvestmentRequest) (*pb.InvestmentResponse, error) {
	if req.Invest == nil {
		return nil, errors.New("investment data is required")
	}

	id, uid, err := utils.ValidateIDs(req.Id, req.Invest.UserId)
	if err != nil {
		return nil, err
	}

	in, err := u.inRepo.GetInvestmentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := helper.ApplyUpdateInFields(req, in); err != nil {
		return nil, err
	}

	syncParams := domain.FileSyncParams{
		UserID:        uid,
		EntityID:      id,
		EntityType:    "investment",
		NewFiles:      req.NewFiles,
		DeleteFileIDs: req.DeleteFileIds,
	}

	if err := u.assetHelper.SyncFiles(ctx, syncParams); err != nil {
		return nil, err
	}

	updatedInvest, err := u.inRepo.UpdateInvestment(ctx, in)
	if err != nil {
		return nil, err
	}

	return &pb.InvestmentResponse{
		Success: true,
		Invest:  mapper.ToInvestProto(updatedInvest),
	}, nil
}

func (u *InvestmentUsecase) DeleteInvestment(ctx context.Context, req *pb.DeleteAssetRequest) (*pb.DeleteAssetResponse, error) {
	id, uid, err := utils.ValidateIDs(req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	if _, err = u.inRepo.GetInvestmentByID(ctx, id); err != nil {
		return nil, err
	}

	if err := u.inRepo.SoftDeleteInvestment(ctx, id, uid); err != nil {
		return nil, err
	}

	go func() {
		bgCtx := context.Background()

		_, err := u.userClient.MarkAssetMessagesDeleted(bgCtx, &userPb.MarkAssetDeletedRequest{
			AssetId: id.String(),
		})

		if err != nil {
			log.Printf("⚠️ Failed to notify User Service via gRPC: %v", err)
		}
	}()

	return &pb.DeleteAssetResponse{
		Success: true,
	}, nil
}

func (u *InvestmentUsecase) CleanupExpiredInvestment(ctx context.Context) error {
	cutoffTime := time.Now().AddDate(0, 0, -7)
	expiredInvest, err := u.inRepo.GetExpiredInvestment(ctx, cutoffTime)
	if err != nil {
		return err
	}

	if len(expiredInvest) == 0 {
		return nil
	}

	for _, inv := range expiredInvest {
		u.assetHelper.CleanupResource(ctx, inv.ID, inv.Files, func(id uuid.UUID) error {
			return u.inRepo.HardDeleteInvestment(ctx, id)
		})
	}

	return nil
}
