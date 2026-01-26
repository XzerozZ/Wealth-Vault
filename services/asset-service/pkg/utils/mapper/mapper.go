package mapper

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
			Id:       f.ID.String(),
			Url:      f.Link,
			FileType: f.FileType,
		}
	}

	return pbFiles
}

func ToRefLandProto(land []domain.Land) []*pb.Reference {
	if len(land) == 0 {
		return []*pb.Reference{}
	}

	ref := make([]*pb.Reference, len(land))
	for i, f := range land {
		ref[i] = &pb.Reference{
			Id:   f.ID.String(),
			Name: f.Name,
		}
	}

	return ref
}

func ToRefBuildingProto(build []domain.Building) []*pb.Reference {
	if len(build) == 0 {
		return []*pb.Reference{}
	}

	ref := make([]*pb.Reference, len(build))
	for i, f := range build {
		ref[i] = &pb.Reference{
			Id:   f.ID.String(),
			Name: f.Name,
		}
	}

	return ref
}
