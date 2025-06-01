package request

type RegisterRequest struct {
	Name       string `json:"name" binding:"required,min=2,max=100"`
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required,min=8,max=32"`
}