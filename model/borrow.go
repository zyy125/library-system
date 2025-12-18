package model

import(
	"time"
)

type BorrowRecord struct {
    ID         uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
    BookID     uint64     `json:"book_id" gorm:"index:idx_book;not null"`
    UserID     uint64     `json:"user_id" gorm:"index:idx_user;not null"`
    BorrowDate time.Time  `json:"borrow_date" gorm:"autoCreateTime"`
    DueDate    time.Time  `json:"due_date" gorm:"not null;index:idx_due_date"`
    ReturnDate *time.Time `json:"return_date,omitempty" gorm:"index:idx_return_date"`
    Status     string     `json:"status" gorm:"type:enum('borrowed','returned','overdue');default:'borrowed';index:idx_status"`
    RenewCount int        `json:"renew_count" gorm:"default:0"`
    Fine       float64    `json:"fine" gorm:"type:decimal(10,2);default:0"`
    CreatedAt  time.Time  `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt  time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

    Book Book `gorm:"foreignKey:BookID"`
    User User `gorm:"foreignKey:UserID"`
}
