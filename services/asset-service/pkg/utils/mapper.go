package utils

import (
	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
)

func ToPbFiles(files []domain.FileAssociate) []*pb.FileInfo {
	if len(files) == 0 {
		return []*pb.FileInfo{}
	}

	pbFiles := make([]*pb.FileInfo, len(files))

	for i, f := range files {
		pbFiles[i] = &pb.FileInfo{
			Id:       f.ID,
			Url:      f.Link,
			FileType: f.FileType,
		}
	}

	return pbFiles
}
