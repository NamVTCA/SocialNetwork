package request

import "time"

type UpdateProfileRequest struct {
	DisplayName string     `json:"displayName,omitempty"`
	Bio         string     `json:"bio,omitempty"`
	Email       string     `json:"email,omitempty"`
	Gender      string     `json:"gender,omitempty"`
	BirthDate   *time.Time `json:"birthDate,omitempty"`
	AvatarURL   string     `json:"avatarUrl,omitempty"`
	CoverURL    string     `json:"coverUrl,omitempty"`
	Location    string     `json:"location,omitempty"`
	Website     string     `json:"website,omitempty"`
	Phone       string     `json:"phone,omitempty"`
}

type UpdateSecurityRequest struct {
    Email       string `json:"email,omitempty"`
    Password    string `json:"password,omitempty"`
    OTP         string `json:"otp" binding:"required"`
    OTPChannel  string `json:"otpChannel" binding:"required,oneof=email phone"`
}

