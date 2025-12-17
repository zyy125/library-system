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

type ChangePwdRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=8,max=32"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=32"`
}

type GetUserListRequest struct {
	Page 		int 	`form:"page"`
	Limit 		int 	`form:"limit"`
	Username 	string 	`form:"username"`
	Role 		string 	`form:"role"`
	Status 		string 	`form:"status"`
}

type CreateUserRequest struct {
	Username    string `json:"username" binding:"required,min=4,max=20"`
	Password    string `json:"password" binding:"required,min=8,max=32"`
	Email       string `json:"email"    binding:"required,email"`
	Phone       string `json:"phone"`
	Role        string `json:"role"`
	BorrowLimit int    `json:"borrow_limit"`
}

type UpdateUserByAdminRequest struct {
	Username 	*string `json:"username" binding:"omitempty,min=4,max=20,alphanumunicode"`
	Email 	 	*string `json:"email" binding:"omitempty,email"`
	Phone    	*string `json:"phone" binding:"omitempty,len=11,numeric"`
	Role        *string `json:"role"`
	Status      *string `json:"status"`
	BorrowLimit *int    `json:"borrow_limit"`
}
