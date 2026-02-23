package helper

import (
	"wealth-vault/user-service/internal/domain"
	"wealth-vault/user-service/internal/infra/storage"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
)

func ApplyUpdateGroupFields(req *pb.UpdateGroupRequest, storage storage.SupabaseStorage, group *domain.Group) ([]string, error) {
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

	if has("group_name") {
		group.GroupName = req.Name
		updateMask = append(updateMask, "GroupName")
	}

	if has("profile") {
		if group.GroupProfile != "" && group.GroupProfile != req.Profile {
			DeleteFilesAsync(storage, []string{group.GroupProfile})
		}
		group.GroupProfile = req.Profile
		updateMask = append(updateMask, "GroupProfile")
	}

	return updateMask, nil
}
