package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Gender string

const (
	GenderMale    Gender = "male"
	GenderFemale  Gender = "female"
	GenderOther   Gender = "other"
	GenderPrivate Gender = "private"
)

type User struct {
	// ID chính và thông tin đăng nhập
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username      string             `bson:"username" json:"username" validate:"required,alphanum,min=3,max=30"`
	Email         string             `bson:"email" json:"email" validate:"required,email"`
	Phone         string             `bson:"phone,omitempty" json:"phone,omitempty" validate:"omitempty,e164"`
	Password      string             `bson:"password,omitempty" json:"password,omitempty"`
	PhoneVerified bool               `bson:"phoneVerified" json:"phoneVerified"`
	EmailVerified bool               `bson:"emailVerified" json:"emailVerified"`

	// Thông tin hồ sơ cá nhân
	DisplayName string     `bson:"displayName,omitempty" json:"displayName,omitempty" validate:"omitempty,max=100"`
	AvatarURL   string     `bson:"avatarUrl,omitempty" json:"avatarUrl,omitempty" validate:"omitempty,url"`
	CoverURL    string     `bson:"coverUrl,omitempty" json:"coverUrl,omitempty" validate:"omitempty,url"`
	Bio         string     `bson:"bio,omitempty" json:"bio,omitempty" validate:"omitempty,max=500"`
	Gender      Gender     `bson:"gender,omitempty" json:"gender,omitempty"`
	BirthDate   *time.Time `bson:"birthDate,omitempty" json:"birthDate,omitempty"`
	Location    string     `bson:"location,omitempty" json:"location,omitempty"`
	Website     string     `bson:"website,omitempty" json:"website,omitempty" validate:"omitempty,url"`

	// Social Graph (quan hệ xã hội)
	Followers []primitive.ObjectID `bson:"followers,omitempty" json:"followers,omitempty"`
	Following []primitive.ObjectID `bson:"following,omitempty" json:"following,omitempty"`

	// Đếm follower/following để truy vấn nhanh, tránh count nhiều lần
	FollowerCount  int `bson:"followerCount" json:"followerCount"`
	FollowingCount int `bson:"followingCount" json:"followingCount"`

	// **Bạn bè & chặn**
	FriendRequests []primitive.ObjectID `bson:"friendRequests,omitempty" json:"friendRequests,omitempty"` // đang chờ xác nhận
	Friends        []primitive.ObjectID `bson:"friends,omitempty" json:"friends,omitempty"`               // đã xác nhận
	BlockedUsers   []primitive.ObjectID `bson:"blockedUsers,omitempty" json:"blockedUsers,omitempty"`     // đã chặn

	// Ẩn/hiện trang cá nhân
	HideProfile bool `bson:"hideProfile" json:"hideProfile"`

	// Bảo mật & quyền hạn
	Roles      []string   `bson:"roles,omitempty" json:"roles,omitempty"` // vd: ["user","admin"]
	IsVerified bool       `bson:"isVerified" json:"isVerified"`
	IsActive   bool       `bson:"isActive" json:"isActive"`
	LastLogin  *time.Time `bson:"lastLogin,omitempty" json:"lastLogin,omitempty"`

	// Timestamps
	CreatedAt time.Time  `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time  `bson:"updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"` // soft delete

	// Dùng cho tìm kiếm full-text (nếu có)
	SearchTerms []string `bson:"searchTerms,omitempty" json:"-"`
}
