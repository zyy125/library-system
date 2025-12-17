package service

import (
	"context"
	"errors"
	"library-system/common"
	"library-system/dto/request"
	"library-system/dto/response"
	"library-system/model"
	"library-system/repository"
	"math"
	"time"

	"gorm.io/gorm"
)

type BookService struct {
	bookRepo     *repository.BookRepository
	categoryRepo *repository.CategoryRepository
}

func NewBookService(bookRepo *repository.BookRepository, categoryRepo *repository.CategoryRepository) *BookService {
	return &BookService{bookRepo: bookRepo, categoryRepo: categoryRepo}
}

func (s *BookService) CreateBook(ctx context.Context, req *request.CreateBookRequest) (*response.CreateBookResponse, error) {
	if _, err := s.bookRepo.GetBookByISBN(ctx, req.ISBN); err == nil {
		return nil, common.ErrISBNExist
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
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
		Stock:       req.Stock,
	}

	if req.Price != nil {
		book.Price = *req.Price
	}
	if req.Description != nil {
		book.Description = *req.Description
	}
	if req.CoverURL != nil {
		book.CoverURL = *req.CoverURL
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
		Price:        req.Price,
		Stock:        book.Stock,
		Available:    book.Stock,
		CoverUrl:     req.CoverURL,
		BorrowCount:  0,
		CreatedAt:    time.Now().UTC().Format(time.RFC3339),
	}

	return resp, nil
}

func (s *BookService) BatchCreateBook(ctx context.Context, req *request.BatchCreateBookRequest) (*response.BatchCreateBookResponse, error) {
	resp := &response.BatchCreateBookResponse{
		FailedItems: make([]response.FailedItems, 0),
	}

	for i, book := range req.Books {
		_, err := s.CreateBook(ctx, &book)
		if err != nil {
			fail := response.FailedItems{
				Index: i,
				ISBN:  book.ISBN,
			}
			switch err {
			case common.ErrISBNExist:
				fail.Error = "ISBN已存在"
			case common.ErrCategoryNotFound:
				fail.Error = "分类不存在"
			default:
				fail.Error = "导入失败"
			}
			resp.FailedItems = append(resp.FailedItems, fail)
			resp.FailedCount++
			continue
		}
		resp.SuccessCount++
	}

	return resp, nil
}

func (s *BookService) GetBookList(ctx context.Context, req *request.GetBookListRequest) (*response.GetBookListResponse, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 10
	}

	books, count, err := s.bookRepo.GetBookList(ctx, req)
	if err != nil {
		return nil, err
	}

	items := make([]response.BookListItem, 0, len(books))
	for _, book := range books {
		items = append(items, response.BookListItem{
			ID:           book.ID,
			Title:        book.Title,
			Author:       book.Author,
			ISBN:         book.ISBN,
			CategoryID:   book.CategoryID,
			CategoryName: book.Category.Name,
			Publisher:    book.Publisher,
			PublishDate:  book.PublishDate,
			Price:        &book.Price,
			Stock:        book.Stock,
			Available:    book.Stock - book.BorrowCount,
			CoverURL:     book.CoverURL,
			BorrowCount:  book.BorrowCount,
		})
	}

	return &response.GetBookListResponse{
		Total:      count,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int(math.Ceil(float64(count) / float64(req.Limit))),
		Books:      items,
	}, nil
}

func (s *BookService) GetBookDetails(ctx context.Context, id uint64) (*response.GetBookDetailsResponse, error) {
    book, err := s.bookRepo.GetBookByID(ctx, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, common.ErrBookNotFound
        }
        return nil, err
    }

    category, err := s.categoryRepo.GetCategoryByID(ctx, book.CategoryID)
    if err != nil {
        return nil, err
    }

    cateDetail := response.CategoryDetails{
        ID: category.ID,
        Name: category.Name,
    }

    resp := &response.GetBookDetailsResponse{
        ID:   uint(book.ID),
        Title: book.Title,
        Author: book.Author,
        ISBN: book.ISBN,
        Category: cateDetail,
        Publisher: book.Publisher,
        PublishDate: book.PublishDate,
        Price: book.Price,
        Stock: book.Stock,
        Available: book.Stock - book.BorrowCount,
        Description: book.Description,
        CoverURL: book.CoverURL,
        BorrowCount: book.BorrowCount,
        Rating: book.Rating,
        CreatedAt: book.CreatedAt,
        UpdatedAt: book.UpdatedAt,
    }

    return resp, nil
}

func (s *BookService) UpdateBook(ctx context.Context, req *request.UpdateBookRequest, id uint64) (*response.UpdateBookResponse, error) {
    if req.Author == nil && req.CategoryID == nil && req.CoverURL == nil && req.Description == nil && req.ISBN == nil && 
    req.Price == nil && req.PublishDate == nil && req.Publisher == nil && req.Stock == nil && req.Title == nil {
        return nil, common.ErrBadRequest
    }

    book, err := s.bookRepo.GetBookByID(ctx, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, common.ErrBookNotFound
        }
        return nil, err
    }

    Updates := make(map[string]interface{})
    if req.Author != nil {
        Updates["author"] = *req.Author
    }
    var categoryName *string    
    if req.CategoryID != nil {
        category, err := s.categoryRepo.GetCategoryByID(ctx, *req.CategoryID)
        if err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                return nil, common.ErrCategoryNotFound
            }
        return nil, err
        }
        categoryName = &category.Name
        Updates["category_id"] = *req.CategoryID
    }
    if req.CoverURL != nil {
        Updates["cover_url"] = *req.CoverURL
    }
    if req.Description != nil {
        Updates["description"] = *req.Description
    }
    if req.ISBN != nil {
        Updates["isbn"] = *req.ISBN
    }
    if req.Price != nil {
        Updates["price"] = *req.Price
    }
    if req.PublishDate != nil {
        Updates["publish_date"] = *req.PublishDate
    }
    if req.Publisher != nil {
        Updates["publisher"] = *req.Publisher
    }
    var available *int
    if req.Stock != nil {
        Updates["stock"] = *req.Stock
        avail := *req.Stock - book.BorrowCount
        available = &avail
    }
    if req.Title != nil {
        Updates["title"] = *req.Title
    }

    if err := s.bookRepo.UpdateBookFields(ctx, id, Updates); err != nil {
        return nil, err
    }

    resp := &response.UpdateBookResponse{
        ID:   book.ID,
        Title: req.Title,
        Author: req.Author,
        ISBN: req.ISBN,
        CategoryID: req.CategoryID,
        CategoryName: categoryName,
        Publisher: req.Publisher,
        PublishDate: req.PublishDate,
        Price: req.Price,
        Stock: req.Stock,
        Available: available,
        Description: req.Description,
        CoverURL: req.CoverURL,
        UpdatedAt: time.Now().UTC().Format(time.RFC3339),
    }

    return resp, nil
}

func (s *BookService) DeleteBook(ctx context.Context, id uint64) error {
	book, err := s.bookRepo.GetBookByID(ctx, id)
	if err != nil {
		return err
	}

	if book.BorrowCount > 0 {
		bizErr := common.NewBizError(400, "无法删除该图书", 400)
		details := make(map[string]interface{})
		details["reason"] = "存在未归还的借阅记录"
		details["unreturned_books"] = book.BorrowCount
		bizErr.WithDetails(details)
		return bizErr
	}

	err = s.bookRepo.DeleteBookByID(ctx, id)
	if err != nil {
		return err
	}

	return nil
}