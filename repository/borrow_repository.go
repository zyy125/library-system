package repository

import (
	"context"
	"library-system/model"
	"library-system/dto/request"
	"time"

	"gorm.io/gorm"
)

type BorrowRepository struct {
	db *gorm.DB
}

func NewBorrowRepository(db *gorm.DB) *BorrowRepository {
	return &BorrowRepository{db: db}
}

func (r *BorrowRepository) DB() *gorm.DB {
	return r.db
} 

func (r *BorrowRepository) CreateBorrowRecord(ctx context.Context, tx *gorm.DB, borrow *model.BorrowRecord) error {
	return gorm.G[model.BorrowRecord](tx).Create(ctx, borrow)
}

func (r *BorrowRepository) GetBorrowRecordByID(ctx context.Context, id uint64) (model.BorrowRecord, error) {
	return gorm.G[model.BorrowRecord](r.db).Where("id = ?", id).First(ctx)
}

func (r *BorrowRepository) GetAllDueRecord(ctx context.Context, tx *gorm.DB, now time.Time) ([]model.BorrowRecord, error) {
	return gorm.G[model.BorrowRecord](tx).Where("status = ? AND return_date IS NULL AND due_date < ?","borrowed", now).Find(ctx)
}

func (r *BorrowRepository) GetDueRecordByUserID(ctx context.Context, tx *gorm.DB, id uint64, now time.Time) ([]model.BorrowRecord, error) {
	return gorm.G[model.BorrowRecord](tx).Where("user_id = ? AND status = ? AND return_date IS NULL AND due_date < ?", id, "borrowed", now).Find(ctx)
}

func (r *BorrowRepository) UpdateFields(ctx context.Context, tx *gorm.DB, id uint64, fields map[string]interface{}) error {
	return tx.WithContext(ctx).Model(&model.BorrowRecord{}).Where("id = ?", id).Updates(fields).Error
}

func (r *BorrowRepository) GetBorrowRecordByIDWithBook(ctx context.Context, id uint64) (model.BorrowRecord, error) {
	return gorm.G[model.BorrowRecord](r.db).Where("id = ?", id).Preload("Book", func(db gorm.PreloadBuilder) error {return nil}).First(ctx)
}
func (r *BorrowRepository) GetBorrowRecordList(ctx context.Context, req *request.GetBorrowRecordListRequest) ([]model.BorrowRecord, int64, error) {
	db := r.db.WithContext(ctx).Model(model.BorrowRecord{}).Preload("Book").Preload("User")

	if req.UserID != nil {
		db = db.Where("user_id = ?", *req.UserID)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.BookID != nil {
		db = db.Where("book_id = ?", *req.BookID)
	}
	if req.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			return nil, 0, err
		}
		db = db.Where("borrow_date >= ?", startDate)
	}
	if req.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return nil, 0, err
		}
		// borrow_date <= end_date (到当天 23:59:59 结束)
		endDate = endDate.Add(24*time.Hour - time.Second) // 2025-12-05 23:59:59
		db = db.Where("borrow_date <= ?", endDate)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortBy := "borrow_date"
	if req.SortBy != nil {
		sortBy = *req.SortBy
	}

	order := "desc"
	if req.Order != nil {
		order = *req.Order
	}

	db = db.Order(sortBy + " " + order)

	page := req.Page
	limit := req.Limit

	offset := (page - 1) * limit
	db = db.Offset(offset).Limit(limit)

	var records []model.BorrowRecord
	if err := db.Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}