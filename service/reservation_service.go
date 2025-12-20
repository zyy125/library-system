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

type ReservationService struct {
	reservationRepo *repository.ReservationRepository
	bookRepo        *repository.BookRepository
	userRepo        *repository.UserRepository
}

func NewReservationService(
	reservationRepo *repository.ReservationRepository,
	bookRepo *repository.BookRepository,
	userRepo *repository.UserRepository,
) *ReservationService {
	return &ReservationService{
		reservationRepo: reservationRepo,
		bookRepo:        bookRepo,
		userRepo:        userRepo,
	}
}

func (s *ReservationService) CreateReservation(ctx context.Context, userID uint64, req *request.CreateReservationRequest) (*response.CreateReservationResponse, error) {
	// 1. 检查图书是否存在
	book, err := s.bookRepo.GetBookByID(ctx, req.BookID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrBookNotFound
		}
		return nil, err
	}

	// 2. 检查图书是否有库存（有库存不能预约）
	available := book.Stock - book.BorrowCount
	if available > 0 {
		return nil, common.ErrReservationFailed // 30007:  预约失败，图书有库存
	}

	// 3. 检查用户是否已预约该图书
	existingReservation, err := s.reservationRepo.GetUserReservationForBook(ctx, userID, req.BookID)
	if err == nil && existingReservation.ID > 0 {
		return nil, common.ErrHasReservationed
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 4. 创建预约记录
	reservation := &model.Reservation{
		BookID:     req.BookID,
		UserID:     userID,
		Status:     model.ReservationStatusWaiting,
		ReservedAt: time.Now(),
	}

	if err := s.reservationRepo.CreateReservation(ctx, reservation); err != nil {
		return nil, err
	}

	// 5. 获取排队位置
	position, err := s.reservationRepo.GetQueuePosition(ctx, reservation.ID)
	if err != nil {
		log.Printf("获取排队位置失败:  %v", err)
		position = 0
	}

	// 6. 构建响应
	expiresAt := time.Now().Add(48 * time.Hour)
	resp := &response.CreateReservationResponse{
		ID:            reservation.ID,
		BookID:        book.ID,
		BookTitle:     book.Title,
		UserID:        userID,
		Status:        reservation.Status,
		QueuePosition: position,
		ReservedAt:    reservation.ReservedAt,
		ExpiresAt:     expiresAt,
	}

	return resp, nil
}

// CancelReservation 取消预约
func (s *ReservationService) CancelReservation(ctx context.Context, userID, reservationID uint64) error {
	// 1. 获取预约记录
	reservation, err := s.reservationRepo.GetReservationByID(ctx, reservationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.ErrReservationNotFound
		}
		return err
	}

	// 2. 检查权限
	if reservation.UserID != userID {
		return &common.BizError{
			Code:    403,
			Message: "无权操作此预约",
		}
	}

	// 3. 检查状态
	if reservation.Status != model.ReservationStatusWaiting &&
		reservation.Status != model.ReservationStatusAvailable {
		return common.ErrReservationHasCanceled
	}

	// 4. 更新状态为已取消
	now := time.Now()
	updates := map[string]interface{}{
		"status":       model.ReservationStatusCancelled,
		"cancelled_at": now,
	}

	return s.reservationRepo.UpdateReservationStatus(ctx, nil, reservationID, updates)
}

// GetMyReservations 获取我的预约列表
func (s *ReservationService) GetMyReservations(ctx context.Context, userID uint64) (*response.GetMyReservationsResponse, error) {
	reservations, err := s.reservationRepo.GetMyReservations(ctx, userID)
	if err != nil {
		return nil, err
	}

	items := make([]response.ReservationItem, 0, len(reservations))
	for _, reservation := range reservations {
		// 获取排队位置
		position := 0
		if reservation.Status == model.ReservationStatusWaiting {
			position, _ = s.reservationRepo.GetQueuePosition(ctx, reservation.ID)
		}

		item := response.ReservationItem{
			ID: reservation.ID,
			Book: response.ReservationBookResponse{
				ID:       reservation.Book.ID,
				Title:    reservation.Book.Title,
				Author:   reservation.Book.Author,
				CoverURL: reservation.Book.CoverURL,
			},
			Status:        reservation.Status,
			QueuePosition: position,
			ReservedAt:    reservation.ReservedAt,
			ExpiresAt:     reservation.ExpiresAt,
		}

		items = append(items, item)
	}

	return &response.GetMyReservationsResponse{
		Reservations: items,
	}, nil
}

// NotifyNextReservation 通知下一个预约者（图书归还时调用）
func (s *ReservationService) NotifyNextReservation(ctx context.Context, tx *gorm.DB, bookID uint64) error {
	// 1. 获取下一个等待的预约
	reservation, err := s.reservationRepo.GetNextWaitingReservation(ctx, bookID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 没有预约，正常情况
			return nil
		}
		return err
	}

	// 2. 更新预约状态为可借阅
	now := time.Now()
	expiresAt := now.Add(48 * time.Hour)

	updates := map[string]interface{}{
		"status":      model.ReservationStatusAvailable,
		"notified_at": now,
		"expires_at":  expiresAt,
	}

	if err := s.reservationRepo.UpdateReservationStatus(ctx, tx, reservation.ID, updates); err != nil {
		return err
	}

	// 3. 发送通知（这里简化处理，实际项目应该调用通知服务）
	log.Printf("📧 通知用户 %d:  您预约的图书《%s》已可借阅，请在 %s 前借阅",
		reservation.UserID, reservation.Book.Title, expiresAt.Format("2006-01-02 15:04"))

	// TODO: 集成邮件/短信通知服务
	// s.notificationService.SendReservationNotification(reservation)

	return nil
}

// ProcessExpiredReservations 处理过期预约（定时任务）
func (s *ReservationService) ProcessExpiredReservations(ctx context.Context) (int, error) {
	expiredReservations, err := s.reservationRepo.GetExpiredReservations(ctx)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, reservation := range expiredReservations {
		// 更新为已过期
		updates := map[string]interface{}{
			"status": model.ReservationStatusExpired,
		}

		if err := s.reservationRepo.UpdateReservationStatus(ctx, nil, reservation.ID, updates); err != nil {
			log.Printf("更新预约%d状态失败: %v", reservation.ID, err)
			continue
		}

		count++

		// 通知下一个预约者
		if err := s.NotifyNextReservation(ctx, nil, reservation.BookID); err != nil {
			log.Printf("通知下一个预约者失败: %v", err)
		}
	}

	return count, nil
}
