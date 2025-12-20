package controller

type Controller struct {
	UserController *UserController
	BookController *BookController
	BorrowController *BorrowController
	ReservationController *ReservationController
}

func NewController(userCtl *UserController, bookCtl *BookController, borrowCtl *BorrowController, reservationCtl *ReservationController) *Controller {
	return &Controller{
		UserController: userCtl,
		BookController: bookCtl,
		BorrowController: borrowCtl,
		ReservationController: reservationCtl,
	}
}