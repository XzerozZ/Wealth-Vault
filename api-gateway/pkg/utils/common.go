package utils

import (
	"strconv"
	"strings"
	"time"
	pb "wealth-vault/api-gateway/pkg/pb/proto/asset"

	"google.golang.org/protobuf/types/known/timestamppb"
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

	cleanName := strings.TrimPrefix(rawName, "LIABILITY_TYPE_")
	return strings.ToLower(cleanName)
}

func GetFolderLiaName(t pb.LiabilityType) string {
	rawName := t.String()
	if t == pb.LiabilityType_LIABILITY_TYPE_UNSPECIFIED {
		return "misc"
	}

	cleanName := strings.TrimPrefix(rawName, "ASSET_TYPE_")
	return strings.ToLower(cleanName)
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
