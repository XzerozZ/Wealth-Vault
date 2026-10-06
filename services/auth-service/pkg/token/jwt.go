package token

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTGenerate struct {
	secretKey string
}

func NewJWTGenerate(secret string) *JWTGenerate {
	return &JWTGenerate{
		secretKey: secret,
	}
}

func (j *JWTGenerate) CreateToken(userID string, email string, tokenType string, duration time.Duration) (string, error) {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"type":    tokenType,
		"exp":     time.Now().Add(duration).Unix(),
		"jti":     hex.EncodeToString(nonce),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTGenerate) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("token is invalid")
		}
		return []byte(j.secretKey), nil
	}

	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token has expired")
		}
		return nil, errors.New("token is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("token is invalid")
	}

	return claims, nil
}
