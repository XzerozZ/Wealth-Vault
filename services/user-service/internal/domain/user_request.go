package domain

type UpdateUserInput struct {
	ID          string
	Firstname   string
	Lastname    string
	Username    string
	Profile     string
	Phonenumber string
	BirthdayStr string
	UpdateMask  []string
}
