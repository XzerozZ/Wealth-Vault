package repository_test

import (
	"context"
	"regexp"
	"testing"
	"wealth-vault/user-service/internal/domain"
	"wealth-vault/user-service/internal/repository"
	testutil "wealth-vault/user-service/test/testdb"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateMessage(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewMsgRepository(mock.DB)

	msgID := uuid.New()
	msgs := []domain.GroupMessage{{ID: msgID}}

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectQuery(`(?i)INSERT INTO "group_messages".*RETURNING`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(msgID))
	mock.Mock.ExpectCommit()

	err := repo.CreateMessage(context.Background(), msgs)

	assert.NoError(t, err)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestCreatePrivateMessage(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewMsgRepository(mock.DB)

	msgID := uuid.New()
	msgs := []domain.PrivateMessage{{ID: msgID}}

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectQuery(`(?i)INSERT INTO "private_messages".*RETURNING`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(msgID))
	mock.Mock.ExpectCommit()

	err := repo.CreatePrivateMessage(context.Background(), msgs)

	assert.NoError(t, err)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestGetGroupMessages(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewMsgRepository(mock.DB)

	groupID := uuid.New().String()
	userID := uuid.New().String()
	senderID := uuid.New()

	queryRegex := `^SELECT .* FROM "group_messages" WHERE group_id = \$1 ORDER BY created_at DESC$`
	mock.Mock.ExpectQuery(queryRegex).
		WithArgs(groupID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "sender_id"}).AddRow(uuid.New(), senderID))

	mock.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)).
		WithArgs(senderID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(senderID, "TestUser"))

	res, err := repo.GetGroupMessages(context.Background(), groupID, userID)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestGetPrivateMessages(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewMsgRepository(mock.DB)

	userID := uuid.New().String()
	friendID := uuid.New().String()
	senderUUID := uuid.New()

	queryRegex := `(?i)SELECT .* FROM "private_messages" WHERE .*sender_id = \$1.* ORDER BY created_at DESC`
	mock.Mock.ExpectQuery(queryRegex).
		WithArgs(userID, friendID, friendID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "sender_id"}).AddRow(uuid.New(), senderUUID))

	mock.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)).
		WithArgs(senderUUID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(senderUUID))

	res, err := repo.GetPrivateMessages(context.Background(), userID, friendID)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}
