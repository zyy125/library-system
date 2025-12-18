package service

import (
	"context"
	"library-system/model"
	"library-system/repository"
	"time"

	"gorm.io/gorm"
)

// OverdueService 逾期检查服务
type OverdueService struct {
	borrowRepo *repository.BorrowRepository
	userRepo   *repository.UserRepository
}

// NewOverdueService 创建逾期服务实例
func NewOverdueService(
	borrowRepo *repository.BorrowRepository,
	userRepo *repository.UserRepository,
) *OverdueService {
	return &OverdueService{
		borrowRepo: borrowRepo,
		userRepo:   userRepo,
	}
}

func (s *OverdueService) RefreshAllUsersOverdue(ctx context.Context) (int, error) {
	now := time.Now()
	updatedCount := 0

	err := s.borrowRepo.DB().Transaction(func(tx *gorm.DB) error {
		dueRecords, err := s.borrowRepo.GetAllDueRecord(ctx, tx, now)
		if err != nil {
			return err
		}

		if len(dueRecords) == 0 {
			return nil
		}

		userOverdueMap := make(map[uint64]int)

		for _, record := range dueRecords {
			// 计算逾期天数
			overdueDays := int(now.Sub(record.DueDate).Hours() / 24)
			// 计算罚金（1元/天）
			fine := float64(overdueDays) * 1.0

			updates := map[string]interface{}{
				"status": "overdue",
				"fine":   fine,
			}
			if err := s.borrowRepo.UpdateFields(ctx, tx, record.ID, updates); err != nil {
				return err
			}
			userOverdueMap[record.UserID]++
			updatedCount++
		}

		for userID, count := range userOverdueMap {
			if err := s.userRepo.IncreaseOverDueCount(ctx, tx, userID, count); err != nil {
				return err
			}
		}

		return nil
	})

	return updatedCount, err
}

func (s *OverdueService) RefreshSingleUserOverdue(ctx context.Context, userID uint64) error {
	now := time.Now()

	// 使用事务
	err := s.borrowRepo.DB().Transaction(func(tx *gorm.DB) error {
		// 1. 查找该用户所有逾期但状态未更新的记录
		var overdueRecords []model.BorrowRecord
		overdueRecords, err := s.borrowRepo.GetDueRecordByUserID(ctx, tx, userID, now)
		if err != nil {
			return err
		}

		if len(overdueRecords) == 0 {
			// 没有新的逾期记录
			return nil
		}

		overdueCount := len(overdueRecords)

		// 2. 更新这些记录
		for _, record := range overdueRecords {
			// 计算逾期天数
			overdueDays := int(now.Sub(record.DueDate).Hours() / 24)
			// 计算罚金（1元/天）
			fine := float64(overdueDays) * 1.0

			updates := map[string]interface{}{
				"status": "overdue",
				"fine":   fine,
			}
			if err := s.borrowRepo.UpdateFields(ctx, tx, record.ID, updates); err != nil {
				return err
			}
		}

		// 3. 更新用户逾期计数
		if err := s.userRepo.IncreaseOverDueCount(ctx, tx, userID, overdueCount); err != nil {
			return err
		}

		return nil
	})

	return err
}
