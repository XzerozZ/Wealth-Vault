package usecase_test

import (
	"context"
	"testing"
	"time"
	"wealth-vault/asset-service/internal/domain"
	"wealth-vault/asset-service/internal/usecase"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	mock_client "wealth-vault/asset-service/test/mock/client"
	mock_event "wealth-vault/asset-service/test/mock/event"
	mock_helper "wealth-vault/asset-service/test/mock/helper"
	mock_repo "wealth-vault/asset-service/test/mock/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInsuranceUsecase(t *testing.T) {
	repo := new(mock_repo.MockInsuranceRepository)
	assetHelper := new(mock_helper.MockAssetHelper)
	publisher := new(mock_event.MockPublisher)
	userClient := new(mock_client.MockUserClient)
	uc := usecase.NewInsuranceUsecase(repo, assetHelper, publisher, userClient)

	ctx := context.Background()
	userID := uuid.New()
	insID := uuid.New()

	t.Run("CreateInsurance_Success", func(t *testing.T) {
		req := &pb.CreateInsuranceRequest{
			UserId: userID.String(),
			Name:   "Life Premium",
		}
		repo.On("CreateInsurance", ctx, mock.AnythingOfType("*domain.Insurance")).Return(nil).Once()

		res, err := uc.CreateInsurance(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, res.Insurance)
		repo.AssertExpectations(t)
	})

	t.Run("GetInsurance_Success", func(t *testing.T) {
		req := &pb.GetAssetRequest{UserId: userID.String()}
		insurances := []*domain.Insurance{{ID: insID, Name: "Life Premium"}}

		repo.On("GetInsurance", ctx, userID).Return(insurances, nil).Once()

		res, err := uc.GetInsurance(ctx, req)
		assert.NoError(t, err)
		assert.Len(t, res.Insurance, 1)
		assert.True(t, res.Success)
	})

	t.Run("GetInsuranceByID_Success", func(t *testing.T) {
		req := &pb.GetAssetByIDRequest{Id: insID.String()}
		repo.On("GetInsuranceByID", ctx, insID).Return(&domain.Insurance{ID: insID}, nil).Once()

		res, err := uc.GetInsuranceByID(ctx, req)
		assert.NoError(t, err)
		assert.True(t, res.Success)
		assert.Equal(t, insID.String(), res.Insurance.Id)
	})

	t.Run("UpdateInsurance_Success", func(t *testing.T) {
		req := &pb.UpdateInsuranceRequest{
			Id: insID.String(),
			Insurance: &pb.Insurance{
				UserId: userID.String(),
				Name:   "Updated Life Premium",
			},
			DeleteFileIds: []string{uuid.New().String()},
		}
		existing := &domain.Insurance{ID: insID, UserID: userID}

		repo.On("GetInsuranceByID", ctx, insID).Return(existing, nil).Once()

		assetHelper.On("SyncFiles", ctx, mock.MatchedBy(func(p domain.FileSyncParams) bool {
			return p.EntityID == insID && p.EntityType == "insurance"
		})).Return(nil).Once()

		repo.On("UpdateInsurance", ctx, mock.Anything).Return(existing, nil).Once()

		res, err := uc.UpdateInsurance(ctx, req)
		assert.NoError(t, err)
		assert.True(t, res.Success)
		assetHelper.AssertExpectations(t)
	})

	t.Run("DeleteInsurance_Success", func(t *testing.T) {
		req := &pb.DeleteAssetRequest{Id: insID.String(), UserId: userID.String()}

		repo.On("GetInsuranceByID", ctx, insID).Return(&domain.Insurance{ID: insID}, nil).Once()
		repo.On("SoftDeleteInsurances", ctx, insID, userID).Return(nil).Once()

		res, err := uc.DeleteInsurance(ctx, req)
		assert.NoError(t, err)
		assert.True(t, res.Success)
	})

	t.Run("CleanupExpiredInsurances_Success", func(t *testing.T) {
		expiredInsurances := []domain.Insurance{
			{ID: insID, Files: []domain.FileAssociate{{ID: uuid.New()}}},
		}

		repo.On("GetExpiredInsurances", ctx, mock.AnythingOfType("time.Time")).Return(expiredInsurances, nil).Once()

		assetHelper.On("CleanupResource", ctx, insID, expiredInsurances[0].Files, mock.Anything).
			Run(func(args mock.Arguments) {
				callback := args.Get(3).(func(uuid.UUID) error)
				callback(insID)
			}).Return().Once()

		repo.On("HardDeleteInsurances", ctx, insID).Return(nil).Once()

		err := uc.CleanupExpiredInsurances(ctx)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("CheckExpiringInsurances_Success", func(t *testing.T) {
		expDate := time.Now().AddDate(0, 0, 7)
		insurances := []*domain.Insurance{
			{
				ID:      insID,
				UserID:  userID,
				Name:    "Life Premium",
				ExpDate: &expDate,
			},
		}

		repo.On("GetExpiringInsurances", ctx, 7).Return(insurances, nil).Once()

		repo.On("GetExpiringInsurances", ctx, mock.Anything).Return([]*domain.Insurance{}, nil)

		publisher.On("Publish", "noti.insurance.expiring", mock.MatchedBy(func(evt domain.InsuranceExpiringEvent) bool {
			return evt.InsuranceID == insID.String() && evt.DaysLeft == 7
		})).Return(nil).Once()

		err := uc.CheckExpiringInsurances(ctx)
		assert.NoError(t, err)
		publisher.AssertExpectations(t)
	})

	t.Run("UpdateInsurance_InvalidID", func(t *testing.T) {
		req := &pb.UpdateInsuranceRequest{Id: "invalid-uuid", Insurance: &pb.Insurance{UserId: userID.String()}}
		res, err := uc.UpdateInsurance(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("UpdateInsurance_MissingData", func(t *testing.T) {
		req := &pb.UpdateInsuranceRequest{Id: insID.String(), Insurance: nil}
		res, err := uc.UpdateInsurance(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, "insurance data is required", err.Error())
		assert.Nil(t, res)
	})
}
