package router

import (
	"library-system/controller"
	"library-system/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(ctl *controller.Controller) *gin.Engine {
	r := gin.Default()

	userCtl := ctl.UserController
	bookCtl := ctl.BookController

	r.Use(middleware.ErrorHandler())
	r.Use(gin.Recovery())
	api := r.Group("/api")
	{
		users := api.Group("/users")
		{
			users.POST("/register", userCtl.Register)
			users.POST("/login", userCtl.Login)
			users.POST("/refresh-token", userCtl.RefreshToken)

			auth := users.Group("", middleware.AuthMiddleware())
			{
				auth.POST("/logout", userCtl.Logout)
				auth.GET("/me", userCtl.GetUserMsg)
				auth.PUT("/me", userCtl.UpdateUser)
				auth.POST("/change-password", userCtl.ChangePwd)

				admin := auth.Group("", middleware.RoleMiddleware())
				{
					admin.GET("", userCtl.GetUserList)
					admin.POST("", userCtl.CreateUser)
					admin.PUT("/:id", userCtl.UpdateUserByAdmin)
					admin.DELETE("/:id", userCtl.DeleteUser)
				}
			}
		}

		books := api.Group("/books")
		{
			auth := books.Group("", middleware.AuthMiddleware())
			{
				admin := auth.Group("", middleware.RoleMiddleware())
				{
					admin.POST("", bookCtl.CreateBook)
				}
			}
		}
	}

	return r
}
