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