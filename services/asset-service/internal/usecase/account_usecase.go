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

type AccountUsecase struct {
	accRepo     repo.AccountRepository
	assetHelper helper.AssetHelper
	userClient  userPb.UserServiceClient
}

func NewAccountUsecase(r repo.AccountRepository, ah helper.AssetHelper, uc userPb.UserServiceClient) *AccountUsecase {
	return &AccountUsecase{
		accRepo:     r,
		assetHelper: ah,
		userClient:  uc,
	}
}

func (u *AccountUsecase) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.AccountResponse, error) {
	uid, err := utils.ParseID(req.UserId)
	if err != nil {
		return nil, err
	}

	acc := mapper.ToAccountDomain(req, uid)

	if err := u.accRepo.CreateAccount(ctx, acc); err != nil {
		return nil, err
	}

	return &pb.AccountResponse{
		Account: mapper.ToBankProto(acc),
	}, nil
}

func (u *AccountUsecase) GetAccount(ctx context.Context, req *pb.GetAssetRequest) (*pb.AccountArrayResponse, error) {
	uid, err := utils.ParseID(req.UserId)
	if err != nil {
		return nil, err
	}

	accounts, err := u.accRepo.GetAccount(ctx, uid)
	if err != nil {
		return nil, err
	}

	return &pb.AccountArrayResponse{
		Success: true,
		Account: mapper.ToBankProtoSlice(accounts),
	}, nil
}

func (u *AccountUsecase) GetAccountByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.AccountArrayResponse, error) {
	ids := utils.ParseUUIDs(req.Ids)
	acc, err := u.accRepo.GetAccountByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return &pb.AccountArrayResponse{
		Account: mapper.ToBankProtoSlice(acc),
	}, nil
}

func (u *AccountUsecase) GetBatchAccountByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.AccountArrayResponse, error) {
	ids := utils.ParseUUIDs(req.Ids)

	acc, err := u.accRepo.GetBatchAccountByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return &pb.AccountArrayResponse{
		Account: mapper.ToBankProtoSlice(acc),
	}, nil
}

func (u *AccountUsecase) GetAccountByID(ctx context.Context, req *pb.GetAssetByIDRequest) (*pb.AccountResponse, error) {
	id, err := utils.ParseID(req.Id)
	if err != nil {
		return nil, err
	}

	acc, err := u.accRepo.GetAccountByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.AccountResponse{
		Success: true,
		Account: mapper.ToBankProto(acc),
	}, nil
}

func (u *AccountUsecase) UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.AccountResponse, error) {
	if req.Acc == nil {
		return nil, errors.New("account data is required")
	}

	id, uid, err := utils.ValidateIDs(req.Id, req.Acc.UserId)
	if err != nil {
		return nil, err
	}

	acc, err := u.accRepo.GetAccountByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := helper.ApplyUpdateAccFields(req, acc); err != nil {
		return nil, err
	}

	syncParams := domain.FileSyncParams{
		Ctx:           ctx,
		UserID:        uid,
		EntityID:      id,
		EntityType:    "account",
		NewFiles:      req.NewFiles,
		DeleteFileIDs: req.DeleteFileIds,
	}

	if err := u.assetHelper.SyncFiles(ctx, syncParams); err != nil {
		return nil, err
	}

	updatedAcc, err := u.accRepo.UpdateAccount(ctx, acc)
	if err != nil {
		return nil, err
	}

	return &pb.AccountResponse{
		Success: true,
		Account: mapper.ToBankProto(updatedAcc),
	}, nil
}

func (u *AccountUsecase) DeleteAccount(ctx context.Context, req *pb.DeleteAssetRequest) (*pb.DeleteAssetResponse, error) {
	id, uid, err := utils.ValidateIDs(req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	if _, err = u.accRepo.GetAccountByID(ctx, id); err != nil {
		return nil, err
	}

	if err := u.accRepo.SoftDeleteAccount(ctx, id, uid); err != nil {
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

func (u *AccountUsecase) CleanupExpiredAccounts(ctx context.Context) error {
	cutoff := time.Now().AddDate(0, 0, -7)
	expiredAccounts, err := u.accRepo.GetExpiredAccounts(ctx, cutoff)
	if err != nil {
		return err
	}

	if len(expiredAccounts) == 0 {
		return nil
	}

	for _, acc := range expiredAccounts {
		u.assetHelper.CleanupResource(ctx, acc.ID, acc.Files, func(id uuid.UUID) error {
			return u.accRepo.HardDeleteAccount(ctx, id)
		})
	}

	return nil
}
