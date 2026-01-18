package utils

import (
	"strconv"
	"strings"
	pb "wealth-vault/api-gateway/pkg/pb/proto/asset"
)

func Parseamount(amountStr string) (float64, error) {
	val, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return 0, err
	}

	return val, nil
}

func GetFolderName(t pb.AssetType) string {
	rawName := t.String()
	if t == pb.AssetType_ASSET_TYPE_UNSPECIFIED {
		return "misc"
	}

	cleanName := strings.TrimPrefix(rawName, "ASSET_TYPE_")
	return strings.ToLower(cleanName)
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
