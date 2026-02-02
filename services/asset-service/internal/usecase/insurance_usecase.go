package usecase

import (
	"context"
	"errors"
	"time"
	"wealth-vault/asset-service/internal/domain"
	repo "wealth-vault/asset-service/internal/repository/interface"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	"wealth-vault/asset-service/pkg/utils"
	helper "wealth-vault/asset-service/pkg/utils/helper"
	"wealth-vault/asset-service/pkg/utils/mapper"

	"github.com/google/uuid"
)

type InsuranceUsecase struct {
	insRepo  repo.InsuranceRepository
	fileRepo repo.FileRepository
	storage  *utils.StorageClient
}

func NewInsuranceUsecase(r repo.InsuranceRepository, fr repo.FileRepository, s *utils.StorageClient) InsuranceUsecase {
	return InsuranceUsecase{
		insRepo:  r,
		fileRepo: fr,
		storage:  s,
	}
}

func (u *InsuranceUsecase) CreateInsurance(ctx context.Context, req *pb.CreateInsuranceRequest) (*pb.InsuranceResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	inType := domain.InsuranceTypeLife
	if val, ok := helper.ProtoToDomainInsType[req.Type]; ok {
		inType = val
	}

	var conDatePtr, expDatePtr *time.Time
	if req.ConDate != nil {
		t := req.ConDate.AsTime()
		conDatePtr = &t
	}

	if req.ExpDate != nil {
		t := req.ExpDate.AsTime()
		expDatePtr = &t
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

	ins := &domain.Insurance{
		UserID:         userID,
		Type:           inType,
		Name:           req.Name,
		PolicyNumber:   req.PolNum,
		CompanyName:    req.CompanyName,
		CoveragePeriod: req.CoveragePeriod,
		CoverageAmount: req.CoverageAmount,
		ConDate:        conDatePtr,
		ExpDate:        expDatePtr,
		Description:    req.Description,
		Files:          domainFiles,
	}

	if err := u.insRepo.CreateInsurance(ctx, ins); err != nil {
		return nil, err
	}

	return &pb.InsuranceResponse{
		Insurance: mapper.ToInsuranceProto(ins),
	}, nil
}

func (u *InsuranceUsecase) GetInsurance(ctx context.Context, req *pb.GetAssetRequest) (*pb.InsuranceArrayResponse, error) {
	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	accounts, err := u.insRepo.GetInsurance(ctx, uid)
	if err != nil {
		return nil, err
	}

	var InvestList []*pb.Insurance
	for _, item := range accounts {
		InvestList = append(InvestList, mapper.ToInsuranceProto(item))
	}

	return &pb.InsuranceArrayResponse{
		Success:   true,
		Insurance: InvestList,
	}, nil
}

func (u *InsuranceUsecase) GetInsuranceByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.InsuranceArrayResponse, error) {
	var ids []uuid.UUID
	for _, idStr := range req.Ids {
		if parsedID, err := uuid.Parse(idStr); err == nil {
			ids = append(ids, parsedID)
		}
	}

	bu, err := u.insRepo.GetInsuranceByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	var pbIns []*pb.Insurance
	for _, a := range bu {
		pbIns = append(pbIns, mapper.ToInsuranceProto(a))
	}

	return &pb.InsuranceArrayResponse{
		Insurance: pbIns,
	}, nil
}

func (u *InsuranceUsecase) GetInsuranceByID(ctx context.Context, req *pb.GetAssetByIDRequest) (*pb.InsuranceResponse, error) {
	id, uid, err := utils.ValidateIDs(req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	in, err := u.insRepo.GetInsuranceByID(ctx, id, uid)
	if err != nil {
		return nil, err
	}

	return &pb.InsuranceResponse{
		Success:   true,
		Insurance: mapper.ToInsuranceProto(in),
	}, nil
}

func (u *InsuranceUsecase) UpdateInsurance(ctx context.Context, req *pb.UpdateInsuranceRequest) (*pb.InsuranceResponse, error) {
	id, uid, err := utils.ValidateIDs(req.Id, req.Insurance.UserId)
	if err != nil {
		return nil, err
	}

	in, err := u.insRepo.GetInsuranceByID(ctx, id, uid)
	if err != nil {
		return nil, err
	}

	if err := helper.ApplyUpdateInsuranceFields(req, in); err != nil {
		return nil, err
	}

	syncParams := domain.FileSyncParams{
		UserID:        uid,
		EntityID:      id,
		EntityType:    "insurance",
		NewFiles:      req.NewFiles,
		DeleteFileIDs: req.DeleteFileIds,
	}

	err = helper.SyncEntityFiles(ctx, u.fileRepo, u.storage, syncParams)
	if err != nil {
		return nil, err
	}

	updatedInsurance, err := u.insRepo.UpdateInsurance(ctx, in)
	if err != nil {
		return nil, err
	}

	return &pb.InsuranceResponse{
		Success:   true,
		Insurance: mapper.ToInsuranceProto(updatedInsurance),
	}, nil
}

func (u *InsuranceUsecase) DeleteInsurance(ctx context.Context, req *pb.DeleteAssetRequest) (*pb.DeleteAssetResponse, error) {
	id, uid, err := utils.ValidateIDs(req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	existingAcc, err := u.insRepo.GetInsuranceByID(ctx, id, uid)
	if err != nil {
		return nil, err
	}

	if err := u.insRepo.DeleteInsurance(ctx, id, uid); err != nil {
		return nil, err
	}

	if len(existingAcc.Files) > 0 {
		fileURLs := make([]string, len(existingAcc.Files))
		for i, f := range existingAcc.Files {
			fileURLs[i] = f.Link
		}

		helper.DeleteFilesAsync(u.storage, fileURLs)
	}

	return &pb.DeleteAssetResponse{
		Success: true,
	}, nil
}
