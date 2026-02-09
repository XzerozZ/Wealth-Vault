package usecase

import (
	"context"
	"errors"
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

type CashUsecase struct {
	cashRepo   repo.CashRepository
	fileRepo   repo.FileRepository
	storage    *utils.StorageClient
	userClient userPb.UserServiceClient
}

func NewCashUsecase(r repo.CashRepository, fr repo.FileRepository, s *utils.StorageClient, userClient userPb.UserServiceClient) CashUsecase {
	return CashUsecase{
		cashRepo:   r,
		fileRepo:   fr,
		storage:    s,
		userClient: userClient,
	}
}

func (u *CashUsecase) CreateCash(ctx context.Context, req *pb.CreateCashRequest) (*pb.CashResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	var domainFiles []domain.FileAssociate
	if len(req.NewFiles) > 0 {
		for _, f := range req.NewFiles {
			domainFiles = append(domainFiles, domain.FileAssociate{
				Link:     f.Url,
				FileType: f.FileType,
				UserID:   userID,
			})
		}
	}

	cash := &domain.Cash{
		UserID:      userID,
		Name:        req.Name,
		Amount:      req.Amount,
		Description: req.Description,
		Files:       domainFiles,
	}

	if err := u.cashRepo.CreateCash(ctx, cash); err != nil {
		return nil, err
	}

	return &pb.CashResponse{
		Cash: mapper.ToCashProto(cash),
	}, nil
}

func (u *CashUsecase) GetCash(ctx context.Context, req *pb.GetAssetRequest) (*pb.CashArrayResponse, error) {
	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	cash, err := u.cashRepo.GetCash(ctx, uid)
	if err != nil {
		return nil, err
	}

	var CashList []*pb.Cash
	for _, item := range cash {
		CashList = append(CashList, mapper.ToCashProto(item))
	}

	return &pb.CashArrayResponse{
		Success: true,
		Cash:    CashList,
	}, nil
}

func (u *CashUsecase) GetCashByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.CashArrayResponse, error) {
	var ids []uuid.UUID
	for _, idStr := range req.Ids {
		if parsedID, err := uuid.Parse(idStr); err == nil {
			ids = append(ids, parsedID)
		}
	}

	bu, err := u.cashRepo.GetCashByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	var pbCash []*pb.Cash
	for _, a := range bu {
		pbCash = append(pbCash, mapper.ToCashProto(a))
	}

	return &pb.CashArrayResponse{
		Cash: pbCash,
	}, nil
}

func (u *CashUsecase) GetBatchCashByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.CashArrayResponse, error) {
	var ids []uuid.UUID
	for _, idStr := range req.Ids {
		if parsedID, err := uuid.Parse(idStr); err == nil {
			ids = append(ids, parsedID)
		}
	}

	bu, err := u.cashRepo.GetBatchCashByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	var pbCash []*pb.Cash
	for _, a := range bu {
		pbCash = append(pbCash, mapper.ToCashProto(a))
	}

	return &pb.CashArrayResponse{
		Cash: pbCash,
	}, nil
}

func (u *CashUsecase) GetCashByID(ctx context.Context, req *pb.GetAssetByIDRequest) (*pb.CashResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	cash, err := u.cashRepo.GetCashByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.CashResponse{
		Success: true,
		Cash:    mapper.ToCashProto(cash),
	}, nil
}

func (u *CashUsecase) UpdateCash(ctx context.Context, req *pb.UpdateCashRequest) (*pb.CashResponse, error) {
	id, uid, err := utils.ValidateIDs(req.Id, req.Cash.UserId)
	if err != nil {
		return nil, err
	}

	cash, err := u.cashRepo.GetCashByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := helper.ApplyUpdateCashFields(req, cash); err != nil {
		return nil, err
	}

	syncParams := domain.FileSyncParams{
		UserID:        uid,
		EntityID:      id,
		EntityType:    "cash",
		NewFiles:      req.NewFiles,
		DeleteFileIDs: req.DeleteFileIds,
	}

	err = helper.SyncEntityFiles(ctx, u.fileRepo, u.storage, syncParams)
	if err != nil {
		return nil, err
	}

	updatedCash, err := u.cashRepo.UpdateCash(ctx, cash)
	if err != nil {
		return nil, err
	}

	return &pb.CashResponse{
		Success: true,
		Cash:    mapper.ToCashProto(updatedCash),
	}, nil
}

func (u *CashUsecase) DeleteCash(ctx context.Context, req *pb.DeleteAssetRequest) (*pb.DeleteAssetResponse, error) {
	id, uid, err := utils.ValidateIDs(req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	_, err = u.cashRepo.GetCashByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := u.cashRepo.SoftDeleteCash(ctx, id, uid); err != nil {
		return nil, err
	}

	return &pb.DeleteAssetResponse{
		Success: true,
	}, nil
}

func (u *CashUsecase) CleanupExpiredCashes(ctx context.Context) error {
	cutoffTime := time.Now().AddDate(0, 0, -7)
	GetExpiredCash, err := u.cashRepo.GetExpiredCash(ctx, cutoffTime)
	if err != nil {
		return err
	}

	if len(GetExpiredCash) == 0 {
		return err
	}

	for _, c := range GetExpiredCash {
		helper.CleanupAssetResource(
			ctx,
			c.ID,
			c.Files,
			u.storage,
			u.userClient,
			func(id uuid.UUID) error {
				return u.cashRepo.HardDeleteCash(ctx, id)
			},
		)
	}

	return nil
}
