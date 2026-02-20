package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Generator interface {
	CreateToken(userID string, email string, tokenType string, duration time.Duration) (string, error)
	VerifyToken(tokenString string) (jwt.MapClaims, error)
}
