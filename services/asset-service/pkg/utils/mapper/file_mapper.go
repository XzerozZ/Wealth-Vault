package mapper

import (
	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"

	"github.com/google/uuid"
)

func ToDomainFiles(pbFiles []*pb.FileInfo, userID uuid.UUID) []domain.FileAssociate {
	if len(pbFiles) == 0 {
		return nil
	}
	files := make([]domain.FileAssociate, len(pbFiles))
	for i, f := range pbFiles {
		files[i] = domain.FileAssociate{
			Link:     f.Url,
			FileType: f.FileType,
			UserID:   userID,
		}
	}
	return files
}
