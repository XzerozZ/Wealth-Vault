package repository_test

import (
	"context"
	"regexp"
	"testing"
	"time"
	"wealth-vault/asset-service/internal/domain"
	"wealth-vault/asset-service/internal/repository"
	testutil "wealth-vault/asset-service/test/testdb"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewAccountRepository(mock.DB)

	account := &domain.Account{
		ID:     uuid.New(),
		UserID: uuid.New(),
		Name:   "Test Account",
	}

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "accounts"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(account.ID))
	mock.Mock.ExpectCommit()

	err := repo.CreateAccount(context.Background(), account)

	assert.NoError(t, err)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestGetAccount(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewAccountRepository(mock.DB)
	uid := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "user_id"}).
		AddRow(uuid.New(), uid)

	mock.Mock.ExpectQuery(`SELECT \* FROM "accounts"`).
		WithArgs(uid).
		WillReturnRows(rows)

	res, err := repo.GetAccount(context.Background(), uid)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestGetAccountByID(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewAccountRepository(mock.DB)
	id := uuid.New()

	accountRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(id, "Savings Account")

	mock.Mock.ExpectQuery(`SELECT \* FROM "accounts"`).
		WithArgs(id, 1).
		WillReturnRows(accountRows)

	mock.Mock.ExpectQuery(`SELECT \* FROM "file_associates"`).
		WithArgs("account", id).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "entity_type", "entity_id"}).
				AddRow(uuid.New(), "account", id),
		)

	res, err := repo.GetAccountByID(context.Background(), id)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, id, res.ID)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestGetAccountByIDs(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewAccountRepository(mock.DB)
	ids := []uuid.UUID{uuid.New(), uuid.New()}

	mock.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "accounts" WHERE id IN ($1,$2)`)).
		WithArgs(ids[0], ids[1]).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(ids[0]).AddRow(ids[1]))

	items, err := repo.GetAccountByIDs(context.Background(), ids)

	assert.NoError(t, err)
	assert.Len(t, items, 2)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestUpdateAccount(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewAccountRepository(mock.DB)
	id := uuid.New()

	account := &domain.Account{
		ID: id,
	}

	mock.Mock.ExpectBegin()

	mock.Mock.ExpectExec(`UPDATE "accounts"`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectQuery(`SELECT \* FROM "accounts"`).
		WithArgs(id, id, 1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(id),
		)

	mock.Mock.ExpectQuery(`SELECT \* FROM "file_associates"`).
		WithArgs("account", id).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "entity_type", "entity_id"}),
		)

	mock.Mock.ExpectCommit()

	res, err := repo.UpdateAccount(context.Background(), account)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestSoftDeleteAccount(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewAccountRepository(mock.DB)

	id := uuid.New()
	uid := uuid.New()

	mock.Mock.ExpectBegin()

	mock.Mock.ExpectExec(`UPDATE "accounts"`).
		WithArgs(sqlmock.AnyArg(), id, uid).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectCommit()

	err := repo.SoftDeleteAccount(context.Background(), id, uid)

	assert.NoError(t, err)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestGetExpiredAccounts(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewAccountRepository(mock.DB)

	timeArg := time.Now()

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(uuid.New())

	mock.Mock.ExpectQuery(`SELECT \* FROM "accounts"`).
		WithArgs(timeArg).
		WillReturnRows(rows)

	mock.Mock.ExpectQuery(`SELECT \* FROM "file_associates"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	res, err := repo.GetExpiredAccounts(context.Background(), timeArg)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestHardDeleteAccount(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewAccountRepository(mock.DB)
	id := uuid.New()

	mock.Mock.ExpectBegin()

	mock.Mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "file_associates" WHERE entity_id = $1 AND entity_type = $2`)).
		WithArgs(id, "account").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "accounts" WHERE id = $1`)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectCommit()

	err := repo.HardDeleteAccount(context.Background(), id)

	assert.NoError(t, err)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}
