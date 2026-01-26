package usecase

import (
	"context"
	"errors"
	"wealth-vault/asset-service/internal/domain"
	repo "wealth-vault/asset-service/internal/repository/interface"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	"wealth-vault/asset-service/pkg/utils"
	helper "wealth-vault/asset-service/pkg/utils/helper"
	"wealth-vault/asset-service/pkg/utils/mapper"

	"github.com/google/uuid"
)

type AccountUsecase struct {
	accRepo  repo.AccountRepository
	fileRepo repo.FileRepository
	storage  *utils.StorageClient
}

func NewAccountUsecase(r repo.AccountRepository, fr repo.FileRepository, s *utils.StorageClient) AccountUsecase {
	return AccountUsecase{
		accRepo:  r,
		fileRepo: fr,
		storage:  s,
	}
}

func (u *AccountUsecase) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.AccountResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	accType := domain.BankTypeSavings
	if val, ok := helper.ProtoToDomainAccType[req.Type]; ok {
		accType = val
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

	acc := &domain.Account{
		UserID:      userID,
		Name:        req.Name,
		Amount:      req.Amount,
		BankName:    req.BankName,
		BankAccount: req.BankAcc,
		Type:        accType,
		Description: req.Description,
		Files:       domainFiles,
	}

	if err := u.accRepo.CreateAccount(ctx, acc); err != nil {
		return nil, err
	}

	return &pb.AccountResponse{
		Account: mapper.ToBankProto(acc),
	}, nil
}

func (u *AccountUsecase) GetAccount(ctx context.Context, req *pb.GetAssetRequest) (*pb.AccountArrayResponse, error) {
	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	accounts, err := u.accRepo.GetAccount(ctx, uid)
	if err != nil {
		return nil, err
	}

	var AccountList []*pb.Account
	for _, item := range accounts {
		AccountList = append(AccountList, mapper.ToBankProto(item))
	}

	return &pb.AccountArrayResponse{
		Success: true,
		Account: AccountList,
	}, nil
}

func (u *AccountUsecase) GetAccountByID(ctx context.Context, req *pb.GetAssetByIDRequest) (*pb.AccountResponse, error) {
	id, uid, err := utils.ValidateIDs(req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	acc, err := u.accRepo.GetAccountByID(ctx, id, uid)
	if err != nil {
		return nil, err
	}

	return &pb.AccountResponse{
		Success: true,
		Account: mapper.ToBankProto(acc),
	}, nil
}

func (u *AccountUsecase) UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.AccountResponse, error) {
	id, uid, err := utils.ValidateIDs(req.Id, req.Acc.UserId)
	if err != nil {
		return nil, err
	}

	acc, err := u.accRepo.GetAccountByID(ctx, id, uid)
	if err != nil {
		return nil, err
	}

	if err := helper.ApplyUpdateAccFields(req, acc); err != nil {
		return nil, err
	}

	syncParams := domain.FileSyncParams{
		UserID:        uid,
		EntityID:      id,
		EntityType:    "account",
		NewFiles:      req.NewFiles,
		DeleteFileIDs: req.DeleteFileIds,
	}

	err = helper.SyncEntityFiles(ctx, u.fileRepo, u.storage, syncParams)
	if err != nil {
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

	existingAcc, err := u.accRepo.GetAccountByID(ctx, id, uid)
	if err != nil {
		return nil, err
	}

	if err := u.accRepo.DeleteAccount(ctx, id, uid); err != nil {
		return nil, err
	}

	if len(existingAcc.Files) > 0 {
		fileURLs := make([]string, len(existingAcc.Files))
		for i, f := range existingAcc.Files {
			fileURLs[i] = f.Link
		}

		helper.DeleteFilesAsync(u.storage, fileURLs)
	}

	return &pb.DeleteAssetResponse{
		Success: true,
	}, nil
}
