package model

import(
	"time"
)

type User struct {
	ID 			uint64 	  	`json:"id" gorm:"primaryKey;autoIncrement"`
	Username 	string 		`json:"username" gorm:"type:varchar(20);unique;not null"`
	Password 	string 		`json:"password" gorm:"type:varchar(255);not null"`
	Email 		string 		`json:"email" gorm:"type:varchar(100);unique;not null"`
	Phone 		string 		`json:"phone" gorm:"varchar(11)"`
	Role 		string 		`json:"role" gorm:"type:enum('admin','user');default:'user'"`
	Status		string 		`json:"status" gorm:"type:enum('active','disabled');default:'active'"`
	BorrowLimit int	   		`json:"borrow_limit" gorm:"defalut:5"`
	CreatedAt	time.Time	`json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time	`json:"updated_at" gorm:"autoUpdateTime"`
}
