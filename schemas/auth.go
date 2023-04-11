package schemas

type LoginFormRequest struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type LogoutResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
