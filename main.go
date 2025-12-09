package main

import (
	"library-system/app"
	"library-system/router"
	"log"
)

func main() {
	userCtl, err := app.InitApp() 
	if err != nil {
		log.Fatal(err)
	}
	
	r := router.SetupRouter(userCtl)

	r.Run(":8080")
}