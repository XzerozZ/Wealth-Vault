package utils_test

import (
	"testing"
	"wealth-vault/user-service/pkg/utils"

	"github.com/stretchr/testify/assert"
)

func TestMaskBankAccount(t *testing.T) {
	t.Run("Should mask correctly for standard 10-digit account", func(t *testing.T) {
		accNum := "1234567890"
		expected := "xxx-x-7890"

		result := utils.MaskBankAccount(accNum)

		assert.Equal(t, expected, result)
	})

	t.Run("Should return original if length is 4 or less", func(t *testing.T) {
		assert.Equal(t, "1234", utils.MaskBankAccount("1234"))
		assert.Equal(t, "", utils.MaskBankAccount(""))
	})

	t.Run("Should handle account numbers slightly longer than 4", func(t *testing.T) {
		accNum := "56789"
		expected := "xxx-x-6789"

		result := utils.MaskBankAccount(accNum)

		assert.Equal(t, expected, result)
	})

	t.Run("Should always keep the last 4 digits even with different inputs", func(t *testing.T) {
		accNum := "000011112222"
		expected := "xxx-x-2222"

		result := utils.MaskBankAccount(accNum)

		assert.Equal(t, expected, result)
	})
}
