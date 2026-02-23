package repository_test

import (
	"context"
	"testing"
	"wealth-vault/user-service/internal/domain"
	"wealth-vault/user-service/internal/repository"
	testutil "wealth-vault/user-service/test/testdb"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateGroup(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewGroupRepository(mock.DB)

	groupID := uuid.New()
	creatorID := uuid.New()
	memberID := uuid.New()

	group := &domain.Group{
		ID:        groupID,
		CreatedBy: creatorID,
	}
	memberIDs := []string{memberID.String(), creatorID.String(), "invalid-uuid"} // ใส่ creator กับ invalid เพื่อเทสเงื่อนไข continue

	mock.Mock.ExpectBegin()

	mock.Mock.ExpectQuery(`(?i)INSERT INTO "groups".*RETURNING`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(groupID))

	mock.Mock.ExpectQuery(`(?i)INSERT INTO "group_members".*RETURNING`).
		WillReturnRows(sqlmock.NewRows([]string{"group_id"}).AddRow(groupID))

	mock.Mock.ExpectCommit()

	err := repo.CreateGroup(context.Background(), group, memberIDs)

	assert.NoError(t, err)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestGetMember(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewGroupRepository(mock.DB)

	groupID := uuid.New()
	userID := uuid.New()

	countRegex := `(?i)SELECT count\(\*\) FROM "users" JOIN group_members ON group_members\.user_id = users\.id WHERE group_members\.group_id = \$1`
	mock.Mock.ExpectQuery(countRegex).
		WithArgs(groupID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

	findRegex := `(?i)SELECT .* FROM "users" JOIN group_members ON group_members\.user_id = users\.id WHERE group_members\.group_id = \$1 ORDER BY users\.username ASC`
	mock.Mock.ExpectQuery(findRegex).
		WithArgs(groupID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(userID, "testuser"))

	members, total, err := repo.GetMember(context.Background(), groupID)

	assert.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, members, 1)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestGetGroup(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewGroupRepository(mock.DB)

	groupID := uuid.New()

	mock.Mock.ExpectQuery(`(?i)SELECT count\(\*\) FROM "group_members" WHERE group_id = \$1`).
		WithArgs(groupID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))

	mock.Mock.ExpectQuery(`(?i)SELECT \* FROM "groups" WHERE id = \$1 ORDER BY "groups"\."id" LIMIT \$2`).
		WithArgs(groupID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(groupID))

	group, total, err := repo.GetGroup(context.Background(), groupID)

	assert.NoError(t, err)
	assert.Equal(t, int64(10), total)
	assert.NotNil(t, group)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestAllGetGroup(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewGroupRepository(mock.DB)

	uid := uuid.New()
	groupID := uuid.New()

	queryRegex := `(?i)SELECT g\.\*, \(SELECT COUNT\(\*\) FROM group_members WHERE group_id = g\.id\) as member_count FROM groups g JOIN group_members gm ON gm\.group_id = g\.id WHERE gm\.user_id = \$1`

	mock.Mock.ExpectQuery(queryRegex).
		WithArgs(uid).
		WillReturnRows(sqlmock.NewRows([]string{"id", "member_count"}).AddRow(groupID, 3))

	results, err := repo.AllGetGroup(context.Background(), uid)

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestIsUserMember(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewGroupRepository(mock.DB)

	groupID := uuid.New()
	userID := uuid.New()

	mock.Mock.ExpectQuery(`(?i)SELECT count\(\*\) FROM "group_members" WHERE group_id = \$1 AND user_id = \$2`).
		WithArgs(groupID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	isMember, err := repo.IsUserMember(context.Background(), groupID, userID)

	assert.NoError(t, err)
	assert.True(t, isMember)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestUpdateGroup(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewGroupRepository(mock.DB)

	groupID := uuid.New()
	group := &domain.Group{ID: groupID}
	logEntry := &domain.GroupLog{ID: uuid.New(), GroupID: groupID}
	mask := []string{"name"}

	mock.Mock.ExpectQuery(`(?i)SELECT count\(\*\) FROM "users" JOIN group_members.*`).
		WithArgs(groupID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectExec(`(?i)UPDATE "groups" SET.*WHERE id = \$2 AND "id" = \$3`).
		WithArgs(sqlmock.AnyArg(), groupID, groupID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.Mock.ExpectCommit()

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectQuery(`(?i)INSERT INTO "group_logs".*RETURNING`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(logEntry.ID))
	mock.Mock.ExpectCommit()

	mock.Mock.ExpectQuery(`(?i)SELECT \* FROM "groups" WHERE id = \$1 AND "groups"\."id" = \$2 ORDER BY "groups"\."id" LIMIT \$3`).
		WithArgs(groupID, groupID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(groupID))

	updatedGroup, total, err := repo.UpdateGroup(context.Background(), group, mask, logEntry)

	assert.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.NotNil(t, updatedGroup)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestRemoveMemberAndTheirSharedItems(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewGroupRepository(mock.DB)

	groupID := uuid.New()
	memberID := uuid.New()
	logEntry := &domain.GroupLog{ID: uuid.New(), GroupID: groupID}

	mock.Mock.ExpectBegin()

	mock.Mock.ExpectExec(`(?i)DELETE FROM "group_items" WHERE group_id = \$1 AND owner_id = \$2`).
		WithArgs(groupID, memberID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectExec(`(?i)DELETE FROM "group_members" WHERE group_id = \$1 AND user_id = \$2`).
		WithArgs(groupID, memberID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.Mock.ExpectQuery(`(?i)INSERT INTO "group_logs".*RETURNING`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(logEntry.ID))

	mock.Mock.ExpectCommit()

	err := repo.RemoveMemberAndTheirSharedItems(context.Background(), groupID, memberID, logEntry)

	assert.NoError(t, err)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestDeleteGroup(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewGroupRepository(mock.DB)

	groupID := uuid.New()

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectExec(`(?i)DELETE FROM "groups" WHERE id = \$1`).
		WithArgs(groupID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.Mock.ExpectCommit()

	err := repo.DeleteGroup(context.Background(), groupID)

	assert.NoError(t, err)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestCreateLog(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewGroupRepository(mock.DB)

	logEntry := &domain.GroupLog{ID: uuid.New()}

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectQuery(`(?i)INSERT INTO "group_logs".*RETURNING`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(logEntry.ID))
	mock.Mock.ExpectCommit()

	err := repo.CreateLog(context.Background(), logEntry)

	assert.NoError(t, err)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}
