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
	"github.com/stretchr/testify/require"
)

func TestCreateInsurance(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewInsuranceRepository(mock.DB)

	policy := &domain.Insurance{
		ID:     uuid.New(),
		UserID: uuid.New(),
	}

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectQuery(`INSERT INTO "insurances"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(policy.ID))
	mock.Mock.ExpectCommit()

	err := repo.CreateInsurance(context.Background(), policy)

	require.NoError(t, err)
	mock.ExpectDone(t)
}

func TestGetInsurance(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewInsuranceRepository(mock.DB)

	uid := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "user_id"}).
		AddRow(uuid.New(), uid)

	mock.Mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "insurances" WHERE user_id = $1 AND "insurances"."deleted_at" IS NULL`,
	)).
		WithArgs(uid).
		WillReturnRows(rows)

	items, err := repo.GetInsurance(context.Background(), uid)

	require.NoError(t, err)
	require.Len(t, items, 1)
	mock.ExpectDone(t)
}

func TestGetInsuranceByIDs(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewInsuranceRepository(mock.DB)

	id := uuid.New()

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(id)

	mock.Mock.ExpectQuery(`SELECT .* FROM "insurances" WHERE id IN`).
		WillReturnRows(rows)

	items, err := repo.GetInsuranceByIDs(context.Background(), []uuid.UUID{id})

	require.NoError(t, err)
	require.Len(t, items, 1)
	mock.ExpectDone(t)
}

func TestGetInsuranceByID(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewInsuranceRepository(mock.DB)

	id := uuid.New()

	insruanceRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(id, "Example Insurance")

	mock.Mock.ExpectQuery(`SELECT \* FROM "insurances"`).
		WithArgs(id, 1).
		WillReturnRows(insruanceRows)

	mock.Mock.ExpectQuery(`SELECT .* FROM "file_associates"`).
		WithArgs("insurance", id).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "entity_type", "entity_id"}).
				AddRow(uuid.New(), "insurance", id),
		)

	item, err := repo.GetInsuranceByID(context.Background(), id)

	require.NoError(t, err)
	require.Equal(t, id, item.ID)

	mock.ExpectDone(t)
}

func TestUpdateInsurance(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewInsuranceRepository(mock.DB)
	id := uuid.New()

	policy := &domain.Insurance{
		ID: id,
	}

	mock.Mock.ExpectBegin()

	mock.Mock.ExpectExec(`UPDATE "insurances"`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectQuery(`SELECT .* FROM "insurances"`).
		WithArgs(id, id, 1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(id),
		)

	mock.Mock.ExpectQuery(`SELECT \* FROM "file_associates"`).
		WithArgs("insurance", id).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "entity_type", "entity_id"}),
		)

	mock.Mock.ExpectCommit()

	_, err := repo.UpdateInsurance(context.Background(), policy)

	require.NoError(t, err)
	mock.ExpectDone(t)
}

func TestSoftDeleteInsurance(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewInsuranceRepository(mock.DB)

	id := uuid.New()
	uid := uuid.New()

	mock.Mock.ExpectBegin()

	mock.Mock.ExpectExec(`UPDATE "insurances"`).
		WithArgs(sqlmock.AnyArg(), id, uid).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectCommit()

	err := repo.SoftDeleteInsurances(context.Background(), id, uid)

	assert.NoError(t, err)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestGetExpiredInsurances(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewInsuranceRepository(mock.DB)

	id := uuid.New()

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(id)

	mock.Mock.ExpectQuery(`SELECT .* FROM "insurances" WHERE deleted_at IS NOT NULL`).
		WillReturnRows(rows)

	mock.Mock.ExpectQuery(`SELECT .* FROM "file_associates"`).
		WithArgs("insurance", id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "entity_type", "entity_id"}))

	items, err := repo.GetExpiredInsurances(context.Background(), time.Now())

	require.NoError(t, err)
	require.Len(t, items, 1)
	mock.ExpectDone(t)
}

func TestHardDeleteInsurances(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewInsuranceRepository(mock.DB)

	id := uuid.New()

	mock.Mock.ExpectBegin()

	mock.Mock.ExpectQuery(`SELECT .* FROM "insurances"`).
		WithArgs(id, 1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(id),
		)

	mock.Mock.ExpectExec("DELETE FROM \"building_insurance\"").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectExec("DELETE FROM \"file_associates\"").
		WithArgs(id, "insurance").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectExec(`DELETE FROM "insurances"`).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectCommit()

	err := repo.HardDeleteInsurances(context.Background(), id)

	require.NoError(t, err)
	mock.ExpectDone(t)
}

func TestGetExpiringInsurances(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewInsuranceRepository(mock.DB)

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(uuid.New())

	mock.Mock.ExpectQuery(`SELECT .* FROM "insurances" WHERE DATE\(exp_date\)`).
		WillReturnRows(rows)

	items, err := repo.GetExpiringInsurances(context.Background(), 5)

	require.NoError(t, err)
	require.Len(t, items, 1)
	mock.ExpectDone(t)
}
