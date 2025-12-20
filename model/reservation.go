package model

import(
	"time"
)

type Reservation struct {
    ID        uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
    BookID    uint64     `json:"book_id" gorm:"index: idx_book;not null"`
    UserID    uint64     `json:"user_id" gorm:"index:idx_user;not null"`
    Status    string     `json:"status" gorm:"type:enum('waiting','available','cancelled','expired','fulfilled');default:'waiting';index:idx_status"`
    
    // 预约时间
    ReservedAt time.Time  `json:"reserved_at" gorm:"autoCreateTime"`
    
    // 通知时间（图书可借时）
    NotifiedAt *time.Time `json:"notified_at"`
    
    // 过期时间（通知后48小时）
    ExpiresAt  *time.Time `json:"expires_at" gorm:"index: idx_expires"`
    
    // 完成时间（实际借书时）
    FulfilledAt *time.Time `json:"fulfilled_at"`
    
    // 取消时间
    CancelledAt *time.Time `json:"cancelled_at"`
    
    CreatedAt  time.Time  `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt  time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

    Book Book `gorm:"foreignKey:BookID"`
    User User `gorm:"foreignKey:UserID"`
}

// 预约状态说明
const (
    ReservationStatusWaiting   = "waiting"    // 等待中（排队）
    ReservationStatusAvailable = "available"  // 可借阅（已通知）
    ReservationStatusCancelled = "cancelled"  // 已取消
    ReservationStatusExpired   = "expired"    // 已过期（48小时未借）
    ReservationStatusFulfilled = "fulfilled"  // 已完成（已借书）
)