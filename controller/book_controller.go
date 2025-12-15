package controller

import (
	"library-system/common"
	"library-system/dto/request"
	"library-system/service"

	"github.com/gin-gonic/gin"
)

type BookController struct {
	bookService *service.BookService
}

func NewBookController(service *service.BookService) *BookController{
	return &BookController{
		bookService: service,
	}
}

func (ctl *BookController) CreateBook(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.CreateBookRequest
	if err := common.ValidateStruct(c, &req); err != nil {
		c.Error(err)
		return
	}

	data, err := ctl.bookService.CreateBook(ctx, &req) 
	if err != nil {
		c.Error(err)
		return
	}

	common.Success(c, 200, "图书添加成功", data)
}