package mapper

import (
	"wealth-vault/api-gateway/internal/domain"
	pb "wealth-vault/api-gateway/pkg/pb/proto/user"
)

func ToUserDomain(p *pb.User) *domain.User {
	if p == nil {
		return nil
	}

	return &domain.User{
		ID:          p.Id,
		Username:    p.Username,
		Email:       p.Email,
		Firstname:   p.Firstname,
		Lastname:    p.Lastname,
		Phonenumber: p.Phonenumber,
		Profile:     p.Profile,
		Birthday:    TimeToPtr(p.Birthday.AsTime()),
		CreatedAt:   p.CreatedAt.AsTime(),
		UpdatedAt:   p.UpdatedAt.AsTime(),
	}
}
