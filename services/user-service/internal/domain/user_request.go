package domain

import "github.com/google/uuid"

type UpdateUserInput struct {
	ID          uuid.UUID
	Firstname   string
	Lastname    string
	Username    string
	Profile     string
	Phonenumber string
	BirthdayStr string
	UpdateMask  []string
}
