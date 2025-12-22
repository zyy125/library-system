package request

type GetBorrowStatsRequest struct {
	StartDate string `form:"start_date"` // 格式：YYYY-MM-DD
	EndDate   string `form:"end_date"`   // 格式：YYYY-MM-DD
	GroupBy   string `form:"group_by"`   // day/week/month，默认day
}

type GetPopularBooksRequest struct {
	Limit      int  `form:"limit"`       // 默认10
	Period     string `form:"period"`    // 7d/30d/90d/all，默认30d
	CategoryID *uint  `form:"category_id"`
}

type GetUserStatsRequest struct {
	UserID uint64 `uri:"user_id" binding:"required"`
}