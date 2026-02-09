package utils

import "github.com/google/uuid"

func ParseUUIDPtr(idStr string) *uuid.UUID {
	if idStr == "" {
		return nil
	}
	if id, err := uuid.Parse(idStr); err == nil {
		return &id
	}
	return nil
}
