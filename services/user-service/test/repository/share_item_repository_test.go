package repository_test

import (
	"context"
	"regexp"
	"testing"
	"wealth-vault/user-service/internal/domain"
	"wealth-vault/user-service/internal/repository"
	testutil "wealth-vault/user-service/test/testdb"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestShareItemtoGroup(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewShareItemRepository(mock.DB)

	id := uuid.New()

	items := []domain.GroupItem{
		{ID: id},
	}

	mock.Mock.ExpectBegin()

	mock.Mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "group_items"`)).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(id),
		)

	mock.Mock.ExpectCommit()

	err := repo.ShareItemtoGroup(context.Background(), items)

	assert.NoError(t, err)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestShareItemtoGroup_Empty(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewShareItemRepository(mock.DB)

	err := repo.ShareItemtoGroup(context.Background(), []domain.GroupItem{})
	assert.NoError(t, err)
}

func TestShareItemtoFriend(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()

	repo := repository.NewShareItemRepository(mock.DB)

	id := uuid.New()

	items := []domain.FriendItem{
		{ID: id},
	}

	mock.Mock.ExpectBegin()

	mock.Mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "friend_items"`)).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(id),
		)

	mock.Mock.ExpectCommit()

	err := repo.ShareItemtoFriend(context.Background(), items)

	assert.NoError(t, err)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestShareItemtoEmail(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewShareItemRepository(mock.DB)

	id := uuid.New()

	items := []domain.EmailItem{
		{ID: id},
	}

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "email_items"`)).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(id),
		)
	mock.Mock.ExpectCommit()

	err := repo.ShareItemtoEmail(context.Background(), items)
	assert.NoError(t, err)
}

func TestIsItemSharedtoGroup(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewShareItemRepository(mock.DB)

	groupID, entityID := uuid.New(), uuid.New()

	mock.Mock.ExpectQuery(`SELECT count\(\*\) FROM "group_items"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	ok, err := repo.IsItemSharedtoGroup(context.Background(), groupID, entityID, "POST")

	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestIsItemSharedtoFriend(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewShareItemRepository(mock.DB)

	mock.Mock.ExpectQuery(`SELECT count\(\*\) FROM "friend_items"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	ok, err := repo.IsItemSharedtoFriend(context.Background(), uuid.New(), uuid.New(), "POST")

	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestIsItemSharedtoEmail(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewShareItemRepository(mock.DB)

	mock.Mock.ExpectQuery(`SELECT count\(\*\) FROM "email_items"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	ok, err := repo.IsItemSharedtoEmail(context.Background(), uuid.New(), "a@b.com", "POST")

	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestGetExistingSharedMap(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewShareItemRepository(mock.DB)

	ownerID, friendID := uuid.New(), uuid.New()
	entityID := uuid.New()

	mock.Mock.ExpectQuery(`SELECT entity_id, entity_type FROM "friend_items"`).
		WillReturnRows(sqlmock.NewRows([]string{"entity_id", "entity_type"}).
			AddRow(entityID, "POST"))

	res, err := repo.GetExistingSharedMap(context.Background(), ownerID, friendID)

	assert.NoError(t, err)
	assert.True(t, res["POST:"+entityID.String()])
}

func TestGetPendingEmails(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewShareItemRepository(mock.DB)

	mock.Mock.ExpectQuery(`SELECT \* FROM "email_items"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

	res, err := repo.GetPendingEmails(context.Background())

	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestDeleteIteminGroup(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewShareItemRepository(mock.DB)

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectExec(`DELETE FROM "group_items"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.Mock.ExpectCommit()

	err := repo.DeleteIteminGroup(context.Background(), uuid.New(), uuid.New())
	assert.NoError(t, err)
}

func TestDeleteAllReferencesByEntityID(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewShareItemRepository(mock.DB)

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectExec(`DELETE FROM "group_items"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.Mock.ExpectExec(`DELETE FROM "friend_items"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.Mock.ExpectExec(`DELETE FROM "email_items"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.Mock.ExpectCommit()

	err := repo.DeleteAllReferencesByEntityID(context.Background(), uuid.New())
	assert.NoError(t, err)
}

func TestAddMember(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewShareItemRepository(mock.DB)

	member := domain.GroupMember{GroupID: uuid.New(), UserID: uuid.New()}
	members := []domain.GroupMember{member}

	mock.Mock.ExpectBegin()

	mock.Mock.ExpectQuery(`INSERT INTO "group_members".*ON CONFLICT.*RETURNING "group_id"`).
		WithArgs(member.UserID, sqlmock.AnyArg(), sqlmock.AnyArg(), member.GroupID).
		WillReturnRows(sqlmock.NewRows([]string{"group_id"}).AddRow(member.GroupID))

	mock.Mock.ExpectCommit()

	err := repo.AddMember(context.Background(), members)
	assert.NoError(t, err)
	assert.NoError(t, mock.Mock.ExpectationsWereMet())
}

func TestBatchCreateViewers(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewShareItemRepository(mock.DB)

	viewers := []domain.GroupItemViewer{
		{GroupItemID: uuid.New(), ViewerID: uuid.New()},
	}

	mock.Mock.ExpectBegin()
	mock.Mock.ExpectExec(`(?i)INSERT INTO "group_item_viewers".*ON CONFLICT`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.Mock.ExpectCommit()

	err := repo.BatchCreateViewers(context.Background(), viewers)
	assert.NoError(t, err)
}

func TestCountItemsByOwner(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewShareItemRepository(mock.DB)

	mock.Mock.ExpectQuery(`SELECT count\(\*\) FROM "group_items"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	count, err := repo.CountItemsByOwner(context.Background(), []string{"id1"}, uuid.New())

	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

func TestGetOwnedItemIDs(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewShareItemRepository(mock.DB)

	id := uuid.New()

	mock.Mock.ExpectQuery(`SELECT "id" FROM "group_items"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))

	res, err := repo.GetOwnedItemIDs(context.Background(), []string{id.String()}, uuid.New())

	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestGetItemSharedTargets(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewShareItemRepository(mock.DB)

	mock.Mock.MatchExpectationsInOrder(false)

	mock.Mock.ExpectQuery(`FROM "group_items"`).
		WillReturnRows(sqlmock.NewRows([]string{"group_id"}).AddRow(uuid.New()))

	mock.Mock.ExpectQuery(`FROM "friend_items"`).
		WillReturnRows(sqlmock.NewRows([]string{"friend_id"}).AddRow(uuid.New()))

	mock.Mock.ExpectQuery(`FROM "email_items"`).
		WillReturnRows(sqlmock.NewRows([]string{"email"}).AddRow("a@b.com"))

	res, err := repo.GetItemSharedTargets(context.Background(), uuid.New(), uuid.New(), "POST")

	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestGetItemsSharedByFriend(t *testing.T) {
	mock := testutil.NewMockDB(t)
	defer mock.Close()
	repo := repository.NewShareItemRepository(mock.DB)

	mock.Mock.ExpectQuery(`(?i)SELECT entity_id, entity_type FROM friend_items .* UNION .*`).
		WillReturnRows(sqlmock.NewRows([]string{"entity_id", "entity_type"}).
			AddRow(uuid.New(), "POST"))

	res, err := repo.GetItemsSharedByFriend(context.Background(), uuid.New(), uuid.New())

	assert.NoError(t, err)
	assert.Len(t, res, 1)
}
