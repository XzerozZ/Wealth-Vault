package domain

import "time"

type Group struct {
	ID           string    `json:"id"`
	GroupName    string    `json:"group_name"`
	GroupProfile string    `json:"group_profile"`
	CreatedBy    string    `json:"created_by"`
	MemberCount  int64     `json:"member_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateGroupRequest struct {
	GroupName    string `json:"group_name" form:"group_name"`
	GroupProfile string
	MemberIDS    []string `json:"member_ids" form:"member_ids"`
}

type UpdateGroupRequest struct {
	GroupName    string `json:"group_name" form:"group_name" mask:"group_name"`
	GroupProfile string `json:"group_profile" form:"group_profile" mask:"group_profile"`
}
