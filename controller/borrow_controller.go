package controller

import (
	"library-system/common"
	"library-system/dto/request"
	"library-system/service"
	"strconv"

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

func (ctl *BorrowController) ReturnBook(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("borrow_id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(err)
		return		
	}

	var req request.ReturnBookRequest
	if err := common.ValidateStruct(c, &req); err != nil {
		c.Error(err)
		return
	}

	data, err := ctl.borrowService.ReturnBook(ctx, id, &req)
	if err != nil {
		c.Error(err)
		return
	}

	common.Success(c, 200, "归还成功", data)
}

func (ctl *BorrowController) RenewBorrow(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("borrow_id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(err)
		return		
	}

	userID, _ := c.Get("user_id")

	var req request.RenewBorrowRequest
	if err := common.ValidateStruct(c, &req); err != nil {
		c.Error(err)
		return
	}

	data, err := ctl.borrowService.RenewBorrow(ctx, userID.(uint64), id, &req)
	if err != nil {
		c.Error(err)
		return
	}

	common.Success(c, 200, "", data)
}

func (ctl *BorrowController) GetBorrowRecordList(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.GetBorrowRecordListRequest
	if err := common.ValidateStruct(c, &req); err != nil {
		c.Error(err)
		return
	}

	// 权限检查：普通用户只能查看自己的记录
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	if role != "admin" {
		// 普通用户强制只查询自己的记录
		uid := userID.(uint64)
		req.UserID = &uid
	}

	data, err := ctl.borrowService.GetBorrowRecordList(ctx, &req)
	if err != nil {
		c. Error(err)
		return
	}

	common.Success(c, 200, "success", data)
}

func (ctl *BorrowController) GetCurrentRecord(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _ := c.Get("user_id")

	data, err := ctl.borrowService.GetCurrentRecord(ctx, userID.(uint64))
	if err != nil {
		c.Error(err)
		return
	}

	common.Success(c, 200, "success", data)
}