package usecase_test

import (
	"context"
	"errors"
	"testing"
	"wealth-vault/user-service/internal/domain"
	usecase "wealth-vault/user-service/internal/usecase/shareItem"
	assetPb "wealth-vault/user-service/pkg/pb/proto/asset"
	mock_client "wealth-vault/user-service/test/mock/client"
	mock_event "wealth-vault/user-service/test/mock/event"
	mock_mail "wealth-vault/user-service/test/mock/mail"
	mock_repo "wealth-vault/user-service/test/mock/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShareItemUsecase_EmailFlows(t *testing.T) {
	t.Run("SendEmailInvitations - success building", func(t *testing.T) {
		mockItem := new(mock_repo.MockShareItemRepository)
		mockGroup := new(mock_repo.MockGroupRepository)
		mockUser := new(mock_repo.MockUserRepository)
		mockMsg := new(mock_repo.MockMsgRepository)
		mockAsset := new(mock_client.MockAssetClient)
		mockMail := new(mock_mail.MockMailClient)
		mockPub := new(mock_event.MockEventPublisher)
		uc := usecase.NewShareItemUsecase(mockItem, mockGroup, mockUser, mockMsg, mockAsset, mockMail, mockPub)

		id := uuid.New()

		items := []domain.EmailItem{
			{
				EntityID:   id,
				EntityType: usecase.AssetTypeBuilding,
				Email:      "test@mail.com",
			},
		}

		mockAsset.On("GetBatchBuilding", mock.Anything, mock.Anything, mock.Anything).
			Return(&assetPb.BuildingArrayResponse{
				Building: []*assetPb.Building{
					{Id: id.String(), Name: "ตึก A"},
				},
			}, nil)

		mockMail.On("SendShareInvitation", mock.Anything, mock.MatchedBy(func(req domain.SendEmailRequest) bool {
			return req.AssetName == "ตึก A"
		})).Return(nil)

		uc.SendEmailInvitations(items)

		mockAsset.AssertExpectations(t)
		mockMail.AssertExpectations(t)
	})

	t.Run("SendEmailInvitations - fallback name", func(t *testing.T) {
		mockAsset := new(mock_client.MockAssetClient)
		mockMail := new(mock_mail.MockMailClient)

		uc := usecase.NewShareItemUsecase(nil, nil, nil, nil, mockAsset, mockMail, nil)

		id := uuid.New()

		items := []domain.EmailItem{
			{
				EntityID:   id,
				EntityType: usecase.AssetTypeBuilding,
				Email:      "fallback@mail.com",
			},
		}

		mockAsset.On("GetBatchBuilding", mock.Anything, mock.Anything, mock.Anything).
			Return(&assetPb.BuildingArrayResponse{}, nil)

		mockMail.On("SendShareInvitation", mock.Anything, mock.MatchedBy(func(req domain.SendEmailRequest) bool {
			return req.AssetName == "รายการทรัพย์สิน"
		})).Return(nil)

		uc.SendEmailInvitations(items)

		mockMail.AssertExpectations(t)
	})

	t.Run("ProcessScheduledEmails - success", func(t *testing.T) {
		mockItem := new(mock_repo.MockShareItemRepository)
		mockGroup := new(mock_repo.MockGroupRepository)
		mockUser := new(mock_repo.MockUserRepository)
		mockMsg := new(mock_repo.MockMsgRepository)
		mockAsset := new(mock_client.MockAssetClient)
		mockMail := new(mock_mail.MockMailClient)
		mockPub := new(mock_event.MockEventPublisher)
		uc := usecase.NewShareItemUsecase(mockItem, mockGroup, mockUser, mockMsg, mockAsset, mockMail, mockPub)

		id := uuid.New()
		emailID := uuid.New()

		pending := []domain.EmailItem{
			{
				ID:         emailID,
				EntityID:   id,
				EntityType: usecase.AssetTypeBuilding,
				Email:      "schedule@mail.com",
			},
		}

		mockItem.On("GetPendingEmails", mock.Anything).
			Return(pending, nil)

		mockAsset.On("GetBatchBuilding", mock.Anything, mock.Anything, mock.Anything).
			Return(&assetPb.BuildingArrayResponse{
				Building: []*assetPb.Building{
					{Id: id.String(), Name: "ตึก B"},
				},
			}, nil)

		mockMail.On("SendShareInvitation", mock.Anything, mock.Anything).
			Return(nil)

		mockItem.On("MarkEmailsAsSent", mock.Anything, []uuid.UUID{emailID}).
			Return(nil)

		err := uc.ProcessScheduledEmails(context.Background())

		assert.NoError(t, err)

		mockItem.AssertExpectations(t)
		mockMail.AssertExpectations(t)
	})

	t.Run("ProcessScheduledEmails - no pending", func(t *testing.T) {
		mockRepo := new(mock_repo.MockShareItemRepository)

		uc := usecase.NewShareItemUsecase(mockRepo, nil, nil, nil, nil, nil, nil)

		mockRepo.On("GetPendingEmails", mock.Anything).
			Return([]domain.EmailItem{}, nil)

		err := uc.ProcessScheduledEmails(context.Background())

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ProcessScheduledEmails - repo error", func(t *testing.T) {
		mockRepo := new(mock_repo.MockShareItemRepository)

		uc := usecase.NewShareItemUsecase(mockRepo, nil, nil, nil, nil, nil, nil)

		mockRepo.On("GetPendingEmails", mock.Anything).
			Return(nil, errors.New("db error"))

		err := uc.ProcessScheduledEmails(context.Background())

		assert.Error(t, err)
	})
}
