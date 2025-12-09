package response

import(
	"time"
)

type UserRegisterResponse struct {
	ID int 				`json:"id"`
	Username string 	`json:"username"`
	Email string 		`json:"email"`
	Role string 		`json:"role"`
	CreatedAt time.Time `json:"created_at"`
}