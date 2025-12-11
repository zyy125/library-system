package request

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=4,max=20,alphanumunicode"`
	Password string `json:"password" binding:"required,min=8,max=32"`
	Email 	 string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"omitempty,len=11,numeric"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required,min=4,max=20,alphanumunicode"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

type UserRefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type UpdateUserRequest struct {
	Username *string `json:"username" binding:"omitempty,min=4,max=20,alphanumunicode"`
	Email 	 *string `json:"email" binding:"omitempty,email"`
	Phone    *string `json:"phone" binding:"omitempty,len=11,numeric"`
}