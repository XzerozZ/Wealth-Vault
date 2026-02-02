package utils

import (
	"strconv"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func Parseamount(amountStr string) (float64, error) {
	val, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return 0, err
	}

	return val, nil
}

func ToProtoTime(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}

	return timestamppb.New(*t)
}

func Unique(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}
	return u
}
