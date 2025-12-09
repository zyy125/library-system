package config

import (
	"os"
)

type Config struct {
	DSN string
}

func Load() *Config {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == ""{
		dsn = "root:@tcp(127.0.0.1:3306)/librarySystem?charset=utf8mb4&parseTime=True&loc=Local"
	}

	return &Config{
		DSN: dsn,
	}
}