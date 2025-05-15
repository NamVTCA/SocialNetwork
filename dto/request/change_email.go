package request

type ChangeEmailRequest struct {
	OldEmail string `json:"oldEmail" binding:"required,email"`
    Password string `json:"password" binding:"required"`
    NewEmail string `json:"newEmail" binding:"required,email"`
}
