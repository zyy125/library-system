package repository

import (
	"context"
	"library-system/model"

	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) GetBookByISBN(ctx context.Context, ISBN string) (model.Book, error) {
	return gorm.G[model.Book](r.db).Where("isbn = ?", ISBN).First(ctx)
}

func (r *BookRepository) CreateBook(ctx context.Context, book *model.Book) error {
	return gorm.G[model.Book](r.db).Create(ctx, book)
}