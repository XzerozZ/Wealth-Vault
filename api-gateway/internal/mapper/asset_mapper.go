package mapper

import (
	"strings"
)

func SafeMapEnum(valueMap map[string]int32, key string, prefix string) int32 {
	key = strings.ToUpper(strings.TrimSpace(key))
	if v, ok := valueMap[key]; ok {
		return v
	}

	if v, ok := valueMap[prefix+key]; ok {
		return v
	}

	return 0
}
