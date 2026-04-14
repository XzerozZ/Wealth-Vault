package repository_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"wealth-vault/user-service/internal/domain"
	"wealth-vault/user-service/internal/repository"
	testutil "wealth-vault/user-service/test/testdb"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewUserRepository(mock.DB)

	user := &domain.User{ID: uuid.New()}

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(user.ID))
	mock.Mock.ExpectCommit()

	err := repo.CreateUser(context.Background(), user)

	assert.NoError(t, err)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestGetUser(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewUserRepository(mock.DB)

	expectedID := uuid.New()

	mock.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(expectedID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	user, err := repo.GetUser(context.Background(), expectedID)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedID, user.ID)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestGetUsersByEmail(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewUserRepository(mock.DB)

	email := "test@gmail.com"
	userID := uuid.New()

	t.Run("GetUsersByEmail - Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "username"}).
			AddRow(userID, email, "testuser")

		mock.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1`)).
			WithArgs(email).
			WillReturnRows(rows)

		res, err := repo.GetUsersByEmail(context.Background(), email)

		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, email, res[0].Email)
		assert.Equal(t, userID, res[0].ID)
	})

	t.Run("GetUsersByEmail - Not Found", func(t *testing.T) {
		mock.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1`)).
			WithArgs("notfound@gmail.com").
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "username"}))

		res, err := repo.GetUsersByEmail(context.Background(), "notfound@gmail.com")

		assert.NoError(t, err)
		assert.Len(t, res, 0)
	})

	t.Run("GetUsersByEmail - DB Error", func(t *testing.T) {
		mock.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1`)).
			WithArgs(email).
			WillReturnError(fmt.Errorf("db connection error"))

		res, err := repo.GetUsersByEmail(context.Background(), email)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestUpdateUser(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewUserRepository(mock.DB)

	user := &domain.User{ID: uuid.New()}
	mask := []string{"Name"}

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.Mock.ExpectCommit()

	mock.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 AND "users"."id" = $2 ORDER BY "users"."id" LIMIT $3`)).
		WithArgs(user.ID, user.ID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(user.ID))

	updatedUser, err := repo.UpdateUser(context.Background(), user, mask)

	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestGetUsersReadyForAutoShare(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewUserRepository(mock.DB)

	userID := uuid.New()
	friendID := uuid.New()

	mock.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE (is_auto_share_enabled = $1 AND is_auto_share_triggered = $2) AND EXTRACT(YEAR FROM age(birthday)) >= auto_share_age`)).
		WithArgs(true, false).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))

	mock.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "friend_lists" WHERE "friend_lists"."user_id" = $1`)).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"friend_id", "user_id"}).AddRow(friendID, userID))

	preloadRegex := `^SELECT .* FROM "users" JOIN friend_lists ON friend_lists\.friend_id = users\.id WHERE "users"\."id" = \$1 AND \(friend_lists\.status = \$2 AND friend_lists\.is_close_friend = \$3\)$`

	mock.Mock.ExpectQuery(preloadRegex).
		WithArgs(friendID, "ACCEPTED", true).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(friendID))

	res, err := repo.GetUsersReadyForAutoShare(context.Background())

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestMarkAutoShareTriggered(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewUserRepository(mock.DB)

	userID := uuid.New()

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "is_auto_share_triggered"=$1,"updated_at"=$2 WHERE id = $3`)).
		WithArgs(true, sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.Mock.ExpectCommit()

	err := repo.MarkAutoShareTriggered(context.Background(), userID)

	assert.NoError(t, err)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestCreateFriendLog(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewUserRepository(mock.DB)

	log := &domain.FriendLog{ID: uuid.New()}

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "friend_logs"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(log.ID))
	mock.Mock.ExpectCommit()

	err := repo.CreateFriendLog(context.Background(), log)

	assert.NoError(t, err)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}
