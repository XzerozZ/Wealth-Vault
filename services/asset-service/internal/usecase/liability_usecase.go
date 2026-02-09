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

type LiabilityUsecase struct {
	liaRepo    repo.LiabilityRepository
	fileRepo   repo.FileRepository
	storage    *utils.StorageClient
	userClient userPb.UserServiceClient
}

func NewLiabilityUsecase(r repo.LiabilityRepository, fr repo.FileRepository, s *utils.StorageClient, userClient userPb.UserServiceClient) LiabilityUsecase {
	return LiabilityUsecase{
		liaRepo:    r,
		fileRepo:   fr,
		storage:    s,
		userClient: userClient,
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

func (u *LiabilityUsecase) GetLiabilityByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.LiabilityArrayResponse, error) {
	var ids []uuid.UUID
	for _, idStr := range req.Ids {
		if parsedID, err := uuid.Parse(idStr); err == nil {
			ids = append(ids, parsedID)
		}
	}

	bu, err := u.liaRepo.GetLiabilityByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	var pbLia []*pb.Liability
	for _, a := range bu {
		pbLia = append(pbLia, mapper.ToLiabilityProto(a))
	}

	return &pb.LiabilityArrayResponse{
		Liability: pbLia,
	}, nil
}

func (u *LiabilityUsecase) GetBatchLiabilityByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.LiabilityArrayResponse, error) {
	var ids []uuid.UUID
	for _, idStr := range req.Ids {
		if parsedID, err := uuid.Parse(idStr); err == nil {
			ids = append(ids, parsedID)
		}
	}

	bu, err := u.liaRepo.GetBatchLiabilityByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	var pbLia []*pb.Liability
	for _, a := range bu {
		pbLia = append(pbLia, mapper.ToLiabilityProto(a))
	}

	return &pb.LiabilityArrayResponse{
		Liability: pbLia,
	}, nil
}

func (u *LiabilityUsecase) GetLiabilityByID(ctx context.Context, req *pb.GetLiabilityByIDRequest) (*pb.LiabilityResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	lia, err := u.liaRepo.GetLiabilityByID(ctx, id)
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

	lia, err := u.liaRepo.GetLiabilityByID(ctx, id)
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

	_, err = u.liaRepo.GetLiabilityByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := u.liaRepo.SoftDeleteLiability(ctx, id, uid); err != nil {
		return nil, err
	}

	return &pb.DeleteLiabilityResponse{
		Success: true,
	}, nil
}

func (u *LiabilityUsecase) CleanupExpiredLiability(ctx context.Context) error {
	cutoffTime := time.Now().AddDate(0, 0, -7)
	GetExpiredLia, err := u.liaRepo.GetExpiredLiability(ctx, cutoffTime)
	if err != nil {
		return err
	}

	if len(GetExpiredLia) == 0 {
		return err
	}

	for _, l := range GetExpiredLia {
		helper.CleanupAssetResource(
			ctx,
			l.ID,
			l.Files,
			u.storage,
			u.userClient,
			func(id uuid.UUID) error {
				return u.liaRepo.HardDeleteLiability(ctx, id)
			},
		)
	}

	return nil
}
