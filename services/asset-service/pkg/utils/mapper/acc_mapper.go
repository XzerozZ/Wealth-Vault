package mapper

import (
	"strings"
	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	"wealth-vault/asset-service/pkg/utils/helper"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var domainToProtoAccType = map[string]pb.BankAccType{
	"SAVINGS":       pb.BankAccType_BANK_ACC_TYPE_SAVINGS,
	"CURRENT":       pb.BankAccType_BANK_ACC_TYPE_CURRENT,
	"FIXED_DEPOSIT": pb.BankAccType_BANK_ACC_TYPE_FIXED_DEPOSIT,
}

func ToAccountDomain(req *pb.CreateAccountRequest, userID uuid.UUID) *domain.Account {
	accType := domain.BankTypeSavings
	if val, ok := helper.ProtoToDomainAccType[req.Type]; ok {
		accType = val
	}

	return &domain.Account{
		UserID:      userID,
		Name:        req.Name,
		Amount:      req.Amount,
		BankName:    req.BankName,
		BankAccount: req.BankAcc,
		Type:        accType,
		Description: req.Description,
		Files:       ToDomainFiles(req.NewFiles, userID),
	}
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
		DeletedAt:   timestamppb.New(d.DeletedAt.Time),
	}

	return res
}

func ToBankProtoSlice(accounts []*domain.Account) []*pb.Account {
	if len(accounts) == 0 {
		return []*pb.Account{}
	}

	res := make([]*pb.Account, len(accounts))
	for i, acc := range accounts {
		res[i] = ToBankProto(acc)
	}
	return res
}
