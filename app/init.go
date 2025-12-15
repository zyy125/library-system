package app

import(
	"library-system/database"
	"library-system/repository"
	"library-system/controller"
	"library-system/service"
)

func InitApp() (*controller.Controller, error) {
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
	userService := service.NewUserService(userRepo)
	userCtl := controller.NewUserController(userService)

	bookRepo := repository.NewBookRepository(db)
	cateRepo := repository.NewCategoryRepository(db)
	bookService := service.NewBookService(bookRepo, cateRepo)
	bookCtl := controller.NewBookController(bookService)

	ctl := controller.NewController(userCtl, bookCtl)
	
	return ctl, nil
}