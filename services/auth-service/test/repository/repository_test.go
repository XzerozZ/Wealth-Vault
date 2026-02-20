package repository_test

import (
	"context"
	"regexp"
	"testing"
	"wealth-vault/auth-service/internal/domain"
	"wealth-vault/auth-service/internal/repository"
	testutil "wealth-vault/auth-service/test/testdb"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewAuthRepository(mock.DB)

	auth := &domain.AuthAccount{UserID: uuid.New(), Email: "test@example.com"}

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "auth_accounts"`)).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(auth.UserID))
	mock.Mock.ExpectCommit()

	err := repo.Register(context.Background(), auth)

	assert.NoError(t, err)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestFindByEmailAndProvider(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewAuthRepository(mock.DB)

	email := "test@example.com"
	provider := "local"
	expectedUserID := uuid.New()

	mock.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "auth_accounts" WHERE email = $1 AND provider = $2`)).
		WithArgs(email, provider, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"email", "provider", "user_id"}).
			AddRow(email, provider, expectedUserID))

	res, err := repo.FindByEmailAndProvider(context.Background(), email, provider)

	assert.NoError(t, err)
	if assert.NotNil(t, res) {
		assert.Equal(t, email, res.Email)
		assert.Equal(t, provider, res.Provider)
		assert.Equal(t, expectedUserID, res.UserID)
	}

	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestCreateSession(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewAuthRepository(mock.DB)

	session := &domain.AuthSession{
		ID:           uuid.New(),
		RefreshToken: "token-123",
		UserID:       uuid.New(),
	}

	mock.Mock.ExpectBegin()

	mock.Mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "auth_sessions"`)).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(uuid.New()),
		)

	mock.Mock.ExpectCommit()

	err := repo.CreateSession(context.Background(), session)

	assert.NoError(t, err)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestRevokeSession(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewAuthRepository(mock.DB)

	token := "token-123"
	mock.Mock.ExpectBegin()
	mock.Mock.ExpectExec(regexp.QuoteMeta(`UPDATE "auth_sessions" SET "revoked"=$1`)).
		WithArgs(true, sqlmock.AnyArg(), token).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.Mock.ExpectCommit()

	err := repo.RevokeSession(context.Background(), token)

	assert.NoError(t, err)
}

func TestDeleteExpiredSessions(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewAuthRepository(mock.DB)

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "auth_sessions" WHERE refresh_expires_at < $1 OR (revoked = $2 AND updated_at < $3)`)).
		WithArgs(sqlmock.AnyArg(), true, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.Mock.ExpectCommit()

	err := repo.DeleteExpiredSessions(context.Background())

	assert.NoError(t, err)
}

func TestGetValidOTP(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewAuthRepository(mock.DB)

	uID := uuid.New()
	code := "123456"

	queryRegex := regexp.QuoteMeta(`SELECT * FROM "auth_otps" WHERE user_id = $1 AND otp = $2 AND expired_at > $3`) + `.*`

	mock.Mock.ExpectQuery(queryRegex).
		WithArgs(uID, code, sqlmock.AnyArg(), 1).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "otp"}).AddRow(uID, code))

	res, err := repo.GetValidOTP(context.Background(), uID, code)

	assert.NoError(t, err)
	if assert.NotNil(t, res) {
		assert.Equal(t, code, res.OTP)
	}

	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestUpdatePassword(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewAuthRepository(mock.DB)

	uID := uuid.New()
	newHash := "new-hashed-password"

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectExec(regexp.QuoteMeta(`UPDATE "auth_accounts" SET "password"=$1,"updated_at"=$2 WHERE user_id = $3`)).
		WithArgs(newHash, sqlmock.AnyArg(), uID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.Mock.ExpectCommit()

	err := repo.UpdatePassword(context.Background(), uID, newHash)

	assert.NoError(t, err)
}
