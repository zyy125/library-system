package main

import (
	"library-system/app"
	"library-system/router"
	"library-system/database"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app, err := app.InitApp() 
	if err != nil {
		log.Fatal(err)
	}
	
	r := router.SetupRouter(app.Controller)

	go func() {
		log.Println("服务器启动在 : 8080")
		if err := r.Run(":8080"); err != nil {
			log. Fatalf("服务器启动失败: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall. SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("服务器正在关闭...")
	app.Scheduler.OverdueScheduler.Stop()
	app.Scheduler.ReservationScheduler.Stop()
	database.CloseRedis()
	log.Println("服务器已关闭")

}