package repository_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"wealth-vault/asset-service/internal/domain" // ปรับ Import path ตามโปรเจคจริงของคุณ
	"wealth-vault/asset-service/internal/repository"
	testutil "wealth-vault/asset-service/test/testdb"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateCash(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewCashRepository(mock.DB)

	cash := &domain.Cash{
		UserID: uuid.New(),
		Name:   "Wallet",
		Amount: 1000,
	}

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "cashes"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
	mock.Mock.ExpectCommit()

	err := repo.CreateCash(context.Background(), cash)

	assert.NoError(t, err)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestGetCash(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewCashRepository(mock.DB)
	uid := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "amount"}).
		AddRow(uuid.New(), uid, "Cash 1", 100).
		AddRow(uuid.New(), uid, "Cash 2", 200)

	mock.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cashes" WHERE user_id = $1`)).
		WithArgs(uid).
		WillReturnRows(rows)

	items, err := repo.GetCash(context.Background(), uid)

	assert.NoError(t, err)
	assert.Len(t, items, 2)
	assert.Equal(t, "Cash 1", items[0].Name)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestGetCashByID(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewCashRepository(mock.DB)

	id := uuid.New()

	mock.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cashes" WHERE id = $1 AND "cashes"."deleted_at" IS NULL ORDER BY "cashes"."id" LIMIT $2`)).
		WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(id, "Test Cash"))

	mock.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "file_associates" WHERE "entity_type" = $1 AND "file_associates"."entity_id" = $2`)).
		WithArgs("cash", id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "file_name", "entity_id", "entity_type"}).
			AddRow(uuid.New(), "slip.jpg", id, "cash"))

	item, err := repo.GetCashByID(context.Background(), id)

	assert.NoError(t, err)

	if assert.NotNil(t, item) {
		assert.Equal(t, id, item.ID)
		assert.Len(t, item.Files, 1)
	}

	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestGetCashByIDs(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewCashRepository(mock.DB)

	ids := []uuid.UUID{uuid.New(), uuid.New()}

	mock.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cashes" WHERE id IN ($1,$2)`)).
		WithArgs(ids[0], ids[1]).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(ids[0]).AddRow(ids[1]))

	items, err := repo.GetCashByIDs(context.Background(), ids)

	assert.NoError(t, err)
	assert.Len(t, items, 2)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestUpdateCash(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewCashRepository(mock.DB)

	id := uuid.New()
	cash := &domain.Cash{
		ID:     id,
		Name:   "Updated Name",
		Amount: 5000,
	}

	mock.Mock.ExpectBegin()

	mock.Mock.ExpectExec(`UPDATE "cashes"`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectQuery(`SELECT \* FROM "cashes"`).
		WithArgs(id, id, 1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(id),
		)

	mock.Mock.ExpectQuery(`SELECT \* FROM "file_associates"`).
		WithArgs("cash", id).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "entity_type", "entity_id"}),
		)

	mock.Mock.ExpectCommit()

	result, err := repo.UpdateCash(context.Background(), cash)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestSoftDeleteCash(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewCashRepository(mock.DB)

	id := uuid.New()
	uid := uuid.New()

	mock.Mock.ExpectBegin()

	mock.Mock.ExpectExec(regexp.QuoteMeta(`UPDATE "cashes" SET "deleted_at"=$1 WHERE (id = $2 AND user_id = $3) AND "cashes"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), id, uid).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectCommit()

	err := repo.SoftDeleteCash(context.Background(), id, uid)

	assert.NoError(t, err)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestGetExpiredCash(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewCashRepository(mock.DB)

	olderThan := time.Now()

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(uuid.New())

	mock.Mock.ExpectQuery(`SELECT \* FROM "cashes"`).
		WithArgs(olderThan).
		WillReturnRows(rows)

	mock.Mock.ExpectQuery(`SELECT \* FROM "file_associates"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	items, err := repo.GetExpiredCash(context.Background(), olderThan)

	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestHardDeleteCash(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewCashRepository(mock.DB)

	id := uuid.New()

	mock.Mock.ExpectBegin()

	mock.Mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "file_associates" WHERE entity_id = $1 AND entity_type = $2`)).
		WithArgs(id, "cash").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "cashes" WHERE id = $1`)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectCommit()

	err := repo.HardDeleteCash(context.Background(), id)

	assert.NoError(t, err)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}
