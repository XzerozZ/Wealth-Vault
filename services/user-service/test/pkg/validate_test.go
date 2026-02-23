package utils_test

import (
	"testing"
	"wealth-vault/user-service/pkg/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestParseUUID(t *testing.T) {
	t.Run("should return parsed uuid when input is valid", func(t *testing.T) {
		input := "550e8400-e29b-41d4-a716-446655440000"
		expected, _ := uuid.Parse(input)

		result, err := utils.ParseUUID(input)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		assert.NotEqual(t, uuid.Nil, result)
	})

	t.Run("should return error when input is invalid string", func(t *testing.T) {
		input := "invalid-uuid-format"

		result, err := utils.ParseUUID(input)

		assert.Error(t, err)
		assert.Equal(t, "invalid uuid", err.Error())
		assert.Equal(t, uuid.Nil, result)
	})

	t.Run("should return error when input is empty string", func(t *testing.T) {
		input := ""

		result, err := utils.ParseUUID(input)

		assert.Error(t, err)
		assert.Equal(t, "invalid uuid", err.Error())
		assert.Equal(t, uuid.Nil, result)
	})

	t.Run("should handle uuid without hyphens if google/uuid supports it", func(t *testing.T) {
		input := "550e8400e29b41d4a716446655440000"

		result, err := utils.ParseUUID(input)

		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, result)
	})
}
