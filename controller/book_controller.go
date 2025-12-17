package controller

import (
	"strconv"
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

	common.Success(c, 201, "图书添加成功", data)
}

func (ctl *BookController) BatchCreateBook(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.BatchCreateBookRequest
	if err := common.ValidateStruct(c, &req); err != nil {
		c.Error(err)
		return
	}

	data, err := ctl.bookService.BatchCreateBook(ctx, &req)
	if err != nil {
		c.Error(err)
		return
	}

	common.Success(c, 201, "批量导入完成", data)
}

func (ctl *BookController) GetBookList(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.GetBookListRequest
	if err := common.ValidateStruct(c, &req); err != nil {
		c.Error(err)
		return
	}

	data, err := ctl.bookService.GetBookList(ctx, &req)
	if err != nil {
		c.Error(err)
		return
	}

	common.Success(c, 200, "success", data)
}

func (ctl *BookController) GetBookDetails(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(err)
		return		
	}

	data, err := ctl.bookService.GetBookDetails(ctx, id)
	if err != nil {
		c.Error(err)
		return			
	}

	common.Success(c, 200, "success", data)
}

func (ctl *BookController) UpdateBook(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(err)
		return		
	}

	var req request.UpdateBookRequest
	if err := common.ValidateStruct(c, &req); err != nil {
		c.Error(err)
		return
	}

	data, err := ctl.bookService.UpdateBook(ctx, &req, id)
	if err != nil {
		c.Error(err)
		return		
	}

	common.Success(c, 200, "图书更新成功", data)
}

func (ctl *BookController) DeleteBook(c *gin.Context) {
	ctx := c.Request.Context()
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(err)
		return		
	}

	err = ctl.bookService.DeleteBook(ctx, id)
	if err != nil {
		c.Error(err)
		return		
	}
	
	common.Success(c, 200, "图书删除成功", gin.H{})
}