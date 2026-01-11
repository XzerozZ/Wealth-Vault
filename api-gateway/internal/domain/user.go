package domain

type UpdateRequest struct {
	Firstname   string `json:"firstname" form:"first_name"`
	Lastname    string `json:"lastname" form:"last_name"`
	Username    string `json:"username" form:"username"`
	Profile     string `json:"profile" form:"profile"`
	Phonenumber string `json:"phone_number" form:"phone_number"`
	Birthday    string `json:"birthday" form:"birthday"` // Format: YYYY-MM-DD
}
