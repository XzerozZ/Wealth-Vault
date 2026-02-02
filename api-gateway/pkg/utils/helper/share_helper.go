package helper

import (
	"time"
	"wealth-vault/api-gateway/internal/domain"
	pb "wealth-vault/api-gateway/pkg/pb/proto/user"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func MapShareTargets(targets []domain.ShareTarget) []*pb.ShareTarget {
	var pbTargets []*pb.ShareTarget

	for _, t := range targets {
		var shareAtProto *timestamppb.Timestamp

		if t.ShareAt != "" {
			parsedTime, err := time.Parse("2006-01-02", t.ShareAt)

			if err == nil {
				shareAtProto = timestamppb.New(parsedTime)
			}
		}

		pbTargets = append(pbTargets, &pb.ShareTarget{
			Id:      t.ID,
			ShareAt: shareAtProto,
		})
	}

	return pbTargets
}
