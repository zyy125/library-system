package main

import (
	"library-system/config"
	"library-system/database"
)

func main() {
	cfg := config.Load()
	database.InitMySQL(cfg.DSN)
	database.MigrateSQL()
}