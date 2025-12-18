package controller

type Controller struct {
	UserController *UserController
	BookController *BookController
	BorrowController *BorrowController

}

func NewController(userCtl *UserController, bookCtl *BookController, borrowCtl *BorrowController) *Controller {
	return &Controller{
		UserController: userCtl,
		BookController: bookCtl,
		BorrowController: borrowCtl,
	}
}