package utils

import (
	"wealth-vault/user-service/internal/domain"

	pb "wealth-vault/user-service/pkg/pb/proto/user"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToUserProto(d *domain.User) *pb.User {
	res := &pb.User{
		Id:          d.ID.String(),
		Email:       d.Email,
		Firstname:   d.Firstname,
		Lastname:    d.Lastname,
		Username:    d.Username,
		Profile:     d.Profile,
		Phonenumber: d.Phonenumber,
		Birthday:    timestamppb.New(*d.Birthday),
		CreatedAt:   timestamppb.New(d.CreatedAt),
		UpdatedAt:   timestamppb.New(d.UpdatedAt),
	}

	return res
}
