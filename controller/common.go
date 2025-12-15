package controller

type Controller struct {
	UserController *UserController
	BookController *BookController
}

func NewController(userCtl *UserController, bookCtl *BookController) *Controller {
	return &Controller{
		UserController: userCtl,
		BookController: bookCtl,
	}
}