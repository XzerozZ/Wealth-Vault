package utils_test

import (
	"testing"
	"wealth-vault/notification-service/pkg/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestParseUUIDPtr(t *testing.T) {
	validID := uuid.New()

	tests := []struct {
		name     string
		input    string
		expected *uuid.UUID
	}{
		{
			name:     "Success - Valid UUID String",
			input:    validID.String(),
			expected: &validID,
		},
		{
			name:     "Edge Case - Empty String",
			input:    "",
			expected: nil,
		},
		{
			name:     "Edge Case - Invalid UUID Format",
			input:    "not-a-uuid-123",
			expected: nil,
		},
		{
			name:     "Edge Case - Short String",
			input:    "abc-123",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ParseUUIDPtr(tt.input)

			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, *tt.expected, *result)
			}
		})
	}
}
