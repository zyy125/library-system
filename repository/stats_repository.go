package repository

import (
	"context"
	"library-system/model"
	"time"

	"gorm.io/gorm"
)

type StatsRepository struct {
	db *gorm.DB
}

func NewStatsRepository(db *gorm.DB) *StatsRepository {
	return &StatsRepository{db: db}
}

func (r *StatsRepository) DB() *gorm.DB {
	return r.db
}

// ========== 系统概览统计 ==========

func (r *StatsRepository) CountTotalBooks(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Book{}).Count(&count).Error
	return count, err
}

func (r *StatsRepository) CountTotalUsers(ctx context.Context) (int64, error) {
	var count int64
	err := r. db.WithContext(ctx).Model(&model.User{}).Count(&count).Error
	return count, err
}

func (r *StatsRepository) CountTotalCategories(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Category{}).Count(&count).Error
	return count, err
}

func (r *StatsRepository) CountBorrowedBooks(ctx context. Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.BorrowRecord{}).
		Where("status IN ? ", []string{"borrowed", "overdue"}).
		Count(&count).Error
	return count, err
}

func (r *StatsRepository) CountAvailableBooks(ctx context.Context) (int64, error) {
	var total int64
	err := r. db.WithContext(ctx).Model(&model.Book{}).
		Select("COALESCE(SUM(stock - borrow_count), 0)").
		Scan(&total).Error
	return total, err
}

func (r *StatsRepository) CountOverdueBooks(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.BorrowRecord{}).
		Where("status = ?", "overdue").
		Count(&count).Error
	return count, err
}

func (r *StatsRepository) CountTotalBorrowRecords(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.BorrowRecord{}).Count(&count).Error
	return count, err
}

func (r *StatsRepository) CountActiveReservations(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Reservation{}).
		Where("status IN ?", []string{"waiting", "available"}).
		Count(&count).Error
	return count, err
}

func (r *StatsRepository) CountActiveUsers30d(ctx context.Context) (int64, error) {
	var count int64
	thirtyDaysAgo := time. Now().AddDate(0, 0, -30)
	err := r.db.WithContext(ctx).Model(&model.BorrowRecord{}).
		Where("borrow_date >= ?", thirtyDaysAgo).
		Distinct("user_id").
		Count(&count).Error
	return count, err
}

// ========== 借阅统计 ==========

type DailyBorrowStats struct {
	Date        string
	BorrowCount int64
	ReturnCount int64
}

func (r *StatsRepository) GetBorrowStatsByDateRange(ctx context.Context, startDate, endDate time.Time, groupBy string) ([]DailyBorrowStats, error) {
	var results []DailyBorrowStats

	dateFormat := "%Y-%m-%d"
	if groupBy == "week" {
		dateFormat = "%Y-%u" // 年-周
	} else if groupBy == "month" {
		dateFormat = "%Y-%m"
	}

	// 借阅统计
	var borrowStats []struct {
		Date  string
		Count int64
	}
	err := r.db.WithContext(ctx).Model(&model.BorrowRecord{}).
		Select("DATE_FORMAT(borrow_date, ?) as date, COUNT(*) as count", dateFormat).
		Where("borrow_date BETWEEN ? AND ?", startDate, endDate).
		Group("date").
		Order("date ASC").
		Scan(&borrowStats).Error
	if err != nil {
		return nil, err
	}

	// 归还统计
	var returnStats []struct {
		Date  string
		Count int64
	}
	err = r.db. WithContext(ctx).Model(&model.BorrowRecord{}).
		Select("DATE_FORMAT(return_date, ?) as date, COUNT(*) as count", dateFormat).
		Where("return_date BETWEEN ? AND ?", startDate, endDate).
		Where("return_date IS NOT NULL").
		Group("date").
		Order("date ASC").
		Scan(&returnStats).Error
	if err != nil {
		return nil, err
	}

	// 合并结果
	borrowMap := make(map[string]int64)
	returnMap := make(map[string]int64)
	dateSet := make(map[string]bool)

	for _, b := range borrowStats {
		borrowMap[b.Date] = b.Count
		dateSet[b.Date] = true
	}
	for _, r := range returnStats {
		returnMap[r. Date] = r.Count
		dateSet[r.Date] = true
	}

	for date := range dateSet {
		results = append(results, DailyBorrowStats{
			Date:        date,
			BorrowCount: borrowMap[date],
			ReturnCount: returnMap[date],
		})
	}

	return results, nil
}

func (r *StatsRepository) CountBorrowsByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.BorrowRecord{}).
		Where("borrow_date BETWEEN ?  AND ?", startDate, endDate).
		Count(&count).Error
	return count, err
}

func (r *StatsRepository) CountReturnsByDateRange(ctx context. Context, startDate, endDate time.Time) (int64, error) {
	var count int64
	err := r.db. WithContext(ctx).Model(&model.BorrowRecord{}).
		Where("return_date BETWEEN ? AND ?", startDate, endDate).
		Where("return_date IS NOT NULL").
		Count(&count).Error
	return count, err
}

// ========== 用户统计 ==========

func (r *StatsRepository) GetUserTotalBorrowCount(ctx context.Context, userID uint64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.BorrowRecord{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}

func (r *StatsRepository) GetUserTotalFine(ctx context.Context, userID uint64) (float64, error) {
	var total float64
	err := r. db.WithContext(ctx).Model(&model.BorrowRecord{}).
		Select("COALESCE(SUM(fine), 0)").
		Where("user_id = ?", userID).
		Scan(&total).Error
	return total, err
}

type CategoryBorrowCount struct {
	CategoryID   uint
	CategoryName string
	BorrowCount  int64
}

func (r *StatsRepository) GetUserFavoriteCategories(ctx context.Context, userID uint64, limit int) ([]CategoryBorrowCount, error) {
	var results []CategoryBorrowCount
	err := r.db.WithContext(ctx).
		Table("borrow_records br").
		Select("b.category_id, c.name as category_name, COUNT(*) as borrow_count").
		Joins("JOIN books b ON br.book_id = b.id").
		Joins("JOIN categories c ON b.category_id = c.id").
		Where("br.user_id = ?", userID).
		Group("b.category_id, c.name").
		Order("borrow_count DESC").
		Limit(limit).
		Scan(&results).Error
	return results, err
}

type MonthlyBorrowCount struct {
	Month string
	Count int64
}

func (r *StatsRepository) GetUserReadingTrend(ctx context.Context, userID uint64, months int) ([]MonthlyBorrowCount, error) {
	var results []MonthlyBorrowCount
	startDate := time.Now().AddDate(0, -months, 0)
	err := r.db.WithContext(ctx).Model(&model.BorrowRecord{}).
		Select("DATE_FORMAT(borrow_date, '%Y-%m') as month, COUNT(*) as count").
		Where("user_id = ? AND borrow_date >= ?", userID, startDate).
		Group("month").
		Order("month ASC").
		Scan(&results).Error
	return results, err
}

// ========== 热门图书 ==========

type PopularBook struct {
	BookID      uint64
	Title       string
	Author      string
	CoverURL    string
	BorrowCount int64
	Rating      float64
}

func (r *StatsRepository) GetPopularBooks(ctx context.Context, limit int, startDate *time.Time, categoryID *uint) ([]PopularBook, error) {
	var results []PopularBook

	query := r.db.WithContext(ctx).
		Table("borrow_records br").
		Select("b.id as book_id, b.title, b.author, b.cover_url, COUNT(*) as borrow_count, b.rating").
		Joins("JOIN books b ON br.book_id = b.id")

	if startDate != nil {
		query = query.Where("br.borrow_date >= ?", startDate)
	}

	if categoryID != nil {
		query = query.Where("b.category_id = ?", *categoryID)
	}

	err := query.
		Group("b.id, b.title, b.author, b.cover_url, b.rating").
		Order("borrow_count DESC").
		Limit(limit).
		Scan(&results).Error

	return results, err
}

// ========== 分类统计 ==========

type CategoryStats struct {
	CategoryID   uint
	CategoryName string
	BookCount    int64
	BorrowCount  int64
}

func (r *StatsRepository) GetCategoryStats(ctx context.Context) ([]CategoryStats, error) {
	var results []CategoryStats

	err := r.db.WithContext(ctx).
		Table("categories c").
		Select(`
			c.id as category_id, 
			c.name as category_name, 
			COUNT(DISTINCT b.id) as book_count,
			COALESCE(SUM(b.borrow_count), 0) as borrow_count
		`).
		Joins("LEFT JOIN books b ON c.id = b.category_id").
		Group("c.id, c. name").
		Order("borrow_count DESC").
		Scan(&results).Error

	return results, err
}