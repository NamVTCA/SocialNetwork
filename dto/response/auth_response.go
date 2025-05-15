package response

type LoginResponse struct {
	Token      string `json:"access_token"`
	ExpiresIn  int    `json:"expires_in"`
	TokenType  string `json:"token_type"`
}