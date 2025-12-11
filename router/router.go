package router

import (
	"github.com/gin-gonic/gin"
	"library-system/controller"
	"library-system/middleware"
)

func SetupRouter(userCtl *controller.UserController) *gin.Engine{
	r := gin.Default()

	r.Use(middleware.ErrorHandler())
	r.Use(gin.Recovery())
	api := r.Group("/api")
	{
		users := api.Group("/users")
		{
			users.POST("/register", userCtl.Register)
			users.POST("/login", userCtl.Login)
			users.POST("/refresh-token", userCtl.RefreshToken)

			// auth := users.Group("", middleware.AuthMiddleware())
			// {

			// }
		}
		
	}

	return r
}