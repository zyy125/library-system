package scheduler

import (
    "context"
    "library-system/service"
    "log"
    "time"

    "github.com/robfig/cron/v3"
)

type ReservationScheduler struct {
    reservationService *service.ReservationService
    cron               *cron. Cron
}

func NewReservationScheduler(reservationService *service.ReservationService) *ReservationScheduler {
    return &ReservationScheduler{
        reservationService: reservationService,
        cron:                cron.New(),
    }
}

// Start 启动定时任务
func (s *ReservationScheduler) Start(cronExpr string) error {
    // 每小时检查一次过期预约
    _, err := s.cron.AddFunc(cronExpr, func() {
        ctx := context.Background()
        startTime := time.Now()

        log. Println("[定时任务] 开始处理过期预约...")

        count, err := s.reservationService.ProcessExpiredReservations(ctx)
        if err != nil {
            log.Printf("[定时任务] 处理失败: %v\n", err)
            return
        }

        duration := time.Since(startTime)
        log.Printf("[定时任务] 完成处理，过期预约数: %d，耗时 %v\n", count, duration)
    })

    if err != nil {
        return err
    }

    s.cron.Start()
    log.Println("[定时任务] 预约过期检查已启动（每小时执行）")

    return nil
}

// Stop 停止定时任务
func (s *ReservationScheduler) Stop() {
    if s.cron != nil {
        s.cron.Stop()
        log.Println("[定时任务] 预约过期检查已停止")
    }
}