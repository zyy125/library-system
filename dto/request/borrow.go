package request

type BorrowBookRequest struct {
	BookId uint64 `json:"book_id" binding:"required"`
	BorrowDays *int `json:"borrow_days"`
}