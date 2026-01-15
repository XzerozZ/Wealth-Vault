package usecase

import (
	"context"
	"wealth-vault/asset-service/internal/domain"
	repo "wealth-vault/asset-service/internal/repository/interface"

	"github.com/google/uuid"
)

type AssetUsecase struct {
	assetRepo repo.AssetRepository
}

func NewAssetUsecase(r repo.AssetRepository) AssetUsecase {
	return AssetUsecase{assetRepo: r}
}

func (u *AssetUsecase) CreateCash(ctx context.Context, cash *domain.Cash) (string, error) {
	cash.ID = uuid.NewString()
	if err := u.assetRepo.CreateCash(ctx, cash); err != nil {
		return "", err
	}

	return cash.ID, nil
}

func (u *AssetUsecase) GetCash(ctx context.Context, uid string) ([]domain.Cash, error) {
	cash, err := u.assetRepo.GetCash(ctx, uid)
	if err != nil {
		return nil, err
	}

	return cash, nil
}

func (u *AssetUsecase) GetCashByID(ctx context.Context, id string, uid string) (*domain.Cash, error) {
	cash, err := u.assetRepo.GetCashByID(ctx, id, uid)
	if err != nil {
		return nil, err
	}

	return cash, nil
}

func (u *AssetUsecase) UpdateCash(ctx context.Context, input *domain.UpdateCashInput) (*domain.Cash, error) {
	updateData := &domain.Cash{
		ID:          input.ID,
		Name:        input.Name,
		Value:       input.Value,
		Description: input.Description,
		UserID:      input.UserID,
	}

	updatedCash, err := u.assetRepo.UpdateCash(ctx, updateData, input.UpdateMask)
	if err != nil {
		return nil, err
	}

	return updatedCash, nil
}
