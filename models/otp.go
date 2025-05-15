package models

type SendOTPRequest struct {
    Identifier string `json:"identifier" binding:"required"` // email hoặc phone
    Channel    string `json:"channel" binding:"required,oneof=email phone"`
    Purpose    string `json:"purpose" binding:"required"`
}

type VerifyOTPRequest struct {
    Identifier string `json:"identifier" binding:"required"`
    Purpose    string `json:"purpose" binding:"required"`
    OTP        string `json:"otp" binding:"required,len=6"`
}

type OTP struct {
    Identifier string
    Purpose    string
    Code       string
    ExpiredAt  int64
}
