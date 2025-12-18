package repository

import (
	"context"
	"library-system/dto/request"
	"library-system/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) DB() *gorm.DB {
	return r.db
}

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	return gorm.G[model.User](r.db).Create(ctx, user)
}

func (r *UserRepository) GetUserByUserID(ctx context.Context, id uint64) (model.User, error) {
	return gorm.G[model.User](r.db).Where("id = ?", id).First(ctx)
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	return gorm.G[model.User](r.db).Where("username = ?", username).First(ctx)
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	return gorm.G[model.User](r.db).Where("email = ?", email).First(ctx)
}

func (r *UserRepository) GetUsersByRole(ctx context.Context, role string) ([]model.User, error) {
	return gorm.G[model.User](r.db).Where("role = ?", role).Find(ctx)
}

func (r *UserRepository) GetUsersByStatus(ctx context.Context, status string) ([]model.User, error) {
	return gorm.G[model.User](r.db).Where("status = ?", status).Find(ctx)
}

func (r *UserRepository) UpdateUserFields(ctx context.Context, db *gorm.DB, id uint64, fields map[string]interface{}) error {
	return db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(fields).Error
}

func (r *UserRepository) GetUserList(ctx context.Context, query *request.GetUserListRequest) ([]model.User, int64, error) {
	var users []model.User
	db := r.db.WithContext(ctx).Model(&model.User{})

	if query.Username != "" {
		db = db.Where("username LIKE ?", "%"+query.Username+"%")
	}
	if query.Role != "" {
		db = db.Where("role = ?", query.Role)
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (query.Page - 1) * query.Limit
	if err := db.
		Offset(offset).
		Limit(query.Limit).
		Order("id DESC").
		Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *UserRepository) DeleteUserByID(ctx context.Context, id uint64) error {
	_, err := gorm.G[model.User](r.db).Where("id = ?", id).Delete(ctx)
	return err
}

func (r *UserRepository) GetUserByIDWithLock(ctx context.Context, tx *gorm.DB, id uint64) (model.User, error) {
	txLock := tx.Clauses(clause.Locking{Strength: "UPDATE"})
	
	return gorm.G[model.User](txLock).Where("id = ?", id).First(ctx)
}

func (r *UserRepository) IncreaseBorrowingCount(ctx context.Context, tx *gorm.DB, id uint64, count int) error {
	err := tx.Model(&model.User{}).Where("id = ?", id).
    UpdateColumn("borrowing_count", gorm.Expr("borrowing_count + ?", 1)).Error

	return err
}

func (r *UserRepository) DecreaseBorrowingCount(ctx context.Context, tx *gorm.DB, id uint64, count int) error {
	err := tx.Model(&model.User{}).Where("id = ?", id).
    UpdateColumn("borrowing_count", gorm.Expr("borrowing_count - ?", 1)).Error

	return err
}

func (r *UserRepository) IncreaseOverDueCount(ctx context.Context, tx *gorm.DB, id uint64, count int) error {
	err := tx.Model(&model.User{}).Where("id = ?", id).
    UpdateColumn("overdue_count", gorm.Expr("overdue_count + ?", count)).Error

	return err
}

func (r *UserRepository) DecreaseOverDueCount(ctx context.Context, tx *gorm.DB, id uint64, count int) error {
	err := tx.Model(&model.User{}).Where("id = ?", id).
    UpdateColumn("overdue_count", gorm.Expr("overdue_count - ?", count)).Error

	return err
}