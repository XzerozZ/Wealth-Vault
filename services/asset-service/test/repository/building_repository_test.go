package repository_test

import (
	"context"
	"database/sql/driver"
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

func TestCreateBuilding(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewBuildingRepository(mock.DB)

	building := &domain.Building{
		ID:   uuid.New(),
		Name: "Test Building",
	}

	mock.Mock.ExpectBegin()

	mock.Mock.ExpectQuery(`INSERT INTO "buildings"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(building.ID),
		)

	mock.Mock.ExpectCommit()

	err := repo.CreateBuilding(context.Background(), building)

	assert.NoError(t, err)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestGetBuilding(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewBuildingRepository(mock.DB)

	uid := uuid.New()
	bid := uuid.New()

	buildingRows := sqlmock.NewRows([]string{
		"id", "user_id",
	}).AddRow(bid, uid)

	locationRows := sqlmock.NewRows([]string{
		"id", "building_id",
	}).AddRow(uuid.New(), bid)

	mock.Mock.ExpectQuery(`SELECT \* FROM "buildings"`).
		WithArgs(uid).
		WillReturnRows(buildingRows)

	mock.Mock.ExpectQuery(`SELECT \* FROM "locations"`).
		WillReturnRows(locationRows)

	result, err := repo.GetBuilding(context.Background(), uid)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestGetBuildingByIDs(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewBuildingRepository(mock.DB)

	id := uuid.New()

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(id)

	mock.Mock.ExpectQuery(`SELECT .* FROM "buildings" WHERE id IN`).
		WillReturnRows(rows)

	items, err := repo.GetBuildingByIDs(context.Background(), []uuid.UUID{id})

	require.NoError(t, err)
	require.Len(t, items, 1)
	mock.ExpectDone(t)
}

func TestGetBuildingByID(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewBuildingRepository(mock.DB)

	buildingID := uuid.New()
	insuranceID := uuid.New()
	landID := uuid.New()
	locationID := uuid.New()

	mock.Mock.MatchExpectationsInOrder(false)

	mock.Mock.ExpectQuery(`SELECT .* FROM "buildings"`).
		WithArgs(buildingID, 1).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"user_id",
				"type",
				"name",
				"area",
				"amount",
				"description",
				"location_id",
			}).
				AddRow(buildingID, uuid.New(), "house", "Example", 100, 1000, "", locationID),
		)

	mock.Mock.ExpectQuery(`SELECT .* FROM "file_associates"`).
		WithArgs("building", buildingID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	mock.Mock.ExpectQuery(`SELECT .* FROM "building_insurance"`).
		WithArgs(buildingID).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"house_id",
				"ins_id",
			}).
				AddRow(buildingID, insuranceID),
		)

	mock.Mock.ExpectQuery(`SELECT .* FROM "insurances"`).
		WithArgs(insuranceID).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"name",
			}).
				AddRow(insuranceID, "Test Insurance"),
		)

	mock.Mock.ExpectQuery(`SELECT .* FROM "building_land"`).
		WithArgs(buildingID).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"house_id",
				"land_id",
			}).
				AddRow(buildingID, landID),
		)

	mock.Mock.ExpectQuery(`SELECT .* FROM "lands"`).
		WithArgs(landID).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"user_id",
				"name",
				"location_id",
			}).
				AddRow(landID, uuid.Nil, "Test Land", locationID),
		)

	mock.Mock.ExpectQuery(`SELECT .* FROM "locations"`).
		WithArgs(locationID).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"address",
			}).
				AddRow(locationID, "Test Address"),
		)

	item, err := repo.GetBuildingByID(context.Background(), buildingID)

	require.NoError(t, err)
	require.Equal(t, buildingID, item.ID)

	mock.ExpectDone(t)
}

func TestUpdateBuilding(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewBuildingRepository(mock.DB)

	buildingID := uuid.New()
	locationID := uuid.New()
	landID := uuid.New()
	insID := uuid.New()

	item := &domain.Building{
		ID:         buildingID,
		Name:       "Updated Name",
		LocationID: locationID,
		Location: domain.Location{
			ID: locationID,
		},
	}

	mock.Mock.MatchExpectationsInOrder(false)
	mock.Mock.ExpectBegin()
	mock.Mock.ExpectQuery("INSERT INTO \"locations\"").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), locationID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(locationID))

	buildingArgs := make([]driver.Value, 11)
	for i := range buildingArgs {
		buildingArgs[i] = sqlmock.AnyArg()
	}
	mock.Mock.ExpectExec("UPDATE \"buildings\"").
		WithArgs(buildingArgs...).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectExec("UPDATE \"locations\"").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), locationID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectQuery("SELECT .* FROM \"buildings\"").
		WithArgs(buildingID, buildingID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(buildingID, "Updated Name"))

	mock.Mock.ExpectQuery("SELECT .* FROM \"file_associates\"").
		WithArgs("building", buildingID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "entity_id"}).AddRow(uuid.New(), buildingID))

	mock.Mock.ExpectQuery("SELECT .* FROM \"locations\"").
		WithArgs(locationID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(locationID))

	mock.Mock.ExpectQuery("SELECT .* FROM \"building_land\"").
		WithArgs(buildingID).
		WillReturnRows(sqlmock.NewRows([]string{"house_id", "land_id"}).AddRow(buildingID, landID))

	mock.Mock.ExpectQuery("SELECT .* FROM \"lands\"").
		WithArgs(landID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(landID))

	mock.Mock.ExpectQuery("SELECT .* FROM \"building_insurance\"").
		WithArgs(buildingID).
		WillReturnRows(sqlmock.NewRows([]string{"house_id", "ins_id"}).AddRow(buildingID, insID))

	mock.Mock.ExpectQuery("SELECT .* FROM \"insurances\"").
		WithArgs(insID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(insID))

	mock.Mock.ExpectCommit()

	_, err := repo.UpdateBuilding(context.Background(), item, nil, nil, nil, nil)

	assert.NoError(t, err)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestSoftDeleteBuilding(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewBuildingRepository(mock.DB)

	id := uuid.New()
	uid := uuid.New()

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectExec(`UPDATE "buildings"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.Mock.ExpectCommit()

	err := repo.SoftDeleteBuilding(context.Background(), id, uid)

	assert.NoError(t, err)
}

func TestGetExpiredBuilding(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewBuildingRepository(mock.DB)
	olderThan := time.Now()
	buildingID := uuid.New()

	mock.Mock.ExpectQuery("SELECT .* FROM \"buildings\"").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "deleted_at"}).
			AddRow(buildingID, "Old Building", time.Now()))

	mock.Mock.ExpectQuery("SELECT .* FROM \"file_associates\"").
		WithArgs("building", buildingID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "entity_id", "file_name"}).
			AddRow(uuid.New(), buildingID, "document.pdf"))

	buildings, err := repo.GetExpiredBuilding(context.Background(), olderThan)

	assert.NoError(t, err)
	assert.Len(t, buildings, 1)
	assert.Len(t, buildings[0].Files, 1)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestHardDeleteBuilding(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewBuildingRepository(mock.DB)

	id := uuid.New()
	locationID := uuid.New()

	mock.Mock.ExpectBegin()

	mock.Mock.ExpectQuery("SELECT .* FROM \"buildings\"").
		WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "location_id"}).AddRow(id, locationID))

	mock.Mock.ExpectExec("DELETE FROM \"building_land\"").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectExec("DELETE FROM \"building_insurance\"").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectExec("DELETE FROM \"file_associates\"").
		WithArgs(id, "building").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectExec("DELETE FROM \"buildings\"").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectExec("DELETE FROM \"locations\"").
		WithArgs(locationID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectCommit()

	err := repo.HardDeleteBuilding(context.Background(), id)

	assert.NoError(t, err)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}
