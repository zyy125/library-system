package app

import(
	"library-system/database"
	"library-system/repository"
	"library-system/controller"
	"library-system/service"
	"library-system/scheduler"

	"fmt"
)

type App struct {
	Controller *controller.Controller
	Scheduler  *scheduler.Scheduler
}

func InitApp() (*App, error) {
	db, err := database.InitMySQL()
	if err != nil {
		return nil, err
	}
	rdb, err := database.InitRedis()
	if err != nil {
		return nil, err
	}

	repository.NewRedis(rdb)

	userRepo := repository.NewUserRepository(db)
	bookRepo := repository.NewBookRepository(db)
	borrowRepo := repository.NewBorrowRepository(db)
	cateRepo := repository.NewCategoryRepository(db)
	reservationRepo := repository.NewReservationRepository(db)

	overdueService := service.NewOverdueService(borrowRepo, userRepo)
	userService := service.NewUserService(userRepo, overdueService)
	bookService := service.NewBookService(bookRepo, cateRepo)
	reservationService := service.NewReservationService(reservationRepo, bookRepo, userRepo)
	borrowService := service.NewBorrowService(borrowRepo, bookRepo, userRepo, reservationRepo, reservationService, overdueService)
	
	overdueScheduler := scheduler.NewOverdueScheduler(overdueService)
	reservationScheduler := scheduler.NewReservationScheduler(reservationService)
	userCtl := controller.NewUserController(userService)
	bookCtl := controller.NewBookController(bookService)
	borrowCtl := controller.NewBorrowController(borrowService)
	reservationCtl := controller.NewReservationController(reservationService)

	ctl := controller.NewController(userCtl, bookCtl, borrowCtl, reservationCtl)


	scheduler := &scheduler.Scheduler{OverdueScheduler: overdueScheduler, ReservationScheduler: reservationScheduler}
	app := &App{
		Controller:  ctl,
		Scheduler:   scheduler,
	}

	// 启动定时任务
	if err := overdueScheduler.Start("0 2 * * *"); err != nil {
		return nil, fmt.Errorf("定时任务启动失败: %v", err)
	}

	if err := reservationScheduler.Start("0 * * * *"); err != nil {
		return nil, fmt.Errorf("定时任务启动失败: %v", err)
	}
	return app, nil
}