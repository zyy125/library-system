package scheduler

import (
	"context"
	"library-system/service"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

// OverdueScheduler 逾期检查定时任务调度器
type OverdueScheduler struct {
	overdueService *service.OverdueService
	cron           *cron.Cron
}

// NewOverdueScheduler 创建调度器
func NewOverdueScheduler(overdueService *service.OverdueService) *OverdueScheduler {
	return &OverdueScheduler{
		overdueService: overdueService,
		cron:           cron.New(),
	}
}

// Start 启动定时任务
func (s *OverdueScheduler) Start(cronExpr string) error {
	// 添加定时任务
	_, err := s.cron.AddFunc(cronExpr, func() {
		ctx := context.Background()
		startTime := time.Now()

		log.Println("[定时任务] 开始检查逾期记录...")

		// 执行逾期检查
		count, err := s.overdueService.RefreshAllUsersOverdue(ctx)
		if err != nil {
			log.Printf("[定时任务] 检查失败: %v\n", err)
			return
		}

		duration := time.Since(startTime)
		log.Printf("[定时任务] 完成检查，更新了 %d 条记录，耗时 %v\n", count, duration)
	})

	if err != nil {
		return err
	}

	// 启动调度器
	s.cron.Start()
	log.Printf("[定时任务] 逾期检查已启动，执行计划:  %s\n", cronExpr)

	return nil
}

// Stop 停止定时任务
func (s *OverdueScheduler) Stop() {
	if s.cron != nil {
		s.cron.Stop()
		log.Println("[定时任务] 逾期检查已停止")
	}
}