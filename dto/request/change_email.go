package request

type ChangeEmailRequest struct {
	OldEmail string `json:"oldEmail" binding:"required,email"`
	NewEmail string `json:"newEmail" binding:"required,email"`
}

type VerifyEmailRequest struct {
	OTP      string `json:"otp" binding:"required"`
}
