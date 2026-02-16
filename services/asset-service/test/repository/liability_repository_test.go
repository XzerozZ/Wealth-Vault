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
	"github.com/stretchr/testify/require"
)

func TestCreateLiability(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewLiabilityRepository(mock.DB)

	lia := &domain.Liability{
		ID:     uuid.New(),
		UserID: uuid.New(),
	}

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "liabilities"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(lia.ID))
	mock.Mock.ExpectCommit()

	err := repo.CreateLiability(context.Background(), lia)
	require.NoError(t, err)

	mock.ExpectDone(t)
}

func TestGetLiability(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewLiabilityRepository(mock.DB)

	uid := uuid.New()

	mock.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "liabilities" WHERE user_id = $1`)).
		WithArgs(uid).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "user_id"}).
				AddRow(uuid.New(), uid),
		)

	result, err := repo.GetLiability(context.Background(), uid)

	require.NoError(t, err)
	require.Len(t, result, 1)

	mock.ExpectDone(t)
}

func TestGetLiabilityByID(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewLiabilityRepository(mock.DB)

	id := uuid.New()

	liaRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(id, "Example Loan")

	mock.Mock.ExpectQuery(`SELECT \* FROM "liabilities"`).
		WithArgs(id, 1).
		WillReturnRows(liaRows)

	mock.Mock.ExpectQuery(`SELECT \* FROM "file_associates"`).
		WithArgs("liability", id).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "entity_type", "entity_id"}).
				AddRow(uuid.New(), "liability", id),
		)

	result, err := repo.GetLiabilityByID(context.Background(), id)

	require.NoError(t, err)
	require.Equal(t, id, result.ID)

	mock.ExpectDone(t)
}

func TestUpdateLiability(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewLiabilityRepository(mock.DB)
	id := uuid.New()
	lia := &domain.Liability{
		ID: id,
	}

	mock.Mock.ExpectBegin()

	mock.Mock.ExpectExec(`UPDATE "liabilities"`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectQuery(`SELECT .* FROM "liabilities"`).
		WithArgs(id, id, 1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(lia.ID),
		)

	mock.Mock.ExpectQuery(`SELECT \* FROM "file_associates"`).
		WithArgs("liability", id).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "entity_type", "entity_id"}),
		)

	mock.Mock.ExpectCommit()

	result, err := repo.UpdateLiability(context.Background(), lia)

	require.NoError(t, err)
	require.Equal(t, lia.ID, result.ID)

	mock.ExpectDone(t)
}

func TestSoftDeleteLiability(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewLiabilityRepository(mock.DB)

	id := uuid.New()
	uid := uuid.New()

	mock.Mock.ExpectBegin()

	mock.Mock.ExpectExec(`UPDATE "liabilities"`).
		WithArgs(sqlmock.AnyArg(), id, uid).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectCommit()

	err := repo.SoftDeleteLiability(context.Background(), id, uid)

	require.NoError(t, err)
	mock.ExpectDone(t)
}

func TestHardDeleteLiability(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewLiabilityRepository(mock.DB)

	id := uuid.New()

	mock.Mock.ExpectBegin()

	mock.Mock.ExpectExec(`DELETE FROM "file_associates"`).
		WithArgs(id, "liability").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectExec(`DELETE FROM "liabilities"`).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectCommit()

	err := repo.HardDeleteLiability(context.Background(), id)

	require.NoError(t, err)
	mock.ExpectDone(t)
}

func TestGetExpiredLiability(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewLiabilityRepository(mock.DB)

	refTime := time.Now()

	mock.Mock.ExpectQuery(`SELECT .* FROM "liabilities"`).
		WithArgs(refTime).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(uuid.New()),
		)

	mock.Mock.ExpectQuery(`SELECT \* FROM "file_associates"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	result, err := repo.GetExpiredLiability(context.Background(), refTime)

	require.NoError(t, err)
	require.Len(t, result, 1)

	mock.ExpectDone(t)
}
