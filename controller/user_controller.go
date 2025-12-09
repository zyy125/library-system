package controller

import (
	"library-system/common"
	"library-system/dto/request"
	"library-system/dto/response"
	"library-system/model"
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

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.Error(err)
		return
	}

	user := model.User{
		Username: req.Username,
		Password: req.Password,
		Email: req.Email,
		Phone: req.Phone,
	}

	err = ctl.userService.Register(ctx, &user)
	if err != nil {
		c.Error(err)
		return
	}

	data := response.UserRegisterResponse{
		ID: int(user.ID),
		Username: user.Username,
		Email: user.Email,
		Role: user.Role,
		CreatedAt: user.CreatedAt,
	}

	common.Success(c, 201, "注册成功", data)
}
