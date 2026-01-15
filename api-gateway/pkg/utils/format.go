package utils

import "strconv"

func Parseamount(amountStr string) (float64, error) {
	val, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return 0, err
	}

	return val, nil
}
