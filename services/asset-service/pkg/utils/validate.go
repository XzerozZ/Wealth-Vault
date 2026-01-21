package utils

import (
	"errors"

	"github.com/google/uuid"
)

func ValidateIDs(idStr, uidStr string) (uuid.UUID, uuid.UUID, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, uuid.Nil, errors.New("invalid liability id format")
	}

	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return uuid.Nil, uuid.Nil, errors.New("invalid user id format")
	}

	return id, uid, nil
}
