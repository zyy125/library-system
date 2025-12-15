package repository

import (
	"context"
	"library-system/model"
	"library-system/dto/request"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
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

func (r *UserRepository) UpdateUserFields(ctx context.Context, id uint64, fields map[string]interface{}) error {
	    return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(fields).Error
}

func (r *UserRepository) List(ctx context.Context, query request.GetUserListRequest) ([]model.User, int64, error) {
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