package usecase

import (
	"context"
	"errors"
	"log"
	"time"
	"wealth-vault/asset-service/internal/domain"
	"wealth-vault/asset-service/internal/event"
	repo "wealth-vault/asset-service/internal/repository/interface"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	"wealth-vault/asset-service/pkg/utils"
	helper "wealth-vault/asset-service/pkg/utils/helper"
	"wealth-vault/asset-service/pkg/utils/mapper"

	"github.com/google/uuid"
)

type InsuranceUsecase struct {
	insRepo     repo.InsuranceRepository
	assetHelper helper.AssetHelper
	publisher   event.EventPublisher
}

func NewInsuranceUsecase(r repo.InsuranceRepository, ah helper.AssetHelper, e event.EventPublisher) InsuranceUsecase {
	return InsuranceUsecase{
		insRepo:     r,
		assetHelper: ah,
		publisher:   e,
	}
}

func (u *InsuranceUsecase) CreateInsurance(ctx context.Context, req *pb.CreateInsuranceRequest) (*pb.InsuranceResponse, error) {
	uid, err := utils.ParseID(req.UserId)
	if err != nil {
		return nil, err
	}

	ins := mapper.ToInsuranceDomain(req, uid)

	if err := u.insRepo.CreateInsurance(ctx, ins); err != nil {
		return nil, err
	}

	return &pb.InsuranceResponse{
		Insurance: mapper.ToInsuranceProto(ins),
	}, nil
}

func (u *InsuranceUsecase) GetInsurance(ctx context.Context, req *pb.GetAssetRequest) (*pb.InsuranceArrayResponse, error) {
	uid, err := utils.ParseID(req.UserId)
	if err != nil {
		return nil, err
	}

	insurances, err := u.insRepo.GetInsurance(ctx, uid)
	if err != nil {
		return nil, err
	}

	return &pb.InsuranceArrayResponse{
		Success:   true,
		Insurance: mapper.ToInsuranceProtoSlice(insurances),
	}, nil
}

func (u *InsuranceUsecase) GetInsuranceByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.InsuranceArrayResponse, error) {
	ids := utils.ParseUUIDs(req.Ids)

	insurances, err := u.insRepo.GetInsuranceByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return &pb.InsuranceArrayResponse{
		Insurance: mapper.ToInsuranceProtoSlice(insurances),
	}, nil
}

func (u *InsuranceUsecase) GetBatchInsuranceByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.InsuranceArrayResponse, error) {
	ids := utils.ParseUUIDs(req.Ids)

	insurances, err := u.insRepo.GetBatchInsuranceByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return &pb.InsuranceArrayResponse{
		Insurance: mapper.ToInsuranceProtoSlice(insurances),
	}, nil
}

func (u *InsuranceUsecase) GetInsuranceByID(ctx context.Context, req *pb.GetAssetByIDRequest) (*pb.InsuranceResponse, error) {
	id, err := utils.ParseID(req.Id)
	if err != nil {
		return nil, err
	}

	in, err := u.insRepo.GetInsuranceByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.InsuranceResponse{
		Success:   true,
		Insurance: mapper.ToInsuranceProto(in),
	}, nil
}

func (u *InsuranceUsecase) UpdateInsurance(ctx context.Context, req *pb.UpdateInsuranceRequest) (*pb.InsuranceResponse, error) {
	if req.Insurance == nil {
		return nil, errors.New("insurance data is required")
	}

	id, uid, err := utils.ValidateIDs(req.Id, req.Insurance.UserId)
	if err != nil {
		return nil, err
	}

	in, err := u.insRepo.GetInsuranceByID(ctx, id)
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

	if err := u.assetHelper.SyncFiles(ctx, syncParams); err != nil {
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

	if _, err = u.insRepo.GetInsuranceByID(ctx, id); err != nil {
		return nil, err
	}

	if err := u.insRepo.SoftDeleteInsurances(ctx, id, uid); err != nil {
		return nil, err
	}

	return &pb.DeleteAssetResponse{
		Success: true,
	}, nil
}

func (u *InsuranceUsecase) CleanupExpiredInsurances(ctx context.Context) error {
	cutoffTime := time.Now().AddDate(0, 0, -7)
	expiredInsurances, err := u.insRepo.GetExpiredInsurances(ctx, cutoffTime)
	if err != nil {
		return err
	}

	if len(expiredInsurances) == 0 {
		return nil
	}

	for _, ins := range expiredInsurances {
		u.assetHelper.CleanupResource(ctx, ins.ID, ins.Files, func(id uuid.UUID) error {
			return u.insRepo.HardDeleteInsurances(ctx, id)
		})
	}

	return nil
}

func (u *InsuranceUsecase) CheckExpiringInsurances(ctx context.Context) error {
	checkDays := []int{30, 21, 14, 7, 1}
	for _, days := range checkDays {
		insurances, err := u.insRepo.GetExpiringInsurances(ctx, days)
		if err != nil {
			log.Printf("❌ Error checking insurances for %d days: %v", days, err)
			continue
		}

		for _, ins := range insurances {

			expDateStr := "ไม่ระบุวันที่"
			if ins.ExpDate != nil {
				expDateStr = ins.ExpDate.Format("2006-01-02")
			}

			evt := domain.InsuranceExpiringEvent{
				UserID:        ins.UserID.String(),
				InsuranceID:   ins.ID.String(),
				InsuranceName: ins.Name,
				DaysLeft:      days,
				ExpDate:       expDateStr,
			}

			topic := "noti.insurance.expiring"
			if err := u.publisher.Publish(topic, evt); err != nil {
				log.Printf("❌ Failed to publish event: %v", err)
			}
		}
	}

	return nil
}
