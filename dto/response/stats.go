package response

// ========== 7.1 系统统计概览 ==========

type StatsOverviewResponse struct {
	TotalBooks        int64 `json:"total_books"`
	TotalUsers        int64 `json:"total_users"`
	TotalCategories   int64 `json:"total_categories"`
	BorrowedBooks     int64 `json:"borrowed_books"`
	AvailableBooks    int64 `json:"available_books"`
	OverdueBooks      int64 `json:"overdue_books"`
	TotalBorrowCount  int64 `json:"total_borrow_count"`
	ReservationsCount int64 `json:"reservations_count"`
	ActiveUsers30d    int64 `json:"active_users_30d"`
}

// ========== 7.2 借阅统计 ==========

type BorrowStatsPeriod struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type BorrowStatsChartItem struct {
	Date        string `json:"date"`
	BorrowCount int64  `json:"borrow_count"`
	ReturnCount int64  `json:"return_count"`
}

type GetBorrowStatsResponse struct {
	Period       BorrowStatsPeriod      `json:"period"`
	TotalBorrows int64                  `json:"total_borrows"`
	TotalReturns int64                  `json:"total_returns"`
	ChartData    []BorrowStatsChartItem `json:"chart_data"`
}

// ========== 7.3 用户借阅统计 ==========

type FavoriteCategory struct {
	CategoryID   uint   `json:"category_id"`
	CategoryName string `json:"category_name"`
	BorrowCount  int64  `json:"borrow_count"`
}

type ReadingTrend struct {
	Month string `json:"month"`
	Count int64  `json:"count"`
}

type GetUserStatsResponse struct {
	UserID             uint64             `json:"user_id"`
	Username           string             `json:"username"`
	TotalBorrowCount   int64              `json:"total_borrow_count"`
	CurrentBorrowing   int                `json:"current_borrowing"`
	OverdueCount       int                `json:"overdue_count"`
	TotalFine          float64            `json:"total_fine"`
	FavoriteCategories []FavoriteCategory `json:"favorite_categories"`
	ReadingTrend       []ReadingTrend     `json:"reading_trend"`
}

// ========== 7.4 热门图书排行 ==========

type PopularBookItem struct {
	Rank        int     `json:"rank"`
	BookID      uint64  `json:"book_id"`
	Title       string  `json:"title"`
	Author      string  `json:"author"`
	CoverURL    string  `json:"cover_url"`
	BorrowCount int64   `json:"borrow_count"`
	Rating      float64 `json:"rating"`
}

type GetPopularBooksResponse struct {
	Period string            `json:"period"`
	Books  []PopularBookItem `json:"books"`
}

// ========== 7.5 分类统计 ==========

type CategoryStatsItem struct {
	CategoryID   uint    `json:"category_id"`
	CategoryName string  `json:"category_name"`
	BookCount    int64   `json:"book_count"`
	BorrowCount  int64   `json:"borrow_count"`
	BorrowRate   float64 `json:"borrow_rate"`
}

type GetCategoryStatsResponse struct {
	Categories []CategoryStatsItem `json:"categories"`
}