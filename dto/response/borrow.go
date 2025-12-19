package response

import "time"

type BorrowBookResponse struct {
	ID            uint64    `json:"id"`
	BookID        uint64    `json:"book_id"`
	BookTitle     string    `json:"book_title"`
	UserID        uint64    `json:"user_id"`
	Username      string    `json:"username"`
	BorrowDate    time.Time `json:"borrow_date"`
	DueDate       time.Time `json:"due_date"`
	Status        string    `json:"status"`
	RenewCount    int       `json:"renew_count"`
	MaxRenewCount int       `json:"max_renew_count"`
}

type ReturnBookResponse struct {
	ID          uint64    `json:"id"`
	BookID      uint64    `json:"book_id"`
	UserID      uint64    `json:"user_id"`
	BorrowDate  time.Time `json:"borrow_date"`
	DueDate     time.Time `json:"due_date"`
	ReturnDate  time.Time `json:"return_date"`
	Status      string    `json:"status"`
	IsOverdue   bool      `json:"is_overdue"`
	OverdueDays int       `json:"overdue_days"`
	Fine        float64   `json:"fine"`
	Condition   *string   `json:"condition"`
}

type RenewBorrowResponse struct {
	ID              uint64    `json:"id"`
	BookID          uint64    `json:"book_id"`
	BookTitle       string    `json:"book_title"`
	OriginalDueDate time.Time `json:"original_due_date"`
	NewDueDate      time.Time `json:"new_due_date"`
	RenewCount      int       `json:"renew_count"`
	MaxRenewCount   int       `json:"max_renew_count"`
}

type GetBorrowRecordListUserResponse struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
}

type GetBorrowRecordListBookResponse struct {
	ID       uint64 `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	ISBN     string `json:"isbn"`
	CoverURL string `json:"cover_url"`
}

type GetBorrowRecordItemResponse struct {
	ID           uint64                          `json:"id"`
	Book         GetBorrowRecordListBookResponse `json:"book"`
	User         GetBorrowRecordListUserResponse `json:"user"`
	BorrowDate   time.Time                       `json:"borrow_date"`
	DueDate      time.Time                       `json:"due_date"`
	ReturnDate   *time.Time                      `json:"return_date"`
	Status       string                          `json:"status"`
	IsOverdue    bool                            `json:"is_overdue"`
	DaysUntilDue int                             `json:"days_until_due"`
	OverdueDays  int                             `json:"overdue_days,omitempty"`
	RenewCount   int                             `json:"renew_count"`
	CanRenew     bool                            `json:"can_renew"`
	Fine         float64                         `json:"fine"`
}

type GetBorrowRecordListResponse struct {
	Total      int64                         `json:"total"`
	Page       int                           `json:"page"`
	Limit      int                           `json:"limit"`
	TotalPages int                           `json:"total_pages"`
	Records    []GetBorrowRecordItemResponse `json:"records"`
}
