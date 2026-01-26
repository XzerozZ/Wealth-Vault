package helper

import (
	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
)

var ProtoToDomainAccType = map[pb.BankAccType]domain.BankType{
	pb.BankAccType_BANK_ACC_TYPE_SAVINGS:       domain.BankTypeSavings,
	pb.BankAccType_BANK_ACC_TYPE_CURRENT:       domain.BankTypeCurrent,
	pb.BankAccType_BANK_ACC_TYPE_FIXED_DEPOSIT: domain.BankTypeFixedDeposit,
}

func ApplyUpdateAccFields(req *pb.UpdateAccountRequest, acc *domain.Account) error {
	paths := req.UpdateMask.GetPaths()
	if len(paths) == 0 {
		return nil
	}

	for _, path := range paths {
		switch path {
		case "name":
			if req.Acc.Name != "" {
				acc.Name = req.Acc.Name
			}

		case "bank_name":
			if req.Acc.BankName != "" {
				acc.BankName = req.Acc.BankName
			}

		case "bank_account":
			if req.Acc.BankAcc != "" {
				acc.BankAccount = req.Acc.BankAcc
			}

		case "amount":
			if req.Acc.Amount != 0 {
				acc.Amount = req.Acc.Amount
			}

		case "description":
			if req.Acc.Description != "" {
				acc.Description = req.Acc.Description
			}

		case "type":
			if req.Acc.Type != pb.BankAccType_BANK_ACC_TYPE_UNSPECIFIED {
				if val, ok := ProtoToDomainAccType[req.Acc.Type]; ok {
					acc.Type = val
				}
			}
		}
	}

	return nil
}
