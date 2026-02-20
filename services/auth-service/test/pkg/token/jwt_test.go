package token_test // ✅ เป็นคนนอก 100%

import (
	"testing"
	"time"

	"wealth-vault/auth-service/pkg/token"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

const (
	secretKey = "super_secret_key_123"
	userID    = "user-1234"
	email     = "test@example.com"
	tokenType = "access"
)

func TestCreateToken(t *testing.T) {
	jwtGen := token.NewJWTGenerate(secretKey)

	t.Run("success", func(t *testing.T) {
		duration := 15 * time.Minute

		tokenStr, err := jwtGen.CreateToken(userID, email, tokenType, duration)

		assert.NoError(t, err)
		assert.NotEmpty(t, tokenStr)

		assert.Contains(t, tokenStr, ".")
	})
}

func TestVerifyToken(t *testing.T) {
	jwtGen := token.NewJWTGenerate(secretKey)

	t.Run("success", func(t *testing.T) {
		tokenStr, err := jwtGen.CreateToken(userID, email, tokenType, 15*time.Minute)
		assert.NoError(t, err)

		claims, err := jwtGen.VerifyToken(tokenStr)

		assert.NoError(t, err)
		assert.NotNil(t, claims)

		assert.Equal(t, userID, claims["user_id"])
		assert.Equal(t, email, claims["email"])
		assert.Equal(t, tokenType, claims["type"])
		assert.IsType(t, float64(0), claims["exp"])
	})

	t.Run("error - token has expired", func(t *testing.T) {
		expiredTokenStr, err := jwtGen.CreateToken(userID, email, tokenType, -1*time.Minute)
		assert.NoError(t, err)

		claims, err := jwtGen.VerifyToken(expiredTokenStr)

		assert.Nil(t, claims)
		assert.Error(t, err)
		assert.Equal(t, "token has expired", err.Error())
	})

	t.Run("error - invalid token format", func(t *testing.T) {
		claims, err := jwtGen.VerifyToken("this.is.not.a.valid.jwt")

		assert.Nil(t, claims)
		assert.Error(t, err)
		assert.Equal(t, "token is invalid", err.Error())
	})

	t.Run("error - invalid signature (wrong secret)", func(t *testing.T) {
		hackerJwtGen := token.NewJWTGenerate("hacker_secret_key_999")
		hackerToken, _ := hackerJwtGen.CreateToken(userID, email, tokenType, 15*time.Minute)

		claims, err := jwtGen.VerifyToken(hackerToken)

		assert.Nil(t, claims)
		assert.Error(t, err)
		assert.Equal(t, "token is invalid", err.Error())
	})

	t.Run("error - alg hacking attempt (none algorithm)", func(t *testing.T) {
		hackerClaims := jwt.MapClaims{
			"user_id": userID,
			"exp":     time.Now().Add(15 * time.Minute).Unix(),
		}

		hackerToken := jwt.NewWithClaims(jwt.SigningMethodNone, hackerClaims)
		hackerTokenStr, _ := hackerToken.SignedString(jwt.UnsafeAllowNoneSignatureType)

		resClaims, err := jwtGen.VerifyToken(hackerTokenStr)

		assert.Nil(t, resClaims)
		assert.Error(t, err)
		assert.Equal(t, "token is invalid", err.Error())
	})
}
