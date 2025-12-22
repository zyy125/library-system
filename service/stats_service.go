package service

import (
	"context"
	"library-system/common"
	"library-system/dto/request"
	"library-system/dto/response"
	"library-system/repository"
	"sort"
	"time"
)

type StatsService struct {
	statsRepo *repository.StatsRepository
	userRepo  *repository.UserRepository
}

func NewStatsService(statsRepo *repository.StatsRepository, userRepo *repository.UserRepository) *StatsService {
	return &StatsService{
		statsRepo: statsRepo,
		userRepo:  userRepo,
	}
}

// GetOverview 获取系统统计概览
func (s *StatsService) GetOverview(ctx context.Context) (*response.StatsOverviewResponse, error) {
	totalBooks, _ := s.statsRepo.CountTotalBooks(ctx)
	totalUsers, _ := s.statsRepo.CountTotalUsers(ctx)
	totalCategories, _ := s.statsRepo.CountTotalCategories(ctx)
	borrowedBooks, _ := s.statsRepo.CountBorrowedBooks(ctx)
	availableBooks, _ := s.statsRepo. CountAvailableBooks(ctx)
	overdueBooks, _ := s.statsRepo.CountOverdueBooks(ctx)
	totalBorrowCount, _ := s.statsRepo. CountTotalBorrowRecords(ctx)
	reservationsCount, _ := s.statsRepo.CountActiveReservations(ctx)
	activeUsers30d, _ := s.statsRepo.CountActiveUsers30d(ctx)

	return &response.StatsOverviewResponse{
		TotalBooks:        totalBooks,
		TotalUsers:        totalUsers,
		TotalCategories:   totalCategories,
		BorrowedBooks:     borrowedBooks,
		AvailableBooks:    availableBooks,
		OverdueBooks:      overdueBooks,
		TotalBorrowCount:  totalBorrowCount,
		ReservationsCount: reservationsCount,
		ActiveUsers30d:    activeUsers30d,
	}, nil
}

// GetBorrowStats 获取借阅统计
func (s *StatsService) GetBorrowStats(ctx context.Context, req *request. GetBorrowStatsRequest) (*response.GetBorrowStatsResponse, error) {
	// 默认值处理
	groupBy := req.GroupBy
	if groupBy == "" {
		groupBy = "day"
	}

	// 解析日期
	var startDate, endDate time.Time
	var err error

	if req.StartDate != "" {
		startDate, err = time. Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, common.ErrBadRequest
		}
	} else {
		startDate = time.Now().AddDate(0, 0, -30) // 默认30天前
	}

	if req.EndDate != "" {
		endDate, err = time.Parse("2006-01-02", req. EndDate)
		if err != nil {
			return nil, common.ErrBadRequest
		}
	} else {
		endDate = time.Now()
	}

	// 设置结束日期为当天23:59:59
	endDate = endDate.Add(24*time.Hour - time.Second)

	// 获取统计数据
	chartData, err := s.statsRepo.GetBorrowStatsByDateRange(ctx, startDate, endDate, groupBy)
	if err != nil {
		return nil, err
	}

	totalBorrows, _ := s. statsRepo.CountBorrowsByDateRange(ctx, startDate, endDate)
	totalReturns, _ := s.statsRepo.CountReturnsByDateRange(ctx, startDate, endDate)

	// 转换为响应格式
	chartItems := make([]response.BorrowStatsChartItem, len(chartData))
	for i, item := range chartData {
		chartItems[i] = response.BorrowStatsChartItem{
			Date:        item.Date,
			BorrowCount: item.BorrowCount,
			ReturnCount: item.ReturnCount,
		}
	}

	// 按日期排序
	sort. Slice(chartItems, func(i, j int) bool {
		return chartItems[i].Date < chartItems[j].Date
	})

	return &response.GetBorrowStatsResponse{
		Period:  response.BorrowStatsPeriod{
			Start: startDate.Format("2006-01-02"),
			End:   endDate.Format("2006-01-02"),
		},
		TotalBorrows: totalBorrows,
		TotalReturns: totalReturns,
		ChartData:    chartItems,
	}, nil
}

// GetUserStats 获取用户借阅统计
func (s *StatsService) GetUserStats(ctx context.Context, userID uint64) (*response.GetUserStatsResponse, error) {
	// 获取用户信息
	user, err := s.userRepo.GetUserByUserID(ctx, userID)
	if err != nil {
		return nil, common.ErrNotFound
	}

	totalBorrowCount, _ := s.statsRepo.GetUserTotalBorrowCount(ctx, userID)
	totalFine, _ := s.statsRepo.GetUserTotalFine(ctx, userID)
	favoriteCategories, _ := s. statsRepo.GetUserFavoriteCategories(ctx, userID, 5)
	readingTrend, _ := s.statsRepo. GetUserReadingTrend(ctx, userID, 6)

	// 转换喜爱分类
	favCategories := make([]response.FavoriteCategory, len(favoriteCategories))
	for i, fc := range favoriteCategories {
		favCategories[i] = response.FavoriteCategory{
			CategoryID:   fc.CategoryID,
			CategoryName: fc.CategoryName,
			BorrowCount:  fc.BorrowCount,
		}
	}

	// 转换阅读趋势
	trend := make([]response.ReadingTrend, len(readingTrend))
	for i, rt := range readingTrend {
		trend[i] = response.ReadingTrend{
			Month: rt.Month,
			Count: rt.Count,
		}
	}

	return &response.GetUserStatsResponse{
		UserID:             userID,
		Username:           user.Username,
		TotalBorrowCount:   totalBorrowCount,
		CurrentBorrowing:   user.BorrowingCount,
		OverdueCount:       user.OverdueCount,
		TotalFine:          totalFine,
		FavoriteCategories: favCategories,
		ReadingTrend:       trend,
	}, nil
}

// GetPopularBooks 获取热门图书排行
func (s *StatsService) GetPopularBooks(ctx context.Context, req *request. GetPopularBooksRequest) (*response.GetPopularBooksResponse, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	period := req.Period
	if period == "" {
		period = "30d"
	}

	// 计算起始日期
	var startDate *time.Time
	switch period {
	case "7d":
		t := time.Now().AddDate(0, 0, -7)
		startDate = &t
	case "30d":
		t := time.Now().AddDate(0, 0, -30)
		startDate = &t
	case "90d": 
		t := time.Now().AddDate(0, 0, -90)
		startDate = &t
	case "all":
		startDate = nil
	default:
		t := time.Now().AddDate(0, 0, -30)
		startDate = &t
	}

	books, err := s.statsRepo.GetPopularBooks(ctx, limit, startDate, req.CategoryID)
	if err != nil {
		return nil, err
	}

	items := make([]response.PopularBookItem, len(books))
	for i, book := range books {
		items[i] = response.PopularBookItem{
			Rank:        i + 1,
			BookID:      book.BookID,
			Title:       book.Title,
			Author:      book. Author,
			CoverURL:     book.CoverURL,
			BorrowCount: book. BorrowCount,
			Rating:      book.Rating,
		}
	}

	return &response.GetPopularBooksResponse{
		Period:  period,
		Books:  items,
	}, nil
}

// GetCategoryStats 获取分类统计
func (s *StatsService) GetCategoryStats(ctx context. Context) (*response.GetCategoryStatsResponse, error) {
	stats, err := s.statsRepo. GetCategoryStats(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]response.CategoryStatsItem, len(stats))
	for i, stat := range stats {
		var borrowRate float64
		if stat.BookCount > 0 {
			borrowRate = float64(stat.BorrowCount) / float64(stat.BookCount)
		}
		items[i] = response.CategoryStatsItem{
			CategoryID:   stat.CategoryID,
			CategoryName: stat.CategoryName,
			BookCount:    stat.BookCount,
			BorrowCount:  stat.BorrowCount,
			BorrowRate:    borrowRate,
		}
	}

	return &response. GetCategoryStatsResponse{
		Categories: items,
	}, nil
}