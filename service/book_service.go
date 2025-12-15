package service

import (
	"context"
	"library-system/common"
	"library-system/dto/request"
	"library-system/dto/response"
	"library-system/model"
	"library-system/repository"
	"time"
	"errors"
	"gorm.io/gorm"
)

type BookService struct {
	bookRepo *repository.BookRepository
	categoryRepo *repository.CategoryRepository
}

func NewBookService(bookRepo *repository.BookRepository, categoryRepo *repository.CategoryRepository) *BookService {
	return &BookService{bookRepo: bookRepo,categoryRepo: categoryRepo}
}

func (s *BookService) CreateBook(ctx context.Context, req *request.CreateBookRequest) (*response.CreateBookResponse, error) { 
	if _, err := s.bookRepo.GetBookByISBN(ctx, req.ISBN); err == nil {
		return nil, common.ErrISBNExist
	} else if !errors.Is(err, gorm.ErrRecordNotFound){
		return nil, err
	}

	category, err := s.categoryRepo.GetCategoryByID(ctx, req.CategoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrCategoryNotFound
		}
		return nil, err		
	}

    var publishTime *time.Time
    if req.PublishDate != nil {
        t, err := time.Parse("2006-01-02", *req.PublishDate)
        if err != nil {
            return nil, err
        }
        publishTime = &t
    }

    book := model.Book{
        Title:       req.Title,
        Author:      req.Author,
        ISBN:        req.ISBN,
        CategoryID:  req.CategoryID,
        Publisher:   req.Publisher,
        PublishDate: publishTime,
        Price:       *req.Price,
        Stock:       req.Stock,
        Description: *req.Description,
        CoverURL:    *req.CoverURL,
    }

	if err := s.bookRepo.CreateBook(ctx, &book); err != nil {
		return nil, err
	}

	resp := &response.CreateBookResponse{
        ID:           book.ID,
        Title:        book.Title,
        Author:       book.Author,
        ISBN:         book.ISBN,
        CategoryID:   book.CategoryID,
        CategoryName: category.Name,
        Publisher:    book.Publisher,
        PublishDate:  req.PublishDate,
        Price:        &book.Price,
        Stock:        book.Stock,
        Available:    book.Stock,
        BorrowCount:  0,
        CreatedAt:    time.Now().UTC().Format(time.RFC3339),
    }

	return resp, nil
}

