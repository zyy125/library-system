package request

type BorrowBookRequest struct {
	BookId     uint64 `json:"book_id" binding:"required"`
	BorrowDays *int   `json:"borrow_days"`
}

type ReturnBookRequest struct {
	Condition *string `json:"condition" binding:"omitempty,oneof=good damaged lost"`
	Remark    *string `json:"remark"    binding:"omitempty,max=255"`
}

type RenewBorrowRequest struct {
	RenewDays *int `json:"renew_days" binding:"omitempty,min=1,max=90"`
}

type GetBorrowRecordListRequest struct {
	UserID    *uint64 `form:"user_id"`
	Status    *string `form:"status" binding:"omitempty,oneof=borrowed returned overdue"`
	BookID    *uint64 `form:"book_id"`
	StartDate *string `form:"start_date"`
	EndDate   *string `form:"end_date"`
	Page      int     `form:"page"  binding:"omitempty,min=1"`
	Limit     int     `form:"limit" binding:"omitempty,min=1,max=100"`
	SortBy    *string `form:"sort_by" binding:"omitempty,oneof=borrow_date due_date return_date"`
	Order     *string `form:"order" binding:"omitempty,oneof=asc desc"`
}
