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
	borrowCtl := ctl.BorrowController
	categoryCtl := ctl.CategoryController
	statsCtl := ctl.StatsController

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
					admin.POST("/batch", bookCtl.BatchCreateBook)
					admin.PUT("/:id", bookCtl.UpdateBook)
					admin.DELETE("/:id", bookCtl.DeleteBook)
				}
			}
			books.GET("/:id", bookCtl.GetBookDetails)
			books.GET("", bookCtl.GetBookList)
		}

		borrow := api.Group("/borrow")
		{
			auth := borrow.Group("", middleware.AuthMiddleware())
			{
				auth.POST("", borrowCtl.BorrowBook)
				auth.POST("/:borrow_id/return", borrowCtl.ReturnBook)
				auth.POST("/:borrow_id/renew", borrowCtl.RenewBorrow)
				auth.GET("", borrowCtl.GetBorrowRecordList)
				auth.GET("/current", borrowCtl.GetCurrentRecord)
			}
		}

		reservations := api.Group("/reservations")
        reservations.Use(middleware.AuthMiddleware())
        {
            reservations.POST("", ctl.ReservationController.CreateReservation)
            reservations.DELETE("/:id", ctl.ReservationController.CancelReservation)
            reservations.GET("/my", ctl.ReservationController.GetMyReservations)
        }

		categories := api. Group("/categories")
		{
			// 公开接口（不需要认证）
			categories.GET("", categoryCtl. GetCategoryList)
			categories.GET("/:id", categoryCtl.GetCategoryDetail)

			// 管理员接口
			auth := categories.Group("", middleware.AuthMiddleware())
			{
				admin := auth.Group("", middleware.RoleMiddleware())
				{
					admin. POST("", categoryCtl.CreateCategory)
					admin.PUT("/:id", categoryCtl.UpdateCategory)
					admin.DELETE("/:id", categoryCtl.DeleteCategory)
				}
			}
		}
	}

	stats := api.Group("/stats")
		{
			// 公开接口
			stats.GET("/popular-books", statsCtl.GetPopularBooks)

			// 需要认证的接口
			auth := stats.Group("", middleware. AuthMiddleware())
			{
				auth.GET("/user/:user_id", statsCtl.GetUserStats)

				// 管理员接口
				admin := auth.Group("", middleware.RoleMiddleware())
				{
					admin.GET("/overview", statsCtl.GetOverview)
					admin.GET("/borrow", statsCtl.GetBorrowStats)
					admin.GET("/categories", statsCtl.GetCategoryStats)
				}
			}
		}

	return r
}
