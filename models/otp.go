package models

type SendOTPRequest struct {
	Identifier string `json:"identifier" binding:"required"` // email hoặc phone
	Channel    string `json:"channel" binding:"required,oneof=email phone"`
	Purpose    string `json:"purpose" binding:"required"`
	CustomKey  string `json:"custom_key,omitempty"` // key tuỳ chỉnh cho từng trường hợp
}

type VerifyOTPRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Purpose    string `json:"purpose" binding:"required"`
	OTP        string `json:"otp" binding:"required,len=6"`
	Channel    string `json:"channel" binding:"required"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	OTP         string `json:"otp" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type OTP struct {
	Identifier string
	Purpose    string
	Code       string
	ExpiredAt  int64
}
