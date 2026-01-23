package domain

import "time"

type User struct {
	ID          string     `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	Firstname   string     `json:"first_name"`
	Lastname    string     `json:"last_name"`
	Phonenumber string     `json:"phone_number"`
	Profile     string     `json:"profile"`
	Birthday    *time.Time `json:"birthday"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type UpdateRequest struct {
	Firstname   string `json:"first_name"    form:"first_name"    mask:"first_name"`
	Lastname    string `json:"last_name"     form:"last_name"     mask:"last_name"`
	Username    string `json:"username"     form:"username"     mask:"username"`
	Profile     string `json:"profile"      form:"profile"      mask:"profile"`
	Phonenumber string `json:"phonenumber"  form:"phonenumber"  mask:"phonenumber"`
	Birthday    string `json:"birthday"     form:"birthday"     mask:"birthday"`
}
