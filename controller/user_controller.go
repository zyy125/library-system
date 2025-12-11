package controller

import (
	"library-system/common"
	"library-system/dto/request"
	"library-system/service"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{userService: service}
}

func (ctl *UserController) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.UserRegisterRequest
	err := common.ValidateStruct(c, &req)
	if err != nil {
		c.Error(err)
		return
	}

	data, err := ctl.userService.Register(ctx, req)
	if err != nil {
		c.Error(err)
		return
	}

	common.Success(c, 201, "注册成功", data)
}

func (ctl *UserController) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.UserLoginRequest
	err := common.ValidateStruct(c, &req)
	if err != nil {
		c.Error(err)
		return
	}

	data, err := ctl.userService.Login(ctx, req)
	if err != nil {
		c.Error(err)
		return
	}

	common.Success(c, 200, "登录成功", data)
}

func (ctl *UserController) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.UserRefreshTokenRequest
	err := common.ValidateStruct(c, &req) 
	if err != nil {
		c.Error(err)
		return
	}

	data, err := ctl.userService.RefreshToken(ctx, req)
	if err != nil {
		c.Error(err)
		return
	}

	common.Success(c, 200, "Token刷新成功", data)
} 

func (ctl *UserController) Logout (c *gin.Context) {
	ctx := c.Request.Context()

	userID, _ := c.Get("user_id")
	tokenID, _ := c.Get("token_id")
	if err := ctl.userService.Logout(ctx, userID.(uint64), tokenID.(string)); err != nil {
		c.Error(err)
		return
	}
	
	common.Success(c, 200, "登出成功", gin.H{})
}

func (ctl *UserController) GetUserMsg (c *gin.Context) {
	ctx := c.Request.Context()

	userID, _ := c.Get("user_id")

	data, err := ctl.userService.GetUserMsg(ctx, userID.(uint64))
	if err != nil {
		c.Error(err)
		return
	}

	common.Success(c, 200,"success" ,data)
}

func (ctl *UserController) UpdateUser (c *gin.Context) {
	ctx := c.Request.Context()

	userID, _ := c.Get("user_id")
	var req request.UpdateUserRequest
	err := common.ValidateStruct(c, &req)
	if err != nil {
		c.Error(err)
		return
	}

	data, err := ctl.userService.UpdateUser(ctx, userID.(uint64), req)
	if err != nil {
		c.Error(err)
		return
	}

	common.Success(c, 200, "更新成功", data)
}