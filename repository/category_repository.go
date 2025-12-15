package repository

import (
	"context"
	"library-system/model"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetCategoryByID(ctx context.Context, id uint) (model.Category, error) {
	return gorm.G[model.Category](r.db).Where("id = ?", id).First(ctx)
}