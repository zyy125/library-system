package scheduler

type Scheduler struct {
	OverdueScheduler *OverdueScheduler
	ReservationScheduler *ReservationScheduler
}