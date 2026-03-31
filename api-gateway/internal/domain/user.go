package domain

import "time"

type User struct {
	ID            string     `json:"id"`
	Username      string     `json:"username"`
	Email         string     `json:"email"`
	Firstname     string     `json:"first_name"`
	Lastname      string     `json:"last_name"`
	Phonenumber   string     `json:"phone_number"`
	Profile       string     `json:"profile"`
	Birthday      *time.Time `json:"birthday"`
	SharedAge     *int32     `json:"shared_age"`
	SharedEnabled *bool      `json:"shared_enabled"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	IsClose       bool       `json:"is_close"`
}

type UpdateRequest struct {
	Firstname     string `json:"first_name"    form:"first_name"    mask:"first_name"`
	Lastname      string `json:"last_name"     form:"last_name"     mask:"last_name"`
	Username      string `json:"username"     form:"username"     mask:"username"`
	Profile       string `json:"profile"      form:"profile"      mask:"profile"`
	Phonenumber   string `json:"phonenumber"  form:"phonenumber"  mask:"phonenumber"`
	Birthday      string `json:"birthday"     form:"birthday"     mask:"birthday"`
	SharedAge     string `json:"shared_age" form:"shared_age"     mask:"shared_age"`
	SharedEnabled string `json:"shared_enabled" form:"shared_enabled"     mask:"share_enabled"`
}

type AddFriendRequest struct {
	RequesterID string `json:"requester_id" form:"requester_id"`
}

type AcceptFriendRequest struct {
	RequesterID string `json:"requester_id" form:"requester_id"`
	Action      string `json:"action" form:"action"`
}

type SetCloseFriendRequest struct {
	FriendID string `json:"friend_id" form:"friend_id" validate:"required"`
	IsClose  string `json:"is_close" form:"is_close"`
}
