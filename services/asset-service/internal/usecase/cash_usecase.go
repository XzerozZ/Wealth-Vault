package usecase

import (
	"context"
	"wealth-vault/asset-service/internal/domain"
	repo "wealth-vault/asset-service/internal/repository/interface"

	"github.com/google/uuid"
)

type CashUsecase struct {
	assetRepo repo.CashRepository
	fileRepo  repo.FileRepository
}

func NewCashUsecase(r repo.CashRepository, fr repo.FileRepository) CashUsecase {
	return CashUsecase{
		assetRepo: r,
		fileRepo:  fr,
	}
}

func (u *CashUsecase) CreateCash(ctx context.Context, cash *domain.Cash) (string, error) {
	cash.ID = uuid.NewString()
	for i := range cash.Files {
		cash.Files[i].ID = uuid.NewString()
	}

	if err := u.assetRepo.CreateCash(ctx, cash); err != nil {
		return "", err
	}

	return cash.ID, nil
}

func (u *CashUsecase) GetCash(ctx context.Context, uid string) ([]domain.Cash, error) {
	cash, err := u.assetRepo.GetCash(ctx, uid)
	if err != nil {
		return nil, err
	}

	return cash, nil
}

func (u *CashUsecase) GetCashByID(ctx context.Context, id string, uid string) (*domain.Cash, error) {
	cash, err := u.assetRepo.GetCashByID(ctx, id, uid)
	if err != nil {
		return nil, err
	}

	return cash, nil
}

func (u *CashUsecase) UpdateCash(ctx context.Context, input *domain.UpdateCashInput) (*domain.Cash, error) {
	updateData := &domain.Cash{
		ID:          input.ID,
		Name:        input.Name,
		Value:       input.Value,
		Description: input.Description,
		UserID:      input.UserID,
	}

	if len(input.DeleteFileIDs) > 0 {
		if err := u.fileRepo.DeleteFiles(ctx, input.DeleteFileIDs, "cash", input.UserID); err != nil {
			return nil, err
		}
	}

	if len(input.NewFiles) > 0 {
		for i := range input.NewFiles {
			input.NewFiles[i].EntityID = input.ID
			input.NewFiles[i].UserID = input.UserID
			input.NewFiles[i].EntityType = "cash"
		}

		if err := u.fileRepo.CreateFiles(ctx, input.NewFiles); err != nil {
			return nil, err
		}
	}

	updatedCash, err := u.assetRepo.UpdateCash(ctx, updateData, input.UpdateMask)
	if err != nil {
		return nil, err
	}

	return updatedCash, nil
}

func (u *CashUsecase) DeleteCash(ctx context.Context, id string, uid string) error {
	if err := u.assetRepo.DeleteCash(ctx, id, uid); err != nil {
		return err
	}

	return nil
}
