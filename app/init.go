package app

import(
	"library-system/database"
	"library-system/repository"
	"library-system/controller"
	"library-system/service"
)

func InitApp() (*controller.UserController, error) {
	db, err := database.InitMySQL()
	if err != nil {
		return nil, err
	}

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userCtl := controller.NewUserController(userService)

	return userCtl, nil
}