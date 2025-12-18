package service

import (
	"context"
	"errors"
	"library-system/common"
	"library-system/dto/request"
	"library-system/dto/response"
	"library-system/model"
	"library-system/repository"
	"log"
	"time"

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
	borrowRepo     *repository.BorrowRepository
	bookRepo       *repository.BookRepository
	userRepo       *repository.UserRepository
	overdueService *OverdueService
}

func NewBorrowService(
	borrowRepo *repository.BorrowRepository,
	bookRepo *repository.BookRepository,
	userRepo *repository.UserRepository,
	overdueService *OverdueService,
) *BorrowService {
	return &BorrowService{
		borrowRepo:     borrowRepo,
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
			BookID:  req.BookId,
			UserID:  userID,
			DueDate: dueDate,
			Status:  "borrowed",
		}

		if err := s.borrowRepo.CreateBorrowRecord(ctx, tx, &borrow); err != nil {
			return err
		}

		if err := s.bookRepo.IncreaseBorrowCount(ctx, tx, req.BookId, 1); err != nil {
			return err
		}

		if err := s.userRepo.IncreaseBorrowingCount(ctx, tx, userID, 1); err != nil {
			return err
		}

		resp = &response.BorrowBookResponse{
			ID:            borrow.ID,
			BookID:        borrow.BookID,
			BookTitle:     book.Title,
			UserID:        userID,
			Username:      user.Username,
			BorrowDate:    nowDate,
			DueDate:       dueDate,
			Status:        borrow.Status,
			RenewCount:    0,
			MaxRenewCount: MaxRenewCount,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *BorrowService) ReturnBook(ctx context.Context, borrowID uint64, req *request.ReturnBookRequest) (*response.ReturnBookResponse, error) {
	borrow, err := s.borrowRepo.GetBorrowRecordByID(ctx, borrowID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrBorrowNotFound
		}
		return nil, err
	}
	if borrow.ReturnDate != nil || borrow.Status == "returned" {
		return nil, &common.BizError{
			Code:    400,
			Message: "该图书已归还",
		}
	}

	var isOverdue bool
	now := time.Now().UTC()
	if now.After(borrow.DueDate) {
		isOverdue = true
	} else {
		isOverdue = false
	}

	overdueDays := 0

	var resp *response.ReturnBookResponse

	err = s.borrowRepo.DB().Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{
			"return_date": now,
			"status":      "returned",
		}
		var fine float64
		if isOverdue {
			// 计算逾期天数
			overdueDays = int(now.Sub(borrow.DueDate).Hours() / 24)
			// 计算罚金（1元/天）
			fine = float64(overdueDays) * 1.0
			updates["fine"] = fine

			if borrow.Status == "overdue" {
				if err := s.userRepo.DecreaseOverDueCount(ctx, tx, borrow.UserID, 1); err != nil {
					return err
				}
			}
		}

		if err := s.userRepo.DecreaseBorrowingCount(ctx, tx, borrow.UserID, 1); err != nil {
			return err
		}
		if err := s.bookRepo.DecreaseBorrowCount(ctx, tx, borrow.BookID, 1); err != nil {
			return err
		}
		if err := s.borrowRepo.UpdateFields(ctx, tx, borrowID, updates); err != nil {
			return err
		}

		resp = &response.ReturnBookResponse{
			ID:          borrowID,
			BookID:      borrow.BookID,
			UserID:      borrow.UserID,
			BorrowDate:  borrow.BorrowDate,
			DueDate:     borrow.DueDate,
			ReturnDate:  now,
			Status:      "returned",
			IsOverdue:   isOverdue,
			OverdueDays: overdueDays,
			Fine:        fine,
			Condition:   req.Condition,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}
