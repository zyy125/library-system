package controller

import (
	"library-system/common"
	"library-system/dto/request"
	"library-system/service"

	"github.com/gin-gonic/gin"
)


type BorrowController struct {
	borrowService *service.BorrowService
}

func NewBorrowController(borrowService *service.BorrowService) *BorrowController {
	return &BorrowController{borrowService: borrowService}
}

func (ctl *BorrowController) BorrowBook(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _ := c.Get("user_id")

	var req request.BorrowBookRequest
	if err := common.ValidateStruct(c, &req); err != nil {
		c.Error(err)
		return
	}

	data, err := ctl.borrowService.BorrowBook(ctx, userID.(uint64), &req)
	if err != nil {
		c.Error(err)
		return
	}

	common.Success(c, 201, "借阅成功", data)
}