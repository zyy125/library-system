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
	"math"
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
	reservationService *ReservationService
	reservationRepo *repository.ReservationRepository
}

func NewBorrowService(
	borrowRepo *repository.BorrowRepository,
	bookRepo *repository.BookRepository,
	userRepo *repository.UserRepository,
	reservationRepo *repository.ReservationRepository,
	reservationService *ReservationService,
	overdueService *OverdueService,
) *BorrowService {
	return &BorrowService{
		borrowRepo:     borrowRepo,
		bookRepo:       bookRepo,
		userRepo:       userRepo,
		reservationRepo: reservationRepo,
		reservationService: reservationService,
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
		// 检查该用户是否有该书的有效预约
        reservation, err := s.reservationRepo.GetUserReservationForBook(ctx, userID, req.BookId)
        if err == nil && reservation.Status == model.ReservationStatusAvailable {
            // 用户有有效预约，标记为已完成
            now := time.Now()
            updates := map[string]interface{}{
                "status":       model.ReservationStatusFulfilled,
                "fulfilled_at": now,
            }
            s.reservationRepo.UpdateReservationStatus(ctx, tx, reservation.ID, updates)
        } else if errors.Is(err, gorm. ErrRecordNotFound) {
            // 没有预约，检查是否有其他人预约
            hasReservation, _ := s.reservationRepo.HasActiveReservation(ctx, req.BookId)
            if hasReservation {
                // 有人预约但不是当前用户
                return common.ErrHasReservation
            }
        }
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

		// 新增：还书后通知下一个预约者
        if err := s.reservationService.NotifyNextReservation(ctx, tx, borrow.BookID); err != nil {
            log.Printf("通知预约者失败: %v", err)
            // 不中断还书流程
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

func (s *BorrowService) RenewBorrow(ctx context.Context, userID uint64, borrowID uint64, req *request.RenewBorrowRequest) (*response.RenewBorrowResponse, error) {
	// 刷新逾期记录
	err := s.overdueService.RefreshSingleUserOverdue(ctx, userID)
	if err != nil {
		// 记录日志，但不阻塞查询
		log.Printf("检查逾期失败: %v", err)
	}

	record, err := s.borrowRepo.GetBorrowRecordByIDWithBook(ctx, borrowID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrBorrowNotFound
		}
		return nil, err
	}

	if record.RenewCount >= MaxRenewCount {
		return nil, common.ErrRenewLimitReached
	}

	if record.Status == "overdue" {
		return nil, common.ErrCannotRenewOverdue
	}

	if record.Status == "returned" {
		return nil, &common.BizError{
			Code:    400,
			Message: "该图书已归还，无法续借",
		}
	}

	var resp *response.RenewBorrowResponse

	renewDays := DefaultBorrowDays
	if req.RenewDays != nil {
		renewDays = *req.RenewDays
	}

	newDueDate := record.DueDate.AddDate(0, 0, renewDays)
	renewCount := record.RenewCount + 1
	updates := map[string]interface{}{
		"renew_count": renewCount,
		"due_date":    newDueDate,
	}
	if err := s.borrowRepo.UpdateFields(ctx, s.borrowRepo.DB(), borrowID, updates); err != nil {
		return nil, err
	}

	resp = &response.RenewBorrowResponse{
		ID:              record.ID,
		BookID:          record.BookID,
		OriginalDueDate: record.DueDate,
		NewDueDate:      newDueDate,
		BookTitle:       record.Book.Title,
		RenewCount:      renewCount,
		MaxRenewCount:   MaxRenewCount,
	}

	return resp, nil
}

func DaysFromToday(t time.Time) int {
	now := time.Now().In(t.Location())

	today := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0,
		now.Location(),
	)

	target := time.Date(
		t.Year(), t.Month(), t.Day(),
		0, 0, 0, 0,
		t.Location(),
	)

	return int(target.Sub(today).Hours() / 24)
}

func (s *BorrowService) GetBorrowRecordList(ctx context.Context, req *request.GetBorrowRecordListRequest) (*response.GetBorrowRecordListResponse, error) {
	if req.UserID != nil {
		err := s.overdueService.RefreshSingleUserOverdue(ctx, *req.UserID)
		if err != nil {
			// 记录日志，但不阻塞查询
			log.Printf("检查用户%d逾期失败: %v", *req.UserID, err)
		}
	}

	if req.Page == 0 {
		req.Page = 1
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	records, total, err := s.borrowRepo.GetBorrowRecordList(ctx, req)
	if err != nil {
		return nil, err
	}

	items := make([]response.GetBorrowRecordItemResponse, 0, len(records))
	for _, record := range records {
		item := response.GetBorrowRecordItemResponse{
			ID: record.ID,
			Book: response.GetBorrowRecordListBookResponse{
				ID:       record.Book.ID,
				Title:    record.Book.Title,
				Author:   record.Book.Author,
				ISBN:     record.Book.ISBN,
				CoverURL: record.Book.CoverURL,
			},
			User: &response.GetBorrowRecordListUserResponse{
				ID:       record.User.ID,
				Username: record.User.Username,
			},
			BorrowDate: record.BorrowDate,
			DueDate:    record.DueDate,
			ReturnDate: record.ReturnDate,
			Status:     record.Status,
			RenewCount: record.RenewCount,
			Fine:       record.Fine,
		}

		if record.Status == "overdue" {
			item.IsOverdue = true
			item.OverdueDays = DaysFromToday(record.DueDate)
		}
		if record.Status == "borrowed" {
			item.IsOverdue = false
			item.DaysUntilDue = DaysFromToday(record.DueDate)
		}

		if record.RenewCount < MaxRenewCount && record.Status == "borrowed" {
			item.CanRenew = true
		}
		items = append(items, item)
	}

	return &response.GetBorrowRecordListResponse{
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int(math.Ceil(float64(total) / float64(req.Limit))),
		Records:    items,
	}, nil
}

func (s *BorrowService) GetCurrentRecord(ctx context.Context, userID uint64) (*response.GetCurrentRecordResponse, error) {
	records, err := s.borrowRepo.GetRecordByUserIDWithPreload(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, common.ErrBorrowNotFound
	}

	var items []response.GetBorrowRecordItemResponse
	var totalFine float64 = 0

	for _, record := range records {
		totalFine += record.Fine
		item := response.GetBorrowRecordItemResponse{
			ID: record.ID,
			User: nil,
			Book: response.GetBorrowRecordListBookResponse{
				ID:       record.Book.ID,
				Title:    record.Book.Title,
				Author:   record.Book.Author,
				ISBN:     record.Book.ISBN,
				CoverURL: record.Book.CoverURL,
			},
			BorrowDate: record.BorrowDate,
			DueDate:    record.DueDate,
			ReturnDate: record.ReturnDate,
			Status:     record.Status,
			RenewCount: record.RenewCount,
			Fine:       record.Fine,
		}

		if record.Status == "overdue" {
			item.IsOverdue = true
			item.OverdueDays = DaysFromToday(record.DueDate)
		}
		if record.Status == "borrowed" {
			item.IsOverdue = false
			item.DaysUntilDue = DaysFromToday(record.DueDate)
		}

		if record.RenewCount < MaxRenewCount && record.Status == "borrowed" {
			item.CanRenew = true
		}
		items = append(items, item)
	}

	user := records[0].User

	return &response.GetCurrentRecordResponse{
		BorrowingCount: user.BorrowingCount,
		BorrowLimit:user.BorrowLimit,
		OverdueCount: user.OverdueCount,
		TotalFine: totalFine,
		Records: items,
	}, nil
}
