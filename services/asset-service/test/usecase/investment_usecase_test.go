package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"wealth-vault/asset-service/internal/domain"
	"wealth-vault/asset-service/internal/usecase"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	mock_helper "wealth-vault/asset-service/test/mock/helper"
	mock_repo "wealth-vault/asset-service/test/mock/repository"
)

func TestInvestmentUsecase(t *testing.T) {
	repo := new(mock_repo.MockInvestmentRepository)
	assetHelper := new(mock_helper.MockAssetHelper)
	uc := usecase.NewInvestmentUsecase(repo, assetHelper)

	ctx := context.Background()
	userID := uuid.New()
	investID := uuid.New()

	t.Run("CreateInvestment_Success", func(t *testing.T) {
		req := &pb.CreateInvestmentRequest{
			UserId: userID.String(),
			Name:   "Stock Portfolio",
		}
		repo.On("CreateInvestment", ctx, mock.AnythingOfType("*domain.Investment")).Return(nil).Once()

		res, err := uc.CreateInvestment(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, res.Invest)
		repo.AssertExpectations(t)
	})

	t.Run("GetInvestment_Success", func(t *testing.T) {
		req := &pb.GetAssetRequest{UserId: userID.String()}
		// ใช้ []*domain.Investment เพื่อให้ตรงกับ Mock Repository
		invests := []*domain.Investment{{ID: investID, Name: "Stock Portfolio"}}

		repo.On("GetInvestment", ctx, userID).Return(invests, nil).Once()

		res, err := uc.GetInvestment(ctx, req)
		assert.NoError(t, err)
		assert.Len(t, res.Invest, 1)
		assert.True(t, res.Success)
	})

	t.Run("GetInvestmentByID_Success", func(t *testing.T) {
		req := &pb.GetAssetByIDRequest{Id: investID.String()}
		repo.On("GetInvestmentByID", ctx, investID).Return(&domain.Investment{ID: investID}, nil).Once()

		res, err := uc.GetInvestmentByID(ctx, req)
		assert.NoError(t, err)
		assert.True(t, res.Success)
		assert.Equal(t, investID.String(), res.Invest.Id)
	})

	t.Run("UpdateInvestment_Success", func(t *testing.T) {
		req := &pb.UpdateInvestmentRequest{
			Id: investID.String(),
			Invest: &pb.Investment{
				UserId: userID.String(),
				Name:   "Updated Portfolio",
			},
		}
		existing := &domain.Investment{ID: investID, UserID: userID}

		repo.On("GetInvestmentByID", ctx, investID).Return(existing, nil).Once()
		assetHelper.On("SyncFiles", ctx, mock.MatchedBy(func(p domain.FileSyncParams) bool {
			return p.EntityID == investID && p.EntityType == "investment"
		})).Return(nil).Once()
		repo.On("UpdateInvestment", ctx, mock.Anything).Return(existing, nil).Once()

		res, err := uc.UpdateInvestment(ctx, req)
		assert.NoError(t, err)
		assert.True(t, res.Success)
		assetHelper.AssertExpectations(t)
	})

	t.Run("DeleteInvestment_Success", func(t *testing.T) {
		req := &pb.DeleteAssetRequest{Id: investID.String(), UserId: userID.String()}

		repo.On("GetInvestmentByID", ctx, investID).Return(&domain.Investment{ID: investID}, nil).Once()
		repo.On("SoftDeleteInvestment", ctx, investID, userID).Return(nil).Once()

		res, err := uc.DeleteInvestment(ctx, req)
		assert.NoError(t, err)
		assert.True(t, res.Success)
	})

	t.Run("CleanupExpiredInvestment_Success", func(t *testing.T) {
		expired := []domain.Investment{
			{ID: investID, Files: []domain.FileAssociate{{ID: uuid.New()}}},
		}

		repo.On("GetExpiredInvestment", ctx, mock.AnythingOfType("time.Time")).Return(expired, nil).Once()

		assetHelper.On("CleanupResource", ctx, investID, expired[0].Files, mock.Anything).
			Run(func(args mock.Arguments) {
				cb := args.Get(3).(func(uuid.UUID) error)
				cb(investID)
			}).Return().Once()

		repo.On("HardDeleteInvestment", ctx, investID).Return(nil).Once()

		err := uc.CleanupExpiredInvestment(ctx)
		assert.NoError(t, err)
	})

	t.Run("UpdateInvestment_MissingData", func(t *testing.T) {
		req := &pb.UpdateInvestmentRequest{Id: investID.String(), Invest: nil}
		res, err := uc.UpdateInvestment(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, "investment data is required", err.Error())
		assert.Nil(t, res)
	})
}
