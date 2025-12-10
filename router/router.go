package router

import (
	"github.com/gin-gonic/gin"
	"library-system/controller"
	"library-system/middleware"
)

func SetupRouter(userCtl *controller.UserController) *gin.Engine{
	r := gin.Default()

	r.Use(middleware.ErrorHandler())
	api := r.Group("/api")
	{
		users := api.Group("/users")
		{
			users.POST("/register", userCtl.Register)
		}
	}

	return r
}