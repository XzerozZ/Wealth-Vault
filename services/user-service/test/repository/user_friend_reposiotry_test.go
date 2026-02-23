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

func TestGetFriendList(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewUserRepository(mock.DB)

	userID := uuid.New()
	friendID := uuid.New()

	mock.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "friend_lists" WHERE user_id = $1 AND status = $2`)).
		WithArgs(userID, "ACCEPTED").
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "friend_id", "status"}).AddRow(userID, friendID, "ACCEPTED"))

	mock.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)).
		WithArgs(friendID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(friendID))

	res, err := repo.GetFriendList(context.Background(), userID)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestAddFriend(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewUserRepository(mock.DB)

	fri := &domain.FriendList{UserID: uuid.New(), FriendID: uuid.New()}

	mock.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "friend_lists" WHERE (user_id = $1 AND friend_id = $2) OR (user_id = $3 AND friend_id = $4)`)).
		WithArgs(fri.UserID, fri.FriendID, fri.FriendID, fri.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "friend_lists"`)).
		WithArgs(fri.UserID, fri.FriendID, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.Mock.ExpectCommit()

	err := repo.AddFriend(context.Background(), fri)

	assert.NoError(t, err)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestCreateFriendship(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewUserRepository(mock.DB)

	fri := &domain.FriendList{UserID: uuid.New(), FriendID: uuid.New()}

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "friend_lists"`)).
		WithArgs(fri.UserID, fri.FriendID, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.Mock.ExpectCommit()

	err := repo.CreateFriendship(context.Background(), fri)

	assert.NoError(t, err)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestRemoveFriend(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewUserRepository(mock.DB)

	userID := uuid.New()
	friendID := uuid.New()

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "friend_lists" WHERE (user_id = $1 AND friend_id = $2) OR (user_id = $3 AND friend_id = $4)`)).
		WithArgs(userID, friendID, friendID, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.Mock.ExpectCommit()

	err := repo.RemoveFriend(context.Background(), userID, friendID)

	assert.NoError(t, err)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestUpdateFriendStatus(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewUserRepository(mock.DB)

	userID := uuid.New()
	friendID := uuid.New()
	status := "ACCEPTED"

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectExec(regexp.QuoteMeta(`UPDATE "friend_lists" SET "status"=$1 WHERE user_id = $2 AND friend_id = $3`)).
		WithArgs(status, friendID, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.Mock.ExpectCommit()

	err := repo.UpdateFriendStatus(context.Background(), userID, friendID, status)

	assert.NoError(t, err)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestCheckFriendship(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewUserRepository(mock.DB)

	userID := uuid.New()
	friendID := uuid.New()

	mock.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "friend_lists" WHERE user_id = $1 AND friend_id = $2 ORDER BY "friend_lists"."user_id" LIMIT $3`)).
		WithArgs(userID, friendID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"status"}).AddRow("PENDING"))

	exists, status, err := repo.CheckFriendship(context.Background(), userID, friendID)

	assert.NoError(t, err)
	assert.True(t, exists)
	assert.Equal(t, "PENDING", status)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestGetIncomingRequests(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewUserRepository(mock.DB)

	userID := uuid.New()
	requesterID := uuid.New()

	mock.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "friend_lists" WHERE friend_id = $1 AND status = $2`)).
		WithArgs(userID, "PENDING").
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "friend_id"}).AddRow(requesterID, userID))

	mock.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)).
		WithArgs(requesterID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(requesterID))

	res, err := repo.GetIncomingRequests(context.Background(), userID)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestSetCloseFriendStatus(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewUserRepository(mock.DB)

	userID := uuid.New()
	friendID := uuid.New()

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectExec(regexp.QuoteMeta(`UPDATE "friend_lists" SET "is_close_friend"=$1 WHERE user_id = $2 AND friend_id = $3`)).
		WithArgs(true, userID, friendID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.Mock.ExpectCommit()

	err := repo.SetCloseFriendStatus(context.Background(), userID, friendID, true)

	assert.NoError(t, err)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestGetCloseFriends(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewUserRepository(mock.DB)

	userID := uuid.New()
	friendID := uuid.New()
	queryRegex := `^SELECT .* FROM "users" JOIN "friend_lists" ON "friend_lists"\."friend_id" = "users"\."id" AND "friend_lists"\."user_id" = \$1 WHERE friend_lists\.is_close_friend = \$2$`

	mock.Mock.ExpectQuery(queryRegex).
		WithArgs(userID, true).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(friendID))

	res, err := repo.GetCloseFriends(context.Background(), userID)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}
