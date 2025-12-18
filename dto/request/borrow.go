package request

type BorrowBookRequest struct {
	BookId uint64 `json:"book_id" binding:"required"`
	BorrowDays *int `json:"borrow_days"`
}

type ReturnBookRequest struct {
	Condition *string `json:"condition" binding:"omitempty,oneof=good damaged lost"`
	Remark    *string `json:"remark"    binding:"omitempty,max=255"`
}