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

func TestCreateLand(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewLandRepository(mock.DB)

	land := &domain.Land{
		ID:     uuid.New(),
		UserID: uuid.New(),
		Name:   "Test Land",
	}

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "lands"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(land.ID))
	mock.Mock.ExpectCommit()

	err := repo.CreateLand(context.Background(), land)

	require.NoError(t, err)
	mock.ExpectDone(t)
}

func TestGetLand(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewLandRepository(mock.DB)

	userID := uuid.New()
	locationID := uuid.New()

	mock.Mock.MatchExpectationsInOrder(false)

	mock.Mock.ExpectQuery(`SELECT .* FROM "lands"`).
		WithArgs(userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "user_id", "location_id"}).
				AddRow(uuid.New(), userID, locationID),
		)

	mock.Mock.ExpectQuery(`SELECT .* FROM "locations"`).
		WithArgs(locationID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(locationID),
		)

	items, err := repo.GetLand(context.Background(), userID)

	require.NoError(t, err)
	require.Len(t, items, 1)
	mock.ExpectDone(t)
}

func TestGetLandByID(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewLandRepository(mock.DB)

	landID := uuid.New()
	locationID := uuid.New()

	mock.Mock.MatchExpectationsInOrder(false)

	// main query
	mock.Mock.ExpectQuery(`SELECT .* FROM "lands"`).
		WithArgs(landID, 1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "location_id"}).
				AddRow(landID, locationID),
		)

	// preload Files
	mock.Mock.ExpectQuery(`SELECT .* FROM "file_associates"`).
		WithArgs("land", landID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	// preload Location
	mock.Mock.ExpectQuery(`SELECT .* FROM "locations"`).
		WithArgs(locationID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(locationID),
		)

	// preload Buildings (many2many join table)
	mock.Mock.ExpectQuery(`SELECT .* FROM "building_land"`).
		WithArgs(landID).
		WillReturnRows(
			sqlmock.NewRows([]string{"land_id", "house_id"}),
		)

	item, err := repo.GetLandByID(context.Background(), landID)

	require.NoError(t, err)
	require.Equal(t, landID, item.ID)
	mock.ExpectDone(t)
}

func TestSoftDeleteLand(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewLandRepository(mock.DB)

	id := uuid.New()
	uid := uuid.New()

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectExec(`UPDATE "lands"`).
		WithArgs(sqlmock.AnyArg(), id, uid).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectCommit()

	err := repo.SoftDeleteLand(context.Background(), id, uid)

	require.NoError(t, err)
	mock.ExpectDone(t)
}

func TestHardDeleteLand(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewLandRepository(mock.DB)

	id := uuid.New()
	locationID := uuid.New()

	mock.Mock.ExpectBegin()

	mock.Mock.ExpectQuery(`SELECT .* FROM "lands"`).
		WithArgs(id, 1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "location_id"}).
				AddRow(id, locationID),
		)

	mock.Mock.ExpectExec(`DELETE FROM "building_land"`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectExec(`DELETE FROM "file_associates"`).
		WithArgs(id, "land").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectExec(`DELETE FROM "lands"`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectCommit()

	err := repo.HardDeleteLand(context.Background(), id)

	require.NoError(t, err)
	mock.ExpectDone(t)
}

func TestGetExpiredLand(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewLandRepository(mock.DB)

	refTime := time.Now()

	mock.Mock.MatchExpectationsInOrder(false)

	mock.Mock.ExpectQuery(`SELECT .* FROM "lands"`).
		WithArgs(refTime).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(uuid.New()),
		)

	mock.Mock.ExpectQuery(`SELECT .* FROM "file_associates"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	items, err := repo.GetExpiredLand(context.Background(), refTime)

	require.NoError(t, err)
	require.Len(t, items, 1)
	mock.ExpectDone(t)
}
