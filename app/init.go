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
	Scheduler  *scheduler.OverdueScheduler
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

	overdueService := service.NewOverdueService(borrowRepo, userRepo)
	userService := service.NewUserService(userRepo, overdueService)
	bookService := service.NewBookService(bookRepo, cateRepo)
	borrowService := service.NewBorrowService(borrowRepo, bookRepo, userRepo, overdueService)

	overdueScheduler := scheduler.NewOverdueScheduler(overdueService)
	userCtl := controller.NewUserController(userService)
	bookCtl := controller.NewBookController(bookService)
	borrowCtl := controller.NewBorrowController(borrowService)

	ctl := controller.NewController(userCtl, bookCtl, borrowCtl)

	app := &App{
		Controller:  ctl,
		Scheduler:   overdueScheduler,
	}

	// 启动定时任务（每天凌晨2点执行）
	if err := overdueScheduler.Start("0 2 * * *"); err != nil {
		return nil, fmt.Errorf("定时任务启动失败: %v", err)
	}

	return app, nil
}