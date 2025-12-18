package repository

import (
	"context"
	"library-system/model"
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

func (r *BorrowRepository) FindAllDueRecord(ctx context.Context, tx *gorm.DB, now time.Time) ([]model.BorrowRecord, error) {
	return gorm.G[model.BorrowRecord](tx).Where("status = ? AND return_date IS NULL AND due_date < ?","borrowed", now).Find(ctx)
}

func (r *BorrowRepository) FindDueRecordByUserID(ctx context.Context, tx *gorm.DB, id uint64, now time.Time) ([]model.BorrowRecord, error) {
	return gorm.G[model.BorrowRecord](tx).Where("user_id = ? AND status = ? AND return_date IS NULL AND due_date < ?", id, "borrowed", now).Find(ctx)
}

func (r *BorrowRepository) UpdateFields(ctx context.Context, tx *gorm.DB, id uint64, fields map[string]interface{}) error {
	return tx.WithContext(ctx).Model(&model.BorrowRecord{}).Where("id = ?", id).Updates(fields).Error
}