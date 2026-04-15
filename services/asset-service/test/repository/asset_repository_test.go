package repository_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"
	"wealth-vault/asset-service/internal/repository"
	testutil "wealth-vault/asset-service/test/testdb"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAssetName(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewAssetRepository(mock.DB)

	assetID := uuid.New()
	userID := uuid.New()
	expectedName := "My Savings Account"

	t.Run("GetAssetName - Success (Account)", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow(expectedName)
		mock.Mock.ExpectQuery(`(?i)SELECT "name" FROM "accounts" WHERE \(id = \$1 AND user_id = \$2\).*`).
			WithArgs(assetID, userID).
			WillReturnRows(rows)

		name, exists, err := repo.CheckExists(context.Background(), "account", assetID, userID)

		assert.NoError(t, err)
		assert.True(t, exists)
		assert.Equal(t, expectedName, name)
	})

	t.Run("GetAssetName - Not Found", func(t *testing.T) {
		mock.Mock.ExpectQuery(`(?i)SELECT "name" FROM "investments" WHERE .*`).
			WithArgs(assetID, userID).
			WillReturnRows(sqlmock.NewRows([]string{"name"}))

		name, exists, err := repo.CheckExists(context.Background(), "investment", assetID, userID)

		assert.NoError(t, err)
		assert.False(t, exists)
		assert.Equal(t, "", name)
	})

	t.Run("GetAssetName - Database Error", func(t *testing.T) {
		mock.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT name FROM`)).
			WillReturnError(fmt.Errorf("connection refused"))

		name, exists, err := repo.CheckExists(context.Background(), "cash", assetID, userID)

		assert.Error(t, err)
		assert.False(t, exists)
		assert.Equal(t, "", name)
	})
}

func TestGetAllAssets_Success(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewAssetRepository(mock.DB)
	uid := uuid.New()

	assetRows := sqlmock.NewRows([]string{
		"id", "type", "name", "value", "created_at",
	}).AddRow(
		uuid.New(), "account", "My Account", 1000, time.Now(),
	)

	mock.Mock.ExpectQuery(`SELECT id, 'account'`).
		WithArgs(uid, uid, uid, uid, uid).
		WillReturnRows(assetRows)

	liabilityRows := sqlmock.NewRows([]string{
		"id", "type", "name", "value", "created_at",
	}).AddRow(
		uuid.New(), "liability", "Car Loan", 500, time.Now(),
	)

	mock.Mock.ExpectQuery(`SELECT id, 'liability'`).
		WithArgs(uid).
		WillReturnRows(liabilityRows)

	assets, liabilities, err := repo.GetAllAssets(context.Background(), uid)

	require.NoError(t, err)
	assert.Len(t, assets, 1)
	assert.Len(t, liabilities, 1)
	assert.Equal(t, "My Account", assets[0].Name)
	assert.Equal(t, "Car Loan", liabilities[0].Name)

	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestGetAllAssets_Empty(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewAssetRepository(mock.DB)
	uid := uuid.New()

	emptyAssetRows := sqlmock.NewRows([]string{
		"id", "type", "name", "value", "created_at",
	})

	mock.Mock.ExpectQuery(`SELECT id, 'account'`).
		WithArgs(uid, uid, uid, uid, uid).
		WillReturnRows(emptyAssetRows)

	emptyLiabilityRows := sqlmock.NewRows([]string{
		"id", "type", "name", "value", "created_at",
	})

	mock.Mock.ExpectQuery(`SELECT id, 'liability'`).
		WithArgs(uid).
		WillReturnRows(emptyLiabilityRows)

	assets, liabilities, err := repo.GetAllAssets(context.Background(), uid)

	require.NoError(t, err)
	assert.Len(t, assets, 0)
	assert.Len(t, liabilities, 0)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestGetAssetCount_Zero(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewAssetRepository(mock.DB)

	uid := uuid.New()

	rows := sqlmock.NewRows([]string{"total_count"}).AddRow(0)

	mock.Mock.ExpectQuery(`SELECT \(`).
		WithArgs(uid, uid, uid, uid, uid, uid).
		WillReturnRows(rows)

	count, err := repo.GetAssetCount(context.Background(), uid)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestGetAssetCount_Success(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewAssetRepository(mock.DB)

	uid := uuid.New()

	rows := sqlmock.NewRows([]string{"total_count"}).AddRow(5)

	mock.Mock.ExpectQuery(`SELECT \(`).
		WithArgs(uid, uid, uid, uid, uid, uid).
		WillReturnRows(rows)

	count, err := repo.GetAssetCount(context.Background(), uid)

	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestGetNetWorthOverview_Success(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewAssetRepository(mock.DB)

	uid := uuid.New()

	rows := sqlmock.NewRows([]string{
		"total_assets",
		"total_liabilities",
	}).AddRow(10000, 2000)

	mock.Mock.ExpectQuery(`SELECT`).
		WithArgs(uid, uid, uid, uid, uid, uid).
		WillReturnRows(rows)

	result, err := repo.GetNetWorthOverview(context.Background(), uid)

	assert.NoError(t, err)
	assert.Equal(t, float64(10000), result.TotalAssets)
	assert.Equal(t, float64(2000), result.TotalLiabilities)
}

func TestGetNetWorthOverview_Zero(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewAssetRepository(mock.DB)

	uid := uuid.New()

	rows := sqlmock.NewRows([]string{
		"total_assets",
		"total_liabilities",
	}).AddRow(0, 0)

	mock.Mock.ExpectQuery(`SELECT`).
		WithArgs(uid, uid, uid, uid, uid, uid).
		WillReturnRows(rows)

	result, err := repo.GetNetWorthOverview(context.Background(), uid)

	assert.NoError(t, err)
	assert.Equal(t, float64(0), result.TotalAssets)
	assert.Equal(t, float64(0), result.TotalLiabilities)
}
