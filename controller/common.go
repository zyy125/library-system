package controller

type Controller struct {
	UserController *UserController
	BookController *BookController
	BorrowController *BorrowController
	ReservationController *ReservationController
	CategoryController    *CategoryController
	StatsController       *StatsController 
}

func NewController(userCtl *UserController, 
	bookCtl *BookController, borrowCtl *BorrowController, 
	reservationCtl *ReservationController, categoryCtl *CategoryController, statsCtl *StatsController,
	) *Controller {
	return &Controller{
		UserController: userCtl,
		BookController: bookCtl,
		BorrowController: borrowCtl,
		ReservationController: reservationCtl,
		CategoryController:    categoryCtl,
		StatsController:       statsCtl,
	}
}