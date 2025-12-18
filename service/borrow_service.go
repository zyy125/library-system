package service

import (
	"context"
	"errors"
	"library-system/common"
	"library-system/dto/request"
	"library-system/dto/response"
	"library-system/model"
	"library-system/repository"
	"time"
	"log"

	"gorm.io/gorm"
)

const (
    // 普通用户默认最大借阅数量
    DefaultBorrowLimit = 5

    // 默认借阅期限（天）
    DefaultBorrowDays = 30

    // 最大续借次数
    MaxRenewCount = 2
)

type BorrowService struct {
	borrowRepo *repository.BorrowRepository
	bookRepo *repository.BookRepository
	userRepo *repository.UserRepository
	overdueService *OverdueService
}

func NewBorrowService(
	borrowRepo *repository.BorrowRepository,
	bookRepo *repository.BookRepository,
	userRepo *repository.UserRepository,
	overdueService *OverdueService,
) *BorrowService {
	return &BorrowService{
		borrowRepo:      borrowRepo,
		bookRepo:       bookRepo,
		userRepo:       userRepo,
		overdueService: overdueService,
	}
}

func (s *BorrowService) BorrowBook(ctx context.Context, userID uint64, req *request.BorrowBookRequest) (*response.BorrowBookResponse, error) {
	// 借书前先刷新逾期记录
	err := s.overdueService.RefreshSingleUserOverdue(ctx, userID)
	if err != nil {
		// 记录日志，但不阻塞查询
		log.Printf("检查逾期失败: %v", err)
	}

	var resp *response.BorrowBookResponse

	err = s.borrowRepo.DB().Transaction(func(tx *gorm.DB) error {
		user, err := s.userRepo.GetUserByIDWithLock(ctx, tx, userID) 
		if err != nil {
			return err
		}

		if user.BorrowingCount >= user.BorrowLimit {
			return common.ErrBorrowLimitReached
		}

		if user.OverdueCount > 0 {
			return common.ErrHasOverdueBooks
		}	

		book, err := s.bookRepo.GetBookByIDWithLock(ctx, tx, req.BookId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return common.ErrBookNotFound
			}
			return err
		}

		if book.Stock <= book.BorrowCount {
			return common.ErrBookOutOfStock
		}

		nowDate := time.Now().UTC()
		var dueDate time.Time
		if req.BorrowDays != nil {
			dueDate = nowDate.AddDate(0, 0, *req.BorrowDays)
		} else {
			dueDate = nowDate.AddDate(0, 0, DefaultBorrowDays)
		}

		borrow := model.BorrowRecord{
			BookID: req.BookId,
			UserID: userID,
			DueDate: dueDate,
			Status: "borrowed",
		}

		if err := s.borrowRepo.CreateBorrowRecord(ctx, tx, &borrow); err != nil {
			return err
		}

		if err := s.bookRepo.IncrementBorrowCount(ctx, tx, req.BookId); err != nil {
			return err
		}

		if err := s.userRepo.IncrementBorrowingCount(ctx, tx, userID, 1); err != nil {
			return err
		}

		resp = &response.BorrowBookResponse{
		ID: borrow.ID,
		BookID: borrow.BookID,
		BookTitle: book.Title,
		UserID: userID,
		Username: user.Username,
		BorrowDate: nowDate,
		DueDate: dueDate,
		Status: borrow.Status,
		RenewCount: 0,
		MaxRenewCount: MaxRenewCount,
	}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}