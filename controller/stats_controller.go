package controller

import (
	"library-system/common"
	"library-system/dto/request"
	"library-system/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StatsController struct {
	statsService *service.StatsService
}

func NewStatsController(service *service.StatsService) *StatsController {
	return &StatsController{
		statsService: service,
	}
}

// GetOverview 获取系统统计概览
// GET /api/stats/overview
func (ctl *StatsController) GetOverview(c *gin.Context) {
	ctx := c.Request.Context()

	data, err := ctl.statsService.GetOverview(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	common.Success(c, 200, "success", data)
}

// GetBorrowStats 获取借阅统计
// GET /api/stats/borrow
func (ctl *StatsController) GetBorrowStats(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.GetBorrowStatsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(common.ErrBadRequest)
		return
	}

	data, err := ctl.statsService.GetBorrowStats(ctx, &req)
	if err != nil {
		c.Error(err)
		return
	}

	common.Success(c, 200, "success", data)
}

// GetUserStats 获取用户借阅统计
// GET /api/stats/user/: user_id
func (ctl *StatsController) GetUserStats(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.Error(common.ErrBadRequest)
		return
	}

	// 权限检查：普通用户只能查看自己的统计
	currentUserID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	if role != "admin" && currentUserID.(uint64) != userID {
		c.Error(common.ErrPermissionDenied)
		return
	}

	data, err := ctl.statsService.GetUserStats(ctx, userID)
	if err != nil {
		c.Error(err)
		return
	}

	common.Success(c, 200, "success", data)
}

// GetPopularBooks 获取热门图书排行
// GET /api/stats/popular-books
func (ctl *StatsController) GetPopularBooks(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.GetPopularBooksRequest
	if err := c. ShouldBindQuery(&req); err != nil {
		c. Error(common.ErrBadRequest)
		return
	}

	data, err := ctl.statsService.GetPopularBooks(ctx, &req)
	if err != nil {
		c.Error(err)
		return
	}

	common.Success(c, 200, "success", data)
}

// GetCategoryStats 获取分类统计
// GET /api/stats/categories
func (ctl *StatsController) GetCategoryStats(c *gin.Context) {
	ctx := c.Request.Context()

	data, err := ctl.statsService.GetCategoryStats(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	common.Success(c, 200, "success", data)
}