package testutil

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MockDB struct {
	DB   *gorm.DB
	Mock sqlmock.Sqlmock
	SQL  *sql.DB
}

func NewMockDB(t *testing.T) *MockDB {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	return &MockDB{
		DB:   gormDB,
		Mock: mock,
		SQL:  sqlDB,
	}
}

func (m *MockDB) Close() {
	m.SQL.Close()
}

func (m *MockDB) ExpectDone(t *testing.T) {
	require.NoError(t, m.Mock.ExpectationsWereMet())
}
