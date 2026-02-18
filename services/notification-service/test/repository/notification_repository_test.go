package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"
	"wealth-vault/notification-service/internal/domain"
	"wealth-vault/notification-service/internal/repository"
	testutil "wealth-vault/notification-service/test/testdb"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNotificationRepository(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewNotificationRepository(mock.DB)

	ctx := context.Background()
	receiverID := uuid.New()
	notiID := uuid.New()

	t.Run("CreateNotification_Success", func(t *testing.T) {
		item := &domain.Notification{
			ID:         notiID,
			EntityType: "building",
			EntityID:   uuid.New(),
			Receiver:   receiverID,
			Channel:    "web",
			Message:    "Test Message",
			IsRead:     false,
		}

		mock.Mock.ExpectBegin()

		mock.Mock.ExpectQuery(`INSERT INTO "notifications"`).
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(item.ID))

		mock.Mock.ExpectCommit()

		err := repo.CreateNotification(ctx, item)
		assert.NoError(t, err)
		assert.NoError(t, mock.Mock.ExpectationsWereMet())
	})

	t.Run("CreateNotification_Error", func(t *testing.T) {
		mock.Mock.ExpectBegin()
		mock.Mock.ExpectQuery(`INSERT INTO "notifications"`).
			WillReturnError(errors.New("db error"))
		mock.Mock.ExpectRollback()

		err := repo.CreateNotification(ctx, &domain.Notification{})

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
		assert.NoError(t, mock.Mock.ExpectationsWereMet())
	})

	t.Run("GetByReceiver_Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "entity_type", "entity_id", "receiver", "sender_id", "channel", "message", "created_at", "is_read",
		}).
			AddRow(notiID, "building", uuid.New(), receiverID, nil, "web", "Hello 1", time.Now(), false).
			AddRow(uuid.New(), "building", uuid.New(), receiverID, nil, "web", "Hello 2", time.Now(), false)

		mock.Mock.ExpectQuery(`SELECT \* FROM "notifications" WHERE receiver = \$1 ORDER BY created_at desc LIMIT \$2`).
			WithArgs(receiverID, 50).
			WillReturnRows(rows)

		list, err := repo.GetByReceiver(ctx, receiverID)

		assert.NoError(t, err)
		if assert.Len(t, list, 2) {
			assert.Equal(t, "Hello 1", list[0].Message)
		}
		assert.NoError(t, mock.Mock.ExpectationsWereMet())
	})

	t.Run("GetByReceiver_Error", func(t *testing.T) {
		mock.Mock.ExpectQuery(`SELECT .* FROM "notifications"`).
			WithArgs(receiverID, 50).
			WillReturnError(errors.New("query error"))

		list, err := repo.GetByReceiver(ctx, receiverID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "query error")
		assert.Nil(t, list)
		assert.NoError(t, mock.Mock.ExpectationsWereMet())
	})
}
