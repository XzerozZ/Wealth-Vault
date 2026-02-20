package mock

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
)

type MockTokenMaker struct {
	mock.Mock
}

func (m *MockTokenMaker) CreateToken(userID string, email string, tokenType string, duration time.Duration) (string, error) {
	args := m.Called(userID, email, tokenType, duration)
	return args.String(0), args.Error(1)
}

func (m *MockTokenMaker) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	args := m.Called(tokenString)
	return args.Get(0).(jwt.MapClaims), args.Error(1)
}
