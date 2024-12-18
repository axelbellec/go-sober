package dtos

type UserSignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSignupResponse struct {
	Message string `json:"message"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type UserMeResponse struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
}
