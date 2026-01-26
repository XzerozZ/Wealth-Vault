package mapper

import (
	"strings"
	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var domainToProtoAccType = map[string]pb.BankAccType{
	"SAVINGS":       pb.BankAccType_BANK_ACC_TYPE_SAVINGS,
	"CURRENT":       pb.BankAccType_BANK_ACC_TYPE_CURRENT,
	"FIXED_DEPOSIT": pb.BankAccType_BANK_ACC_TYPE_FIXED_DEPOSIT,
}

func ToBankProto(d *domain.Account) *pb.Account {
	normalizedType := strings.ToUpper(strings.TrimSpace(string(d.Type)))
	accTypeEnum, ok := domainToProtoAccType[normalizedType]
	if !ok {
		accTypeEnum = pb.BankAccType_BANK_ACC_TYPE_UNSPECIFIED
	}

	res := &pb.Account{
		Id:          d.ID.String(),
		UserId:      d.UserID.String(),
		Name:        d.Name,
		BankName:    d.BankName,
		BankAcc:     d.BankAccount,
		Type:        accTypeEnum,
		Amount:      d.Amount,
		Description: d.Description,
		Files:       ToPbFiles(d.Files),
		CreatedAt:   timestamppb.New(d.CreatedAt),
		UpdatedAt:   timestamppb.New(d.UpdatedAt),
	}

	return res
}
