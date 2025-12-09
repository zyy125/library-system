package router

import (
	"github.com/gin-gonic/gin"
	"library-system/controller"
)

func SetupRouter(userCtl *controller.UserController) *gin.Engine{
	r := gin.Default()

	api := r.Group("/api")
	{
		users := api.Group("/users")
		{
			users.POST("/register", userCtl.Register)
		}
	}

	return r
}