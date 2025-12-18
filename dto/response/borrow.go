package response

import "time"

type BorrowBookResponse struct {
    ID             uint64    	`json:"id"`
    BookID         uint64    	`json:"book_id"`
    BookTitle      string    	`json:"book_title"`
    UserID         uint64    	`json:"user_id"`
    Username       string    	`json:"username"`
    BorrowDate     time.Time    `json:"borrow_date"`
    DueDate        time.Time    `json:"due_date"`
    Status         string    	`json:"status"`
    RenewCount     int      	`json:"renew_count"`
    MaxRenewCount  int       	`json:"max_renew_count"`
}

type ReturnBookResponse struct {
	ID           uint64    `json:"id"`
	BookID       uint64    `json:"book_id"`
	UserID       uint64    `json:"user_id"`
	BorrowDate   time.Time `json:"borrow_date"`
	DueDate      time.Time `json:"due_date"`
	ReturnDate   time.Time `json:"return_date"`
	Status       string    `json:"status"`
	IsOverdue    bool      `json:"is_overdue"`
	OverdueDays  int       `json:"overdue_days"`
	Fine         float64   `json:"fine"`
	Condition    *string    `json:"condition"`
}
