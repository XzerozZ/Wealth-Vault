package mapper_test

import (
	"testing"
	"time"

	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	"wealth-vault/asset-service/pkg/utils/mapper"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestToAccountDomain(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name     string
		input    *pb.CreateAccountRequest
		expected *domain.Account
	}{
		{
			name: "Success - Valid Mapping",
			input: &pb.CreateAccountRequest{
				Name:        "Savings Account",
				Amount:      1000.50,
				BankName:    "KBank",
				BankAcc:     "123-456",
				Type:        pb.BankAccType_BANK_ACC_TYPE_SAVINGS,
				Description: "Monthly savings",
			},
			expected: &domain.Account{
				UserID:      userID,
				Name:        "Savings Account",
				Amount:      1000.50,
				BankName:    "KBank",
				BankAccount: "123-456",
				Type:        domain.BankTypeSavings,
				Description: "Monthly savings",
			},
		},
		{
			name: "Fallback - Unknown Type should use default",
			input: &pb.CreateAccountRequest{
				Type: pb.BankAccType(999),
			},
			expected: &domain.Account{
				UserID: userID,
				Type:   domain.BankTypeSavings,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.ToAccountDomain(tt.input, userID)

			assert.Equal(t, tt.expected.UserID, result.UserID)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.Type, result.Type)
			assert.Equal(t, tt.expected.Amount, result.Amount)
		})
	}
}

func TestToBankProto(t *testing.T) {
	accID := uuid.New()
	userID := uuid.New()
	now := time.Now()

	tests := []struct {
		name     string
		input    *domain.Account
		expected pb.BankAccType
	}{
		{
			name: "Map SAVINGS correctly",
			input: &domain.Account{
				ID:     accID,
				UserID: userID,
				Type:   "SAVINGS",
			},
			expected: pb.BankAccType_BANK_ACC_TYPE_SAVINGS,
		},
		{
			name: "Map Case-Insensitive with spaces",
			input: &domain.Account{
				Type: " fixed_deposit ",
			},
			expected: pb.BankAccType_BANK_ACC_TYPE_FIXED_DEPOSIT,
		},
		{
			name: "Map Unknown Type to UNSPECIFIED",
			input: &domain.Account{
				Type: "BITCOIN_WALLET",
			},
			expected: pb.BankAccType_BANK_ACC_TYPE_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.CreatedAt = now
			result := mapper.ToBankProto(tt.input)

			assert.Equal(t, tt.expected, result.Type)
		})
	}
}

func TestToBankProtoSlice(t *testing.T) {
	t.Run("Empty slice should return empty slice not nil", func(t *testing.T) {
		result := mapper.ToBankProtoSlice([]*domain.Account{})
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("Should map all elements", func(t *testing.T) {
		accounts := []*domain.Account{
			{Name: "Acc 1", CreatedAt: time.Now()},
			{Name: "Acc 2", CreatedAt: time.Now()},
		}
		result := mapper.ToBankProtoSlice(accounts)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, "Acc 1", result[0].Name)
		assert.Equal(t, "Acc 2", result[1].Name)
	})
}
