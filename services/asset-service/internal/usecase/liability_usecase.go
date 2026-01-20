package usecase

import (
	"context"
	"errors"
	"time"
	"wealth-vault/asset-service/internal/domain"
	repo "wealth-vault/asset-service/internal/repository/interface"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	"wealth-vault/asset-service/pkg/utils"

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

func (u *LiabilityUsecase) CreateLiability(ctx context.Context, req *pb.CreateLiabilityRequest) (*pb.CreateLiabilityResponse, error) {
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
		Type:         domain.LiabilityType(req.Type.String()),
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

	return &pb.CreateLiabilityResponse{
		Success: true,
		Id:      liability.ID.String(),
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
		LiaList = append(LiaList, utils.ToLiabilityProto(item))
	}

	return &pb.LiabilityArrayResponse{
		Success:   true,
		Liability: LiaList,
	}, nil
}

func (u *LiabilityUsecase) GetLiabilityByID(ctx context.Context, req *pb.GetLiabilityByIDRequest) (*pb.LiabilityResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, errors.New("invalid asset id")
	}

	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	lia, err := u.liaRepo.GetLiabilityByID(ctx, id, uid)
	if err != nil {
		return nil, err
	}

	return &pb.LiabilityResponse{
		Success:   true,
		Liability: utils.ToLiabilityProto(lia),
	}, nil
}

func (u *LiabilityUsecase) UpdateLiability(ctx context.Context, req *pb.UpdateLiabilityRequest) (*pb.LiabilityResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, errors.New("invalid asset id")
	}

	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	lia, err := u.liaRepo.GetLiabilityByID(ctx, id, uid)
	if err != nil {
		return nil, err
	}

	var updateMask []string
	has := func(target string) bool {
		if req.UpdateMask == nil || len(req.UpdateMask.Paths) == 0 {
			return true
		}

		for _, p := range req.UpdateMask.Paths {
			if p == target {
				return true
			}
		}

		return false
	}

	if has("name") {
		lia.Name = req.Name
		updateMask = append(updateMask, "Name")
	}

	if has("creditor") {
		lia.Creditor = req.Creditor
		updateMask = append(updateMask, "Creditor")
	}

	if has("principal") {
		lia.Principal = req.Principal
		updateMask = append(updateMask, "Principal")
	}

	if has("interest_rate") {
		lia.InterestRate = req.InterestRate
		updateMask = append(updateMask, "InterestRate")
	}

	if has("description") {
		lia.Description = req.Description
		updateMask = append(updateMask, "Description")
	}

	if has("started_at") && req.StartAt != nil {
		t := req.StartAt.AsTime()
		lia.StartAt = &t
		updateMask = append(updateMask, "StartAt")
	}

	if has("ended_at") && req.EndAt != nil {
		t := req.EndAt.AsTime()
		lia.EndAt = &t
		updateMask = append(updateMask, "EndAt")
	}

	if len(req.DeleteFileIds) > 0 {
		var fileUUIDs []uuid.UUID
		for _, idStr := range req.DeleteFileIds {
			parsedID, err := uuid.Parse(idStr)
			if err == nil {
				fileUUIDs = append(fileUUIDs, parsedID)
			}
		}

		filesToDelete, err := u.fileRepo.GetFilesByIDs(ctx, fileUUIDs)
		if err == nil {
			var validURLs []string
			for _, f := range filesToDelete {
				if f.UserID == uid {
					validURLs = append(validURLs, f.Link)
				}
			}

			if len(validURLs) > 0 {
				utils.DeleteFilesAsync(u.storage, validURLs)
			}
		}

		err = u.fileRepo.DeleteFiles(ctx, req.DeleteFileIds, id, uid)
		if err != nil {
			return nil, err
		}
	}

	if len(req.NewFiles) > 0 {
		var filesToCreate []domain.FileAssociate

		for _, f := range req.NewFiles {
			filesToCreate = append(filesToCreate, domain.FileAssociate{
				ID:         uuid.New(),
				EntityID:   id,
				EntityType: "liability",
				UserID:     lia.UserID,
				Link:       f.Url,
				FileType:   f.FileType,
			})
		}

		if err := u.fileRepo.CreateFiles(ctx, filesToCreate); err != nil {
			return nil, err
		}
	}

	updatedLia, err := u.liaRepo.UpdateLiability(ctx, lia, updateMask)
	if err != nil {
		return nil, err
	}

	return &pb.LiabilityResponse{
		Success:   true,
		Liability: utils.ToLiabilityProto(updatedLia),
	}, nil
}

func (u *LiabilityUsecase) DeleteLiability(ctx context.Context, req *pb.DeleteLiabilityRequest) (*pb.DeleteLiabilityResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, errors.New("invalid asset id")
	}

	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.New("invalid user id")
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

		utils.DeleteFilesAsync(u.storage, fileURLs)
	}

	return &pb.DeleteLiabilityResponse{
		Success: true,
	}, nil
}
