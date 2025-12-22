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

func (r *CategoryRepository) DB() *gorm.DB {
	return r.db
}

func (r *CategoryRepository) GetCategoryByID(ctx context.Context, id uint) (model.Category, error) {
	var category model.Category
	err := r. db.WithContext(ctx).Where("id = ?", id).First(&category).Error
	return category, err
}

func (r *CategoryRepository) GetCategoryByName(ctx context. Context, name string) (model.Category, error) {
	var category model.Category
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&category).Error
	return category, err
}

func (r *CategoryRepository) GetCategoryList(ctx context.Context) ([]model.Category, error) {
	var categories []model.Category
	err := r.db. WithContext(ctx).Order("id ASC").Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) GetCategoryListWithBookCount(ctx context. Context) ([]model.Category, []int64, error) {
	var categories []model.Category
	err := r.db. WithContext(ctx).Order("id ASC").Find(&categories).Error
	if err != nil {
		return nil, nil, err
	}

	counts := make([]int64, len(categories))
	for i, cat := range categories {
		var count int64
		r.db.WithContext(ctx).Model(&model.Book{}).Where("category_id = ?", cat.ID).Count(&count)
		counts[i] = count
	}

	return categories, counts, nil
}

func (r *CategoryRepository) GetChildCategories(ctx context.Context, parentID uint) ([]model.Category, error) {
	var categories []model.Category
	err := r. db.WithContext(ctx).Where("parent_id = ?", parentID).Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) GetBookCountByCategoryID(ctx context.Context, categoryID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Book{}).Where("category_id = ?", categoryID).Count(&count).Error
	return count, err
}

func (r *CategoryRepository) CreateCategory(ctx context. Context, category *model.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *CategoryRepository) UpdateCategory(ctx context. Context, id uint, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.Category{}).Where("id = ?", id).Updates(updates).Error
}

func (r *CategoryRepository) DeleteCategory(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Category{}, id).Error
}

func (r *CategoryRepository) HasBooks(ctx context.Context, categoryID uint) (bool, int64, error) {
	var count int64
	err := r.db. WithContext(ctx).Model(&model.Book{}).Where("category_id = ?", categoryID).Count(&count).Error
	return count > 0, count, err
}

func (r *CategoryRepository) HasChildren(ctx context.Context, categoryID uint) (bool, error) {
	var count int64
	err := r.db. WithContext(ctx).Model(&model.Category{}).Where("parent_id = ?", categoryID).Count(&count).Error
	return count > 0, err
}