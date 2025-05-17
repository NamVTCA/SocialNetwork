package request

type ForgotPasswordRequest struct {
	Email string `json:"identifier" binding:"required"` // email hoặc phone
}

type ResetPasswordRequest struct {
	Email string `json:"identifier" binding:"required"` // email hoặc phone
	OTP         string `json:"otp" binding:"required,len=6"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
