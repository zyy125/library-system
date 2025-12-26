package controller

type Controller struct {
	UserController        *UserController
	BookController        *BookController
	BorrowController      *BorrowController
	ReservationController *ReservationController
	CategoryController    *CategoryController
	StatsController       *StatsController
}

type Option func(*Controller)

func WithUser(userCtl *UserController) Option {
	return func(c *Controller) {
		c.UserController = userCtl
	}
}

func WithBook(bookCtl *BookController) Option {
	return func(c *Controller) {
		c.BookController = bookCtl
	}
}

func WithBorrow(borrowCtl *BorrowController) Option {
	return func(c *Controller) {
		c.BorrowController = borrowCtl
	}
}

func WithReservation(reserv *ReservationController) Option {
	return func(c *Controller) {
		c.ReservationController = reserv
	}
}

func WithCategory(cate *CategoryController) Option {
	return func(c *Controller) {
		c.CategoryController = cate
	}
}

func WithStats(stats *StatsController) Option {
	return func(c *Controller) {
		c.StatsController = stats
	}
}

func NewController(opts ...Option) *Controller {
	ctl := &Controller{}

	for _, opt := range opts {
		opt(ctl)
	}

	return ctl
}
