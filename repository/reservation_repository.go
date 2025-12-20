package repository

import (
    "context"
    "library-system/model"
    "time"

    "gorm.io/gorm"
)

type ReservationRepository struct {
    db *gorm.DB
}

func (r *ReservationRepository) DB() *gorm.DB {
	return r.db
} 

func NewReservationRepository(db *gorm.DB) *ReservationRepository {
    return &ReservationRepository{db: db}
}

// CreateReservation 创建预约
func (r *ReservationRepository) CreateReservation(ctx context.Context, reservation *model. Reservation) error {
    return gorm.G[model.Reservation](r.db).Create(ctx, reservation)
}

// GetReservationByID 根据ID获取预约
func (r *ReservationRepository) GetReservationByID(ctx context.Context, id uint64) (model.Reservation, error) {
    return gorm.G[model.Reservation](r.db).Where("id = ?", id).
	Preload("Book", func(db gorm.PreloadBuilder) error {return nil}).
	Preload("User", func(db gorm.PreloadBuilder) error {return nil}).
	First(ctx)
}

// GetUserReservationForBook 检查用户是否已预约该图书
func (r *ReservationRepository) GetUserReservationForBook(ctx context.Context, userID, bookID uint64) (model.Reservation, error) {
	return gorm.G[model.Reservation](r.db).Where("user_id = ? AND book_id = ?  AND status = ?", 
            userID, bookID, model.ReservationStatusWaiting).First(ctx)
}

// GetMyReservations 获取我的预约列表
func (r *ReservationRepository) GetMyReservations(ctx context.Context, userID uint64) ([]model.Reservation, error) {
    var reservations []model.Reservation
    err := r.db.WithContext(ctx).
        Preload("Book").
        Where("user_id = ? AND status IN ? ", userID, 
            []string{model.ReservationStatusWaiting, model.ReservationStatusAvailable}).
        Order("reserved_at ASC").
        Find(&reservations).Error
    return reservations, err
}

// GetQueuePosition 获取预约排队位置
func (r *ReservationRepository) GetQueuePosition(ctx context.Context, reservationID uint64) (int, error) {
    var reservation model.Reservation
    if err := r.db.WithContext(ctx).First(&reservation, reservationID).Error; err != nil {
        return 0, err
    }

    var position int64
    err := r.db. WithContext(ctx).
        Model(&model.Reservation{}).
        Where("book_id = ? AND status = ?  AND reserved_at < ?", 
            reservation.BookID, model.ReservationStatusWaiting, reservation.ReservedAt).
        Count(&position).Error

    return int(position) + 1, err
}

// GetNextWaitingReservation 获取下一个等待的预约（图书归还时调用）
func (r *ReservationRepository) GetNextWaitingReservation(ctx context. Context, bookID uint64) (*model.Reservation, error) {
    var reservation model.Reservation
    err := r.db.WithContext(ctx).
        Where("book_id = ? AND status = ? ", bookID, model.ReservationStatusWaiting).
        Order("reserved_at ASC"). // 先预约先得
        First(&reservation).Error
    return &reservation, err
}

// UpdateReservationStatus 更新预约状态
func (r *ReservationRepository) UpdateReservationStatus(ctx context.Context, tx *gorm.DB, id uint64, updates map[string]interface{}) error {
    return tx.WithContext(ctx).
        Model(&model.Reservation{}).
        Where("id = ? ", id).
        Updates(updates).Error
}

// GetExpiredReservations 获取过期的预约（定时任务用）
func (r *ReservationRepository) GetExpiredReservations(ctx context.Context) ([]model.Reservation, error) {
    var reservations []model.Reservation
    now := time.Now()
    
    err := r.db.WithContext(ctx).
        Where("status = ? AND expires_at < ? ", model.ReservationStatusAvailable, now).
        Find(&reservations).Error
    
    return reservations, err
}

// CountActiveReservations 统计活跃预约数
func (r *ReservationRepository) CountActiveReservations(ctx context.Context) (int64, error) {
    var count int64
    err := r. db.WithContext(ctx).
        Model(&model.Reservation{}).
        Where("status IN ? ", []string{model.ReservationStatusWaiting, model. ReservationStatusAvailable}).
        Count(&count).Error
    return count, err
}

// HasActiveReservation 检查图书是否有活跃预约
func (r *ReservationRepository) HasActiveReservation(ctx context.Context, bookID uint64) (bool, error) {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&model.Reservation{}).
        Where("book_id = ?  AND status IN ?", bookID, 
            []string{model.ReservationStatusWaiting, model.ReservationStatusAvailable}).
        Count(&count).Error
    return count > 0, err
}
