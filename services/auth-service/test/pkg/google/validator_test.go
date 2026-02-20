package google_test

import (
	"context"
	"testing"

	"wealth-vault/auth-service/pkg/google"

	"github.com/stretchr/testify/assert"
)

func TestGoogleValidator_Validate(t *testing.T) {
	clientID := "test-client-id.apps.googleusercontent.com"

	validator := google.NewGoogleValidator(clientID)

	t.Run("error - invalid token format", func(t *testing.T) {
		fakeToken := "this_is_not_a_valid_jwt_token"

		payload, err := validator.Validate(context.Background(), fakeToken)

		assert.Nil(t, payload)
		assert.Error(t, err)

		assert.Contains(t, err.Error(), "idtoken")
	})

	t.Run("error - invalid signature", func(t *testing.T) {
		fakeJWT := "eyJhbGciOiJSUzI1NiIsImtpZCI6IjEyMyIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.invalid_signature"

		payload, err := validator.Validate(context.Background(), fakeJWT)

		assert.Nil(t, payload)
		assert.Error(t, err)
	})
}
