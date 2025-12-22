package controller

import (
	"library-system/common"
	"library-system/dto/request"
	"library-system/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	categoryService *service.CategoryService
}

func NewCategoryController(service *service.CategoryService) *CategoryController {
	return &CategoryController{
		categoryService: service,
	}
}

// GetCategoryList 获取分类列表
// GET /api/categories
func (ctl *CategoryController) GetCategoryList(c *gin.Context) {
	ctx := c.Request. Context()

	var req request.GetCategoryListRequest
	if err := c. ShouldBindQuery(&req); err != nil {
		c.Error(common.ErrBadRequest)
		return
	}

	data, err := ctl.categoryService.GetCategoryList(ctx, &req)
	if err != nil {
		c. Error(err)
		return
	}

	common.Success(c, 200, "success", data)
}

// GetCategoryDetail 获取分类详情
// GET /api/categories/: id
func (ctl *CategoryController) GetCategoryDetail(c *gin.Context) {
	ctx := c.Request. Context()

	idStr := c. Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(common.ErrBadRequest)
		return
	}

	data, err := ctl.categoryService. GetCategoryDetail(ctx, uint(id))
	if err != nil {
		c.Error(err)
		return
	}

	common.Success(c, 200, "success", data)
}

// CreateCategory 创建分类
// POST /api/categories
func (ctl *CategoryController) CreateCategory(c *gin.Context) {
	ctx := c.Request. Context()

	var req request.CreateCategoryRequest
	if err := common.ValidateStruct(c, &req); err != nil {
		c.Error(err)
		return
	}

	data, err := ctl.categoryService.CreateCategory(ctx, &req)
	if err != nil {
		c. Error(err)
		return
	}

	common.Success(c, 201, "分类创建成功", data)
}

// UpdateCategory 更新分类
// PUT /api/categories/:id
func (ctl *CategoryController) UpdateCategory(c *gin.Context) {
	ctx := c.Request. Context()

	idStr := c. Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(common.ErrBadRequest)
		return
	}

	var req request.UpdateCategoryRequest
	if err := common. ValidateStruct(c, &req); err != nil {
		c.Error(err)
		return
	}

	data, err := ctl.categoryService. UpdateCategory(ctx, uint(id), &req)
	if err != nil {
		c.Error(err)
		return
	}

	common.Success(c, 200, "分类更新成功", data)
}

// DeleteCategory 删除分类
// DELETE /api/categories/:id
func (ctl *CategoryController) DeleteCategory(c *gin. Context) {
	ctx := c.Request.Context()

	idStr := c. Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(common.ErrBadRequest)
		return
	}

	err = ctl.categoryService.DeleteCategory(ctx, uint(id))
	if err != nil {
		c.Error(err)
		return
	}

	common.Success(c, 200, "分类删除成功", gin.H{})
}