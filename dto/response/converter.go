package response

import "socialnetwork/models"

func ToUserDetailResponse(u *models.User) *UserDetailResponse {
	return &UserDetailResponse{
		ID:          u.ID.Hex(),
		Username:    u.Username,
		Email:       u.Email,
		DisplayName: u.DisplayName,
		AvatarURL:   u.AvatarURL,
		CoverURL:    u.CoverURL,
		Bio:         u.Bio,
		Gender:      string(u.Gender),
		Phone:       u.Phone,
		BirthDate:      u.BirthDate,
		Location:       u.Location,
		Website:        u.Website,
		IsVerified:     u.IsVerified,
		FollowerCount:  u.FollowerCount,
		FollowingCount: u.FollowingCount,
		CreatedAt:      u.CreatedAt,
	}
}
