package utils

import (
	"crypto/rand"
	"math/big"
)

func GenerateOTP(length int) (string, error) {
	const numbers = "0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(numbers))))
		if err != nil {
			return "", err
		}
		result[i] = numbers[num.Int64()]
	}

	return string(result), nil
}
