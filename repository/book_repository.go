package repository

import (
	"context"
	"library-system/dto/request"
	"library-system/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) DB() *gorm.DB {
	return r.db
}

func (r *BookRepository) GetBookByISBN(ctx context.Context, ISBN string) (model.Book, error) {
	return gorm.G[model.Book](r.db).Where("isbn = ?", ISBN).First(ctx)
}

func (r *BookRepository) GetBookByID(ctx context.Context, id uint64) (model.Book, error) {
	return gorm.G[model.Book](r.db).Where("id = ?", id).First(ctx)
}

func (r *BookRepository) CreateBook(ctx context.Context, book *model.Book) error {
	return gorm.G[model.Book](r.db).Create(ctx, book)
}

func (r *BookRepository) GetBookList(ctx context.Context, req *request.GetBookListRequest) ([]model.Book, int64, error) {
	db := r.db.WithContext(ctx).Model(&model.Book{}).Preload("Category")

	if req.Title != nil {
		db = db.Where("title LIKE ?", "%"+*req.Title+"%")
	}
	if req.Author != nil {
		db = db.Where("author LIKE ?", "%"+*req.Author+"%")
	}
	if req.ISBN != nil {
		db = db.Where("isbn = ?", req.ISBN)
	}
	if req.Publisher != nil {
		db = db.Where("publisher = ?", req.Publisher)
	}
	if req.CategoryID != nil {
		db = db.Where("category_id = ?", req.CategoryID)
	}
	if req.AvailableOnly != nil {
		if *req.AvailableOnly {
			db = db.Where("stock - borrow_count > 0")
		} else {
			db = db.Where("stock - borrow_count <= 0")
		}
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortBy := "created_at"
	order := "desc"

	if req.SortBy != nil {
		sortBy = *req.SortBy
	}
	if req.Order != nil {
		order = *req.Order
	}

	db = db.Order(sortBy + " " + order)

	page := req.Page
	limit := req.Limit

	offset := (page - 1) * limit

	var books []model.Book
	if err := db.Offset(offset).Limit(limit).Find(&books).Error; err != nil {
		return nil, 0, err
	}

	return books, total, nil
}

func (r *BookRepository) UpdateBookFields(ctx context.Context, id uint64, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(model.Book{}).Where("id = ?", id).Updates(fields).Error
}

func (r *BookRepository) DeleteBookByID(ctx context.Context, id uint64) error {
	_, err := gorm.G[model.Book](r.db).Where("id = ?", id).Delete(ctx)
	return err
}

func (r *BookRepository) GetBookByIDWithLock(ctx context.Context, tx *gorm.DB, id uint64) (model.Book, error) {
	txLock := tx.Clauses(clause.Locking{Strength: "UPDATE"})
	
	return gorm.G[model.Book](txLock).Where("id = ?", id).First(ctx)
}

func (r *BookRepository) IncreaseBorrowCount(ctx context.Context, tx *gorm.DB, id uint64, count int) error {
	err := tx.Model(&model.Book{}).Where("id = ?", id).
    UpdateColumn("borrow_count", gorm.Expr("borrow_count + ?", 1)).Error
	return err
}

func (r *BookRepository) DecreaseBorrowCount(ctx context.Context, tx *gorm.DB, id uint64, count int) error {
	err := tx.Model(&model.Book{}).Where("id = ?", id).
    UpdateColumn("borrow_count", gorm.Expr("borrow_count - ?", 1)).Error
	return err
}