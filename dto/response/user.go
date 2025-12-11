package response

import(
	"time"
)

type UserRegisterResponse struct {
	ID 				uint64 			`json:"id"`
	Username	 	string 			`json:"username"`
	Email 			string 			`json:"email"`
	Role 			string 			`json:"role"`
	CreatedAt 		time.Time 		`json:"created_at"`
}

type UserResponse struct {
	ID 				uint64 			`json:"id"`
	Username 		string 			`json:"username"`
	Email 			string 			`json:"email"`
	Role 			string 			`json:"role"`
}

type UserLoginResponse struct {
	AccessToken 	string 			`json:"access_token"`
	RefreshToken 	string 			`json:"refresh_token"`
	TokenType 		string 			`json:"token_type"`
	ExpiresIn 		int				`json:"expires_in"`
	User 			UserResponse   	`josn:"user"`
}

type UserTokenRefreshResponse struct {
	AccessToken 	string 			`json:"access_token"`
	RefreshToken 	string 			`json:"refresh_token"`
	TokenType 		string 			`json:"token_type"`
	ExpiresIn 		int				`json:"expires_in"`
	User 			UserResponse   	`josn:"user"`
}

type GetUserMsgResponse struct {
	ID 				uint64 	  	`json:"id"`
	Username 		string 		`json:"username"`
	Email 			string 		`json:"email"`
	Phone 			string 		`json:"phone"`
	Role 			string 		`json:"role"`
	Status			string 		`json:"status"`
	BorrowLimit 	int	   		`json:"borrow_limit"`
	BorrowingCount 	int			`json:"borrowing_count"`
    OverdueCount 	int			`json:"overdue_count"`
	CreatedAt		time.Time	`json:"created_at"`
}

type UpdateUserResponse struct {
	ID 				uint64 			`json:"id"`
	Username 		string 			`json:"username,omitempty"`
	Email 			string 			`json:"email,omitempty"`
	Phone 			string 			`json:"phone,omitempty"`
	UpdatedAt   	string			`json:"updated_at"`
}