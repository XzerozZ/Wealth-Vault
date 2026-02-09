package utils

import (
	"wealth-vault/user-service/internal/domain"

	pb "wealth-vault/user-service/pkg/pb/proto/user"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToUserProto(d *domain.User) *pb.User {
	if d == nil {
		return nil
	}

	var birthdayPb *timestamppb.Timestamp
	if d.Birthday != nil {
		birthdayPb = timestamppb.New(*d.Birthday)
	}

	res := &pb.User{
		Id:            d.ID.String(),
		Email:         d.Email,
		Firstname:     d.Firstname,
		Lastname:      d.Lastname,
		Username:      d.Username,
		Profile:       d.Profile,
		Phonenumber:   d.Phonenumber,
		Birthday:      birthdayPb,
		Sharedage:     int32(d.AutoShareAge),
		Sharedenabled: d.IsAutoShareEnabled,
		CreatedAt:     timestamppb.New(d.CreatedAt),
		UpdatedAt:     timestamppb.New(d.UpdatedAt),
	}

	return res
}

func ToGroupProto(g *domain.Group) *pb.Group {
	if g == nil {
		return nil
	}

	return &pb.Group{
		Id:          g.ID.String(),
		Name:        g.GroupName,
		Profile:     g.GroupProfile,
		UserId:      g.CreatedBy.String(),
		CreatedAt:   timestamppb.New(g.CreatedAt),
		UpdatedAt:   timestamppb.New(g.UpdatedAt),
		MemberCount: 0,
	}
}
