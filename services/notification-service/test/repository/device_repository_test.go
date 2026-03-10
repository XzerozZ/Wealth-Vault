package repository_test

import (
	"context"
	"testing"
	"time"

	"wealth-vault/notification-service/internal/domain"
	"wealth-vault/notification-service/internal/repository"
	testutil "wealth-vault/notification-service/test/testdb"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeviceRepository(t *testing.T) {
	t.Run("RegisterDevice_Success", func(t *testing.T) {
		mockDB := testutil.NewMockDB(t)
		defer mockDB.Close()
		repo := repository.NewDeviceRepository(mockDB.DB)

		userID := uuid.New()
		fcmToken := "fcm-token-123"

		mockDB.Mock.ExpectBegin()
		mockDB.Mock.ExpectQuery(`(?i)INSERT INTO "device_tokens"`).
			WithArgs(
				userID,
				fcmToken,
				"ios",
				"iPhone 15",
				true,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

		mockDB.Mock.ExpectCommit()

		device := domain.DeviceToken{
			ID:         uuid.New(),
			UserID:     userID,
			Token:      fcmToken,
			Platform:   "ios",
			DeviceName: "iPhone 15",
			IsActive:   true,
		}

		err := repo.RegisterDevice(context.Background(), &device)

		assert.NoError(t, err)
		mockDB.ExpectDone(t)
	})

	t.Run("GetActiveTokens_Success", func(t *testing.T) {
		mock := testutil.NewMockDB(t)
		defer mock.Close()

		repo := repository.NewDeviceRepository(mock.DB)

		userID := uuid.New()

		rows := sqlmock.NewRows([]string{
			"id", "user_id", "token", "platform", "device_name", "is_active", "created_at", "updated_at",
		}).AddRow(
			uuid.New(),
			userID,
			"fcm-token-123",
			"ios",
			"iPhone 15",
			true,
			time.Now(),
			time.Now(),
		)

		mock.Mock.ExpectQuery(`SELECT .* FROM .*device_tokens.*`).
			WithArgs(userID).
			WillReturnRows(rows)

		tokens, err := repo.GetActiveTokens(context.Background(), userID)

		assert.NoError(t, err)
		assert.Len(t, tokens, 1)
		assert.Equal(t, "fcm-token-123", tokens[0].Token)

		assert.NoError(t, mock.Mock.ExpectationsWereMet())
	})

	t.Run("UnregisterDevice_Success", func(t *testing.T) {
		mock := testutil.NewMockDB(t)
		defer mock.Close()

		repo := repository.NewDeviceRepository(mock.DB)

		userID := uuid.New()
		token := "fcm-token-123"

		mock.Mock.ExpectBegin()

		mock.Mock.ExpectExec(`UPDATE .*device_tokens.*`).
			WithArgs(false, sqlmock.AnyArg(), userID, token).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.Mock.ExpectCommit()

		err := repo.UnregisterDevice(context.Background(), userID, token)

		assert.NoError(t, err)
		assert.NoError(t, mock.Mock.ExpectationsWereMet())
	})

	t.Run("MarkTokenInactive_Success", func(t *testing.T) {
		mock := testutil.NewMockDB(t)
		defer mock.Close()

		repo := repository.NewDeviceRepository(mock.DB)

		token := "fcm-token-123"

		mock.Mock.ExpectBegin()

		mock.Mock.ExpectExec(`UPDATE .*device_tokens.*`).
			WithArgs(false, sqlmock.AnyArg(), token).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.Mock.ExpectCommit()

		err := repo.MarkTokenInactive(context.Background(), token)

		assert.NoError(t, err)
		assert.NoError(t, mock.Mock.ExpectationsWereMet())
	})
}
