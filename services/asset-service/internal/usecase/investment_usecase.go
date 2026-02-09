package usecase

import (
	"context"
	"errors"
	"time"
	"wealth-vault/asset-service/internal/domain"
	repo "wealth-vault/asset-service/internal/repository/interface"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	userPb "wealth-vault/asset-service/pkg/pb/proto/user"
	"wealth-vault/asset-service/pkg/utils"
	helper "wealth-vault/asset-service/pkg/utils/helper"
	"wealth-vault/asset-service/pkg/utils/mapper"

	"github.com/google/uuid"
)

type InvestmentUsecase struct {
	inRepo     repo.InvestmentRepository
	fileRepo   repo.FileRepository
	storage    *utils.StorageClient
	userClient userPb.UserServiceClient
}

func NewInvestmentUsecase(r repo.InvestmentRepository, fr repo.FileRepository, s *utils.StorageClient, userClient userPb.UserServiceClient) InvestmentUsecase {
	return InvestmentUsecase{
		inRepo:     r,
		fileRepo:   fr,
		storage:    s,
		userClient: userClient,
	}
}

func (u *InvestmentUsecase) CreateInvestment(ctx context.Context, req *pb.CreateInvestmentRequest) (*pb.InvestmentResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	inType := domain.InvestTypeStockUS
	if val, ok := helper.ProtoToDomainInType[req.Type]; ok {
		inType = val
	}

	var domainFiles []domain.FileAssociate
	if len(req.NewFiles) > 0 {
		for _, f := range req.NewFiles {
			domainFiles = append(domainFiles, domain.FileAssociate{
				Link:     f.Url,
				FileType: f.FileType,
				UserID:   userID,
			})
		}
	}

	in := &domain.Investment{
		UserID:       userID,
		Name:         req.Name,
		Symbol:       req.Symbol,
		Type:         inType,
		BrokerName:   req.BrokerName,
		Quantity:     req.Quantity,
		CostPerPrice: req.CostPrice,
		Description:  req.Description,
		Files:        domainFiles,
	}

	if err := u.inRepo.CreateInvestment(ctx, in); err != nil {
		return nil, err
	}

	return &pb.InvestmentResponse{
		Invest: mapper.ToInvestProto(in),
	}, nil
}

func (u *InvestmentUsecase) GetInvestment(ctx context.Context, req *pb.GetAssetRequest) (*pb.InvestmentArrayResponse, error) {
	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	accounts, err := u.inRepo.GetInvestment(ctx, uid)
	if err != nil {
		return nil, err
	}

	var InvestList []*pb.Investment
	for _, item := range accounts {
		InvestList = append(InvestList, mapper.ToInvestProto(item))
	}

	return &pb.InvestmentArrayResponse{
		Success: true,
		Invest:  InvestList,
	}, nil
}

func (u *InvestmentUsecase) GetInvestmentByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.InvestmentArrayResponse, error) {
	var ids []uuid.UUID
	for _, idStr := range req.Ids {
		if parsedID, err := uuid.Parse(idStr); err == nil {
			ids = append(ids, parsedID)
		}
	}

	bu, err := u.inRepo.GetInvestmentByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	var pbIns []*pb.Investment
	for _, a := range bu {
		pbIns = append(pbIns, mapper.ToInvestProto(a))
	}

	return &pb.InvestmentArrayResponse{
		Invest: pbIns,
	}, nil
}

func (u *InvestmentUsecase) GetBatchInvestmentByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.InvestmentArrayResponse, error) {
	var ids []uuid.UUID
	for _, idStr := range req.Ids {
		if parsedID, err := uuid.Parse(idStr); err == nil {
			ids = append(ids, parsedID)
		}
	}

	bu, err := u.inRepo.GetBatchInvestmentByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	var pbIns []*pb.Investment
	for _, a := range bu {
		pbIns = append(pbIns, mapper.ToInvestProto(a))
	}

	return &pb.InvestmentArrayResponse{
		Invest: pbIns,
	}, nil
}

func (u *InvestmentUsecase) GetInvestmentByID(ctx context.Context, req *pb.GetAssetByIDRequest) (*pb.InvestmentResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	in, err := u.inRepo.GetInvestmentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.InvestmentResponse{
		Success: true,
		Invest:  mapper.ToInvestProto(in),
	}, nil
}

func (u *InvestmentUsecase) UpdateInvestment(ctx context.Context, req *pb.UpdateInvestmentRequest) (*pb.InvestmentResponse, error) {
	id, uid, err := utils.ValidateIDs(req.Id, req.Invest.UserId)
	if err != nil {
		return nil, err
	}

	in, err := u.inRepo.GetInvestmentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := helper.ApplyUpdateInFields(req, in); err != nil {
		return nil, err
	}

	syncParams := domain.FileSyncParams{
		UserID:        uid,
		EntityID:      id,
		EntityType:    "investment",
		NewFiles:      req.NewFiles,
		DeleteFileIDs: req.DeleteFileIds,
	}

	err = helper.SyncEntityFiles(ctx, u.fileRepo, u.storage, syncParams)
	if err != nil {
		return nil, err
	}

	updatedInvest, err := u.inRepo.UpdateInvestment(ctx, in)
	if err != nil {
		return nil, err
	}

	return &pb.InvestmentResponse{
		Success: true,
		Invest:  mapper.ToInvestProto(updatedInvest),
	}, nil
}

func (u *InvestmentUsecase) DeleteInvestment(ctx context.Context, req *pb.DeleteAssetRequest) (*pb.DeleteAssetResponse, error) {
	id, uid, err := utils.ValidateIDs(req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	_, err = u.inRepo.GetInvestmentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := u.inRepo.SoftDeleteInvestment(ctx, id, uid); err != nil {
		return nil, err
	}

	return &pb.DeleteAssetResponse{
		Success: true,
	}, nil
}

func (u *InvestmentUsecase) CleanupExpiredInvestment(ctx context.Context) error {
	cutoffTime := time.Now().AddDate(0, 0, -7)
	GetExpiredInvest, err := u.inRepo.GetExpiredInvestment(ctx, cutoffTime)
	if err != nil {
		return err
	}

	if len(GetExpiredInvest) == 0 {
		return err
	}

	for _, i := range GetExpiredInvest {
		helper.CleanupAssetResource(
			ctx,
			i.ID,
			i.Files,
			u.storage,
			u.userClient,
			func(id uuid.UUID) error {
				return u.inRepo.HardDeleteInvestment(ctx, id)
			},
		)
	}

	return nil
}
