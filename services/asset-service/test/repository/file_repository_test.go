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

func TestDeleteFiles(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewFileRepository(mock.DB)

	ctx := context.Background()

	fileID := uuid.New()

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "file_associates" WHERE id IN ($1)`,
	)).
		WithArgs(fileID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.Mock.ExpectCommit()

	err := repo.DeleteFiles(ctx, []uuid.UUID{fileID})

	require.NoError(t, err)
	require.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestCreateFiles(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewFileRepository(mock.DB)

	ctx := context.Background()

	entityID := uuid.New()
	userID := uuid.New()

	files := []domain.FileAssociate{
		{
			EntityID:   entityID,
			EntityType: "building",
			Link:       "https://example.com",
			FileType:   "image",
			UserID:     userID,
		},
	}

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "file_associates" ("entity_id","entity_type","link","file_type","user_id","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`,
	)).
		WithArgs(
			entityID,
			"building",
			"https://example.com",
			"image",
			userID,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(uuid.New()),
		)
	mock.Mock.ExpectCommit()

	err := repo.CreateFiles(ctx, files)

	require.NoError(t, err)
	require.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestGetFilesByIDs(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewFileRepository(mock.DB)

	ctx := context.Background()

	fileID := uuid.New()
	entityID := uuid.New()
	userID := uuid.New()

	rows := sqlmock.NewRows([]string{
		"id",
		"entity_id",
		"entity_type",
		"link",
		"file_type",
		"user_id",
		"created_at",
		"updated_at",
	}).AddRow(
		fileID,
		entityID,
		"building",
		"https://example.com",
		"image",
		userID,
		time.Now(),
		time.Now(),
	)

	mock.Mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "file_associates" WHERE id IN ($1)`,
	)).
		WithArgs(fileID).
		WillReturnRows(rows)

	files, err := repo.GetFilesByIDs(ctx, []uuid.UUID{fileID})

	require.NoError(t, err)
	require.Len(t, files, 1)
	require.Equal(t, fileID, files[0].ID)
	require.NoError(t, mock.Mock.ExpectationsWereMet())
}
