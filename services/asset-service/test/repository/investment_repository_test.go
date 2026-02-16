package repository_test

import (
	"context"
	"testing"
	"time"
	"wealth-vault/asset-service/internal/domain"
	"wealth-vault/asset-service/internal/repository"
	testutil "wealth-vault/asset-service/test/testdb"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateInvestment(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewInvestmentRepository(mock.DB)

	invest := &domain.Investment{
		ID:     uuid.New(),
		UserID: uuid.New(),
		Name:   "Test",
	}

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectQuery(`INSERT INTO "investments"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(invest.ID))
	mock.Mock.ExpectCommit()

	err := repo.CreateInvestment(context.Background(), invest)

	require.NoError(t, err)
	mock.ExpectDone(t)
}

func TestGetInvestment(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewInvestmentRepository(mock.DB)

	userID := uuid.New()

	mock.Mock.ExpectQuery(`SELECT .* FROM "investments"`).
		WithArgs(userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "user_id"}).
				AddRow(uuid.New(), userID),
		)

	items, err := repo.GetInvestment(context.Background(), userID)

	require.NoError(t, err)
	require.Len(t, items, 1)
	mock.ExpectDone(t)
}

func TestGetInvestmentByIDs(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewInvestmentRepository(mock.DB)

	id := uuid.New()

	mock.Mock.ExpectQuery(`SELECT .* FROM "investments"`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(id),
		)

	items, err := repo.GetInvestmentByIDs(context.Background(), []uuid.UUID{id})

	require.NoError(t, err)
	require.Len(t, items, 1)
	mock.ExpectDone(t)
}

func TestGetInvestmentByID(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewInvestmentRepository(mock.DB)

	id := uuid.New()

	mock.Mock.MatchExpectationsInOrder(false)

	mock.Mock.ExpectQuery(`SELECT .* FROM "investments"`).
		WithArgs(id, 1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(id),
		)

	mock.Mock.ExpectQuery(`SELECT .* FROM "file_associates"`).
		WithArgs("investment", id).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}),
		)

	item, err := repo.GetInvestmentByID(context.Background(), id)

	require.NoError(t, err)
	require.Equal(t, id, item.ID)
	mock.ExpectDone(t)
}

func TestUpdateInvestment(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewInvestmentRepository(mock.DB)

	id := uuid.New()
	invest := &domain.Investment{
		ID: id,
	}

	mock.Mock.ExpectBegin()

	mock.Mock.ExpectExec(`UPDATE "investments"`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectQuery(`SELECT .* FROM "investments"`).
		WithArgs(id, id, 1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(invest.ID),
		)

	mock.Mock.ExpectQuery(`SELECT .* FROM "file_associates"`).
		WithArgs("investment", invest.ID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "entity_type", "entity_id"}),
		)

	mock.Mock.ExpectCommit()

	_, err := repo.UpdateInvestment(context.Background(), invest)

	require.NoError(t, err)
	mock.ExpectDone(t)
}

func TestSoftDeleteInvestment(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewInvestmentRepository(mock.DB)

	id := uuid.New()
	uid := uuid.New()

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectExec(`UPDATE "investments"`).
		WithArgs(sqlmock.AnyArg(), id, uid).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectCommit()

	err := repo.SoftDeleteInvestment(context.Background(), id, uid)

	require.NoError(t, err)
	mock.ExpectDone(t)
}

func TestHardDeleteInvestment(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewInvestmentRepository(mock.DB)

	id := uuid.New()

	mock.Mock.ExpectBegin()

	mock.Mock.ExpectExec(`DELETE FROM "file_associates"`).
		WithArgs(id, "investment").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectExec(`DELETE FROM "investments"`).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectCommit()

	err := repo.HardDeleteInvestment(context.Background(), id)

	require.NoError(t, err)
	mock.ExpectDone(t)
}

func TestGetExpiredInvestment(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewInvestmentRepository(mock.DB)

	timeRef := time.Now()

	mock.Mock.MatchExpectationsInOrder(false)

	mock.Mock.ExpectQuery(`SELECT .* FROM "investments"`).
		WithArgs(timeRef).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(uuid.New()),
		)

	mock.Mock.ExpectQuery(`SELECT .* FROM "file_associates"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	items, err := repo.GetExpiredInvestment(context.Background(), timeRef)

	require.NoError(t, err)
	require.Len(t, items, 1)
	mock.ExpectDone(t)
}
