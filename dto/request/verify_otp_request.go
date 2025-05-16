package request

type VerifyOTPRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Purpose    string `json:"purpose" binding:"required"`
	OTP        string `json:"otp" binding:"required,len=6"`
	Channel    string `json:"channel" binding:"required"`
}
