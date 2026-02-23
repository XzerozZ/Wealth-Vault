package utils

import (
	"errors"

	"github.com/google/uuid"
)

func ParseUUID(id string) (uuid.UUID, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, errors.New("invalid uuid")
	}
	return uid, nil
}
