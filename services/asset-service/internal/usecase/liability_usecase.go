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

type LiabilityUsecase struct {
	liaRepo  repo.LiabilityRepository
	fileRepo repo.FileRepository
	storage  *utils.StorageClient
}

func NewLiabilityUsecase(r repo.LiabilityRepository, fr repo.FileRepository, s *utils.StorageClient) LiabilityUsecase {
	return LiabilityUsecase{
		liaRepo:  r,
		fileRepo: fr,
		storage:  s,
	}
}

func (u *LiabilityUsecase) CreateLiability(ctx context.Context, req *pb.CreateLiabilityRequest) (*pb.LiabilityResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	var startAt *time.Time
	if req.StartAt != nil {
		t := req.StartAt.AsTime()
		startAt = &t
	}

	var endAt *time.Time
	if req.EndAt != nil {
		t := req.EndAt.AsTime()
		endAt = &t
	}

	liType := domain.LiabilityTypeLoan
	if val, ok := helper.ProtoToDomainLiType[req.Type]; ok {
		liType = val
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

	liability := &domain.Liability{
		UserID:       userID,
		Type:         liType,
		Name:         req.Name,
		Creditor:     req.Creditor,
		Principal:    req.Principal,
		InterestRate: req.InterestRate,
		Description:  req.Description,
		StartAt:      startAt,
		EndAt:        endAt,
		Files:        domainFiles,
	}

	if err := u.liaRepo.CreateLiability(ctx, liability); err != nil {
		return nil, err
	}

	return &pb.LiabilityResponse{
		Success:   true,
		Liability: mapper.ToLiabilityProto(liability),
	}, nil
}

func (u *LiabilityUsecase) GetLiability(ctx context.Context, req *pb.GetLiabilityRequest) (*pb.LiabilityArrayResponse, error) {
	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	lias, err := u.liaRepo.GetLiability(ctx, uid)
	if err != nil {
		return nil, err
	}

	var LiaList []*pb.Liability
	for _, item := range lias {
		LiaList = append(LiaList, mapper.ToLiabilityProto(item))
	}

	return &pb.LiabilityArrayResponse{
		Success:   true,
		Liability: LiaList,
	}, nil
}

func (u *LiabilityUsecase) GetLiabilityByID(ctx context.Context, req *pb.GetLiabilityByIDRequest) (*pb.LiabilityResponse, error) {
	id, uid, err := utils.ValidateIDs(req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	lia, err := u.liaRepo.GetLiabilityByID(ctx, id, uid)
	if err != nil {
		return nil, err
	}

	return &pb.LiabilityResponse{
		Success:   true,
		Liability: mapper.ToLiabilityProto(lia),
	}, nil
}

func (u *LiabilityUsecase) UpdateLiability(ctx context.Context, req *pb.UpdateLiabilityRequest) (*pb.LiabilityResponse, error) {
	id, uid, err := utils.ValidateIDs(req.Id, req.Liability.CreatedBy)
	if err != nil {
		return nil, err
	}

	lia, err := u.liaRepo.GetLiabilityByID(ctx, id, uid)
	if err != nil {
		return nil, err
	}

	if err := helper.ApplyUpdateFields(req, lia); err != nil {
		return nil, err
	}

	syncParams := domain.FileSyncParams{
		UserID:        uid,
		EntityID:      id,
		EntityType:    "liability",
		NewFiles:      req.NewFiles,
		DeleteFileIDs: req.DeleteFileIds,
	}

	err = helper.SyncEntityFiles(ctx, u.fileRepo, u.storage, syncParams)
	if err != nil {
		return nil, err
	}

	updatedLia, err := u.liaRepo.UpdateLiability(ctx, lia)
	if err != nil {
		return nil, err
	}

	return &pb.LiabilityResponse{
		Success:   true,
		Liability: mapper.ToLiabilityProto(updatedLia),
	}, nil
}

func (u *LiabilityUsecase) DeleteLiability(ctx context.Context, req *pb.DeleteLiabilityRequest) (*pb.DeleteLiabilityResponse, error) {
	id, uid, err := utils.ValidateIDs(req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	existingCash, err := u.liaRepo.GetLiabilityByID(ctx, id, uid)
	if err != nil {
		return nil, err
	}

	if err := u.liaRepo.DeleteLiability(ctx, id, uid); err != nil {
		return nil, err
	}

	if len(existingCash.Files) > 0 {
		fileURLs := make([]string, len(existingCash.Files))
		for i, f := range existingCash.Files {
			fileURLs[i] = f.Link
		}

		helper.DeleteFilesAsync(u.storage, fileURLs)
	}

	return &pb.DeleteLiabilityResponse{
		Success: true,
	}, nil
}
