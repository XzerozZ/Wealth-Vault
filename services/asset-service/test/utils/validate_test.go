package utils_test

import (
	"testing"
	"time"

	"wealth-vault/asset-service/pkg/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestValidateIDs(t *testing.T) {
	validID := uuid.New().String()
	validUID := uuid.New().String()
	invalidID := "not-a-uuid"

	tests := []struct {
		name    string
		idStr   string
		uidStr  string
		wantErr bool
		errMs   string
	}{
		{
			name:    "Success - Both Valid",
			idStr:   validID,
			uidStr:  validUID,
			wantErr: false,
		},
		{
			name:    "Error - Invalid Asset ID",
			idStr:   invalidID,
			uidStr:  validUID,
			wantErr: true,
			errMs:   "invalid liability id format",
		},
		{
			name:    "Error - Invalid User ID",
			idStr:   validID,
			uidStr:  invalidID,
			wantErr: true,
			errMs:   "invalid user id format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, uid, err := utils.ValidateIDs(tt.idStr, tt.uidStr)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMs, err.Error())
				assert.Equal(t, uuid.Nil, id)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, id)
				assert.NotEqual(t, uuid.Nil, uid)
			}
		})
	}
}

func TestParseUUIDs(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()

	t.Run("Should parse valid and skip invalid", func(t *testing.T) {
		input := []string{id1.String(), "invalid", id2.String()}
		result := utils.ParseUUIDs(input)

		assert.Len(t, result, 2)
		assert.Contains(t, result, id1)
		assert.Contains(t, result, id2)
	})

	t.Run("Empty input should return nil slice", func(t *testing.T) {
		result := utils.ParseUUIDs([]string{})
		assert.Nil(t, result)
	})
}

func TestToTimePtr(t *testing.T) {
	now := time.Now()
	ts := timestamppb.New(now)

	t.Run("Success - Convert Timestamp", func(t *testing.T) {
		result := utils.ToTimePtr(ts)
		assert.NotNil(t, result)
		assert.Equal(t, now.Unix(), result.Unix())
	})

	t.Run("Edge Case - Nil Input", func(t *testing.T) {
		result := utils.ToTimePtr(nil)
		assert.Nil(t, result)
	})
}
