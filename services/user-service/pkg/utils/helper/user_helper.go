package helper

import (
	"wealth-vault/user-service/internal/domain"
	"wealth-vault/user-service/internal/infra/storage"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
)

func ApplyUpdateUserFields(req *pb.UpdateUserRequest, storage storage.SupabaseStorage, user *domain.User) ([]string, error) {
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

	if has("first_name") {
		user.Firstname = req.Firstname
		updateMask = append(updateMask, "Firstname")
	}

	if has("last_name") {
		user.Lastname = req.Lastname
		updateMask = append(updateMask, "Lastname")
	}

	if has("username") {
		user.Username = req.Username
		updateMask = append(updateMask, "Username")
	}

	if has("profile") {
		if user.Profile != "" && user.Profile != req.Profile {
			DeleteFilesAsync(storage, []string{user.Profile})
		}
		user.Profile = req.Profile
		updateMask = append(updateMask, "Profile")
	}

	if has("phonenumber") {
		user.Phonenumber = req.Phonenumber
		updateMask = append(updateMask, "Phonenumber")
	}

	if has("birthday") && req.Birthday != nil {
		t := req.Birthday.AsTime()
		user.Birthday = &t
		updateMask = append(updateMask, "Birthday")
	}

	if has("share_enabled") && req.Sharedenabled != nil {
		user.IsAutoShareEnabled = *req.Sharedenabled
		updateMask = append(updateMask, "AutoShareAge")
	}

	if has("shared_age") && req.Sharedage != nil {
		user.AutoShareAge = int(*req.Sharedage)
		updateMask = append(updateMask, "IsAutoShareEnabled")
	}

	return updateMask, nil
}
