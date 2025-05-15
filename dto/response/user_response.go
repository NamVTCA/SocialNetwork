package response

import (
	"time"
	
	)

// UserDetailResponse - Cho thông tin chi tiết cá nhân
type UserDetailResponse struct {
	ID           string     `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	Phone        string     `json:"phone,omitempty"`
	DisplayName  string     `json:"displayName"`
	AvatarURL    string     `json:"avatarUrl"`
	CoverURL     string     `json:"coverUrl,omitempty"`
	Bio          string     `json:"bio,omitempty"`
	Gender       string     `json:"gender,omitempty"`
	BirthDate    *time.Time `json:"birthDate,omitempty"`
	Location     string     `json:"location,omitempty"`
	Website      string     `json:"website,omitempty"`
	IsVerified   bool       `json:"isVerified"`
	FollowerCount int       `json:"followerCount"`
	FollowingCount int      `json:"followingCount"`
	CreatedAt    time.Time  `json:"createdAt"`
}

// UserPublicResponse - Cho thông tin công khai
type UserPublicResponse struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	AvatarURL   string `json:"avatarUrl"`
	Bio         string `json:"bio,omitempty"`
	FollowerCount int  `json:"followerCount"`
}

// UserListResponse - Cho danh sách user
type UserListResponse struct {
	Users      []UserPublicResponse `json:"users"`
	Total      int                  `json:"total"`
	Page       int                  `json:"page"`
	Limit      int                  `json:"limit"`
}