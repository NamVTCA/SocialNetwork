package request

type SendOTPRequest struct {
	Identifier string `json:"identifier" binding:"required"` // email hoặc phone
	Channel    string `json:"channel" binding:"required,oneof=email phone"`
	Purpose    string `json:"purpose" binding:"required"` // eg: verify_email, reset_password
}
